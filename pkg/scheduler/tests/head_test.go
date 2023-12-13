package tests

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/common"
	"github.com/apache/yunikorn-core/pkg/entrypoint"
	"github.com/apache/yunikorn-core/pkg/kubeacon"
	siCommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
	"strconv"
	"testing"
	"time"
)

func TestHead(t *testing.T) {

	configData := `
partitions:
  - name: default
    queues:
      - name: root
        submitacl: "*"
`

	// Not implemented
	memberMgr := &kubeacon.MemberManager{}

	context := entrypoint.StartAllServices()
	response, err := context.RMProxy.RegisterResourceManager(
		&si.RegisterResourceManagerRequest{
			RmID:        "test",
			PolicyGroup: "",
			Version:     "0.0.1",
			Config:      configData,
		}, memberMgr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("registered")
	fmt.Println(response)

	virtualNodeID1, err := common.GetVirtualNodeID("compute-cluster-1", "root.sandbox")
	if err != nil {
		t.Error(err)
	}
	virtualNodeID2, err := common.GetVirtualNodeID("compute-cluster-2", "root.sandbox")
	if err != nil {
		t.Error(err)
	}
	err = context.RMProxy.UpdateNode(&si.NodeRequest{Nodes: []*si.NodeInfo{
		{
			// virtual node ID
			NodeID: virtualNodeID1,
			Attributes: map[string]string{
				"cost-tier": "low",
			},
			// Max queue resource
			SchedulableResource: &si.Resource{
				Resources: map[string]*si.Quantity{
					"memory": {Value: 100000000},
					"vcore":  {Value: 20000},
				},
			},
			ExistingAllocations: []*si.Allocation{
				{
					// allocation key is the applicationID
					// tracks the "current" allocated resources for the app
					ApplicationID: "job-1",
					AllocationKey: "job-1",
					AllocationID:  "job-1-allocation-1",
					NodeID:        virtualNodeID1,
					AllocationTags: map[string]string{
						siCommon.CreationTime: strconv.FormatInt(time.Now().UnixMilli(), 10),
					},
					ResourcePerAlloc: &si.Resource{
						Resources: map[string]*si.Quantity{
							"memory": {Value: 50000000},
							"vcore":  {Value: 10000},
						},
					},
				},
			},
			Action: si.NodeInfo_CREATE,
		}, {
			NodeID: virtualNodeID2,
			Attributes: map[string]string{
				"cost-tier": "high",
			},
			SchedulableResource: &si.Resource{
				Resources: map[string]*si.Quantity{
					"memory": {Value: 200000000},
					"vcore":  {Value: 40000},
				},
			},
			Action: si.NodeInfo_CREATE,
		},
	},
		RmID: "test",
	})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	part := context.Scheduler.GetClusterContext().GetPartition("[test]default")
	for _, node := range part.GetNodes() {
		fmt.Println(node)
	}

	internalQueue := part.GetQueue("root.sandbox")
	fmt.Println(internalQueue.Name)
	fmt.Println(internalQueue.GetMaxResource().String())
}
