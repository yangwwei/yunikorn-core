package head

import (
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
)

// NoOptCallback implements the RM callback API but doesn't nothing
type NoOptCallback struct{}

func (m *NoOptCallback) UpdateAllocation(response *si.AllocationResponse) error {
	return nil
}

func (m *NoOptCallback) UpdateApplication(response *si.ApplicationResponse) error {
	return nil
}

func (m *NoOptCallback) UpdateNode(response *si.NodeResponse) error {
	return nil
}

func (m *NoOptCallback) Predicates(args *si.PredicatesArgs) error {
	return nil
}

func (m *NoOptCallback) PreemptionPredicates(args *si.PreemptionPredicatesArgs) *si.PreemptionPredicatesResponse {
	return nil
}

func (m *NoOptCallback) SendEvent(events []*si.EventRecord) {
	return
}

func (m *NoOptCallback) UpdateContainerSchedulingState(request *si.UpdateContainerSchedulingStateRequest) {
	return
}
