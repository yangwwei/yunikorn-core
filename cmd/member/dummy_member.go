package main

import (
	"context"
	"fmt"
	"github.com/apache/yunikorn-core/pkg/common"
	"github.com/apache/yunikorn-core/pkg/log"
	siCommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"strconv"
	"time"
)

const (
	address = "localhost:3333"
)

func main() {
	if err := runApp(); err != nil {
		log.Log(log.Member).Fatal("unable to start member",
			zap.Error(err))
	}
}

func runApp() error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log(log.Member).Fatal("unable to connect to the head",
			zap.Error(err))
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Log(log.Member).Fatal("failed to close the connection")
		}
	}(conn)

	c := si.NewSchedulerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*100000)
	defer cancel()
	_, err = c.RegisterResourceManager(ctx, &si.RegisterResourceManagerRequest{
		RmID:        "test",
		PolicyGroup: "",
		Version:     "0.0.1",
	})
	if err != nil {
		return fmt.Errorf("unable to register: %v", err)
	}
	log.Log(log.Member).Info("Member registered")

	stream, err := c.UpdateNode(ctx)
	if err != nil {
		return fmt.Errorf("error on update: %v", err)
	}
	done := make(chan bool)

	virtualNodeID1, _ := common.GetVirtualNodeID("compute-cluster-1", "root.sandbox")
	virtualNodeID2, _ := common.GetVirtualNodeID("compute-cluster-2", "root.sandbox")
	nodeRequest := &si.NodeRequest{Nodes: []*si.NodeInfo{
		{
			// virtual node ID
			NodeID: virtualNodeID1,
			Attributes: map[string]string{
				"cost-tier": "low",
			},
			// Max queue resource
			SchedulableResource: &si.Resource{
				Resources: map[string]*si.Quantity{
					"vcore": {Value: 200},
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
						"queue":               "root.sandbox", // TODO: clarify requirement, this is required!
						siCommon.CreationTime: strconv.FormatInt(time.Now().UnixMilli(), 10),
					},
					ResourcePerAlloc: &si.Resource{
						Resources: map[string]*si.Quantity{
							"vcore": {Value: 100},
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
					"vcore": {Value: 100},
				},
			},
			ExistingAllocations: []*si.Allocation{
				{
					// allocation key is the applicationID
					// tracks the "current" allocated resources for the app
					ApplicationID: "job-2",
					AllocationKey: "job-2",
					AllocationID:  "job-2-allocation-1",
					NodeID:        virtualNodeID2,
					AllocationTags: map[string]string{
						"queue":               "root.sandbox", // TODO: clarify requirement, this is required!
						siCommon.CreationTime: strconv.FormatInt(time.Now().UnixMilli(), 10),
					},
					ResourcePerAlloc: &si.Resource{
						Resources: map[string]*si.Quantity{
							"vcore": {Value: 10},
						},
					},
				},
				{
					// allocation key is the applicationID
					// tracks the "current" allocated resources for the app
					ApplicationID: "job-3",
					AllocationKey: "job-3",
					AllocationID:  "job-3-allocation-1",
					NodeID:        virtualNodeID2,
					AllocationTags: map[string]string{
						"queue":               "root.sandbox", // TODO: clarify requirement, this is required!
						siCommon.CreationTime: strconv.FormatInt(time.Now().UnixMilli(), 10),
					},
					ResourcePerAlloc: &si.Resource{
						Resources: map[string]*si.Quantity{
							"vcore": {Value: 20},
						},
					},
				},
			},
			Action: si.NodeInfo_CREATE,
		},
	},
		RmID: "test",
	}

	// Connect to server and send streaming
	// first goroutine sends requests
	if err := stream.Send(nodeRequest); err != nil {
		log.Log(log.Member).Fatal("failed to send the node request", zap.Error(err))
	}

	log.Log(log.Member).Info("node request sent")

	// second goroutine receives data from stream
	// and saves result in max variable
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Log(log.Member).Fatal("failed to get response", zap.Error(err))
			}
			log.Log(log.Member).Info("Responded by server",
				zap.String("response", response.String()))
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Log(log.Member).Fatal("error closing", zap.Error(err))
		}
		close(done)
	}()

	<-done
	log.Log(log.Member).Info("client has shutdown")
	return nil
}
