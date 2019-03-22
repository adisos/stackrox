package image

import (
	"github.com/stackrox/rox/central/image/index/mappings"
	policyDatastore "github.com/stackrox/rox/central/policy/datastore"
	"github.com/stackrox/rox/central/searchbasedpolicies"
	"github.com/stackrox/rox/central/searchbasedpolicies/matcher"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/sync"
)

type setImpl struct {
	lock sync.RWMutex

	policyIDToPolicy  map[string]*storage.Policy
	policyIDToMatcher map[string]searchbasedpolicies.Matcher
	policyStore       policyDatastore.DataStore
}

// ForEach runs the given function on all present policies.
func (p *setImpl) ForEach(fe func(*storage.Policy, searchbasedpolicies.Matcher) error) error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for id, m := range p.policyIDToMatcher {
		if err := fe(p.policyIDToPolicy[id], m); err != nil {
			return err
		}
	}
	return nil
}

// UpsertPolicy adds or updates a policy in the set.
func (p *setImpl) UpsertPolicy(policy *storage.Policy) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	cloned := protoutils.CloneStoragePolicy(policy)

	searchBasedMatcher, err := matcher.ForPolicy(cloned, mappings.OptionsMap, nil)
	if err != nil {
		return err
	}

	p.policyIDToPolicy[cloned.GetId()] = cloned
	p.policyIDToMatcher[cloned.GetId()] = searchBasedMatcher
	return nil
}

// RemovePolicy removes a policy from the set.
func (p *setImpl) RemovePolicy(policyID string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, exists := p.policyIDToPolicy[policyID]; exists {
		delete(p.policyIDToPolicy, policyID)
		delete(p.policyIDToMatcher, policyID)
	}
	return nil
}

// RemoveNotifier removes a given notifier from any policies in the set that use it.
func (p *setImpl) RemoveNotifier(notifierID string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, policy := range p.policyIDToPolicy {
		filtered := policy.GetNotifiers()[:0]
		for _, n := range policy.GetNotifiers() {
			if n != notifierID {
				filtered = append(filtered, n)
			}
		}
		policy.Notifiers = filtered

		err := p.policyStore.UpdatePolicy(policy)
		if err != nil {
			return err
		}
	}

	return nil
}
