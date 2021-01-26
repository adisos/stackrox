package admissioncontroller

import (
	"context"
	"io"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/features"
	pkgGRPC "github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/grpc/authz/idcheck"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	authorizer = idcheck.AdmissionControlOnly()
)

type managementService struct {
	settingsStream     concurrency.ReadOnlyValueStream
	sensorEventsStream concurrency.ReadOnlyValueStream

	alertHandler AlertHandler
	admCtrlMgr   SettingsManager
}

// NewManagementService retrieves a new admission control management service, that allows pushing config updates out
// to admission control service replicas.
func NewManagementService(mgr SettingsManager, alertHandler AlertHandler) pkgGRPC.APIService {
	return &managementService{
		settingsStream:     mgr.SettingsStream(),
		sensorEventsStream: mgr.SensorEventsStream(),

		alertHandler: alertHandler,
		admCtrlMgr:   mgr,
	}
}

func (s *managementService) RegisterServiceServer(srv *grpc.Server) {
	sensor.RegisterAdmissionControlManagementServiceServer(srv, s)
}

func (s *managementService) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error {
	return nil
}

func (s *managementService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func (s *managementService) runRecv(
	stream sensor.AdmissionControlManagementService_CommunicateServer,
	msgC chan<- *sensor.MsgFromAdmissionControl,
	errC chan<- error) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			errC <- err
			return
		}

		select {
		case <-stream.Context().Done():
			return
		case msgC <- msg:
		}
	}
}

func (s *managementService) sendCurrentSettings(stream sensor.AdmissionControlManagementService_CommunicateServer, settingsIt concurrency.ValueStreamIter) error {
	settings, _ := settingsIt.Value().(*sensor.AdmissionControlSettings)
	if settings == nil {
		return nil
	}
	return stream.Send(&sensor.MsgToAdmissionControl{
		Msg: &sensor.MsgToAdmissionControl_SettingsPush{
			SettingsPush: settings,
		},
	})
}

func (s *managementService) Communicate(stream sensor.AdmissionControlManagementService_CommunicateServer) error {
	if err := stream.SendHeader(metadata.MD{}); err != nil {
		return errors.Wrap(err, "sending header metadata")
	}

	settingsIt := s.settingsStream.Iterator(false)

	if err := s.sendCurrentSettings(stream, settingsIt); err != nil {
		return errors.Wrap(err, "sending initial settings")
	}

	var sensorEventIt concurrency.ValueStreamIter
	if features.K8sEventDetection.Enabled() {
		if err := s.sync(stream); err != nil {
			return errors.Wrap(err, "syncing resources")
		}
		sensorEventIt = s.sensorEventsStream.Iterator(true)
	}

	recvdMsgC := make(chan *sensor.MsgFromAdmissionControl)
	recvErrC := make(chan error, 1)
	go s.runRecv(stream, recvdMsgC, recvErrC)

	for {
		var sensorEventItrDoneC <-chan struct{}
		if sensorEventIt != nil {
			sensorEventItrDoneC = sensorEventIt.Done()
		}

		select {
		case err := <-recvErrC:
			recvErrC = nil // we won't receive anything more on this channel
			if err != nil && err != io.EOF {
				return errors.Wrap(err, "receiving message from admission control service")
			}
		case <-recvdMsgC:
			log.Warn("Received message from admission control service, not sure what to do with it...")
		case <-settingsIt.Done():
			settingsIt = settingsIt.TryNext()
			if err := s.sendCurrentSettings(stream, settingsIt); err != nil {
				return errors.Wrap(err, "sending settings push")
			}
		case <-sensorEventItrDoneC:
			sensorEventIt = sensorEventIt.TryNext()
			if err := s.sendSensorEvent(stream, sensorEventIt); err != nil {
				return errors.Wrap(err, "sending sensor events to admission control service")
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func (s *managementService) PolicyAlerts(_ context.Context, alerts *sensor.AdmissionControlAlerts) (*types.Empty, error) {
	if !features.K8sEventDetection.Enabled() {
		return nil, errors.New("support for kubernetes events policies is not enabled")
	}
	go s.alertHandler.ProcessAlerts(alerts)
	return &types.Empty{}, nil
}

func (s *managementService) sendSensorEvent(stream sensor.AdmissionControlManagementService_CommunicateServer, iter concurrency.ValueStreamIter) error {
	obj, _ := iter.Value().(*sensor.AdmCtrlUpdateResourceRequest)
	if obj == nil {
		return nil
	}

	return stream.Send(&sensor.MsgToAdmissionControl{
		Msg: &sensor.MsgToAdmissionControl_UpdateResourceRequest{
			UpdateResourceRequest: obj,
		},
	})
}

func (s *managementService) sync(stream sensor.AdmissionControlManagementService_CommunicateServer) error {
	for _, msg := range s.admCtrlMgr.GetResourcesForSync() {
		err := stream.Send(&sensor.MsgToAdmissionControl{
			Msg: &sensor.MsgToAdmissionControl_UpdateResourceRequest{
				UpdateResourceRequest: msg,
			},
		})
		if err != nil {
			return err
		}
	}

	err := stream.Send(&sensor.MsgToAdmissionControl{
		Msg: &sensor.MsgToAdmissionControl_UpdateResourceRequest{
			UpdateResourceRequest: &sensor.AdmCtrlUpdateResourceRequest{
				Resource: &sensor.AdmCtrlUpdateResourceRequest_Synced{
					Synced: &sensor.AdmCtrlUpdateResourceRequest_ResourcesSynced{},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
