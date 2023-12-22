package member

import (
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
)

// MemberManager implements the RM callback API and it is running with member cluster YK instances
type MemberManager struct {
}

func (m *MemberManager) UpdateAllocation(response *si.AllocationResponse) error {
	return nil
}

func (m *MemberManager) UpdateApplication(response *si.ApplicationResponse) error {
	return nil
}

func (m *MemberManager) UpdateNode(response *si.NodeResponse) error {
	return nil
}

func (m *MemberManager) Predicates(args *si.PredicatesArgs) error {
	return nil
}

func (m *MemberManager) PreemptionPredicates(args *si.PreemptionPredicatesArgs) *si.PreemptionPredicatesResponse {
	return nil
}

func (m *MemberManager) SendEvent(events []*si.EventRecord) {

}

func (m *MemberManager) UpdateContainerSchedulingState(request *si.UpdateContainerSchedulingStateRequest) {

}
