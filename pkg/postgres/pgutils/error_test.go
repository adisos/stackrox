package pgutils

import (
	"context"
	"io"
	"testing"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stretchr/testify/assert"
)

func TestWrappedErrors(t *testing.T) {
	multiErrorEmpty := errorhelpers.NewErrorList("any")

	multiError1Error := errorhelpers.NewErrorList("hello")
	multiError1Error.AddError(context.DeadlineExceeded)

	multiErrorCanceled := errorhelpers.NewErrorList("hello")
	multiErrorCanceled.AddError(context.Canceled)
	multiErrorCanceled.AddErrors(errors.New("other error"))

	multiErrorDeadlineExceeded := errorhelpers.NewErrorList("hello")
	multiErrorDeadlineExceeded.AddError(context.DeadlineExceeded)
	multiErrorDeadlineExceeded.AddErrors(errors.New("other error"))

	cases := []struct {
		err       error
		transient bool
	}{
		{
			err:       errors.New("hello"),
			transient: false,
		},
		{
			err:       errors.Wrap(errors.New("hello"), "hello"),
			transient: false,
		},
		{
			err:       errors.Wrap(context.Canceled, "hello"),
			transient: false,
		},
		{
			err:       errors.Wrap(context.DeadlineExceeded, "hello"),
			transient: false,
		},
		{
			err:       errors.Wrap(errors.Wrap(io.EOF, "1"), "2"),
			transient: true,
		},
		{
			err:       errors.Wrap(errors.Wrap(errors.New("nothing"), "1"), "2"),
			transient: false,
		},
		{
			err:       multiErrorEmpty.ToError(),
			transient: false,
		},
		{
			err:       multiError1Error.ToError(),
			transient: false,
		},
		{
			err:       multiErrorCanceled.ToError(),
			transient: false,
		},
		{
			err:       multiErrorDeadlineExceeded.ToError(),
			transient: false,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.transient, IsTransientError(c.err))
	}
}
