package head

import (
	"context"
	"github.com/apache/yunikorn-core/pkg/common"
	"github.com/apache/yunikorn-core/pkg/entrypoint"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
	"go.uber.org/zap"
	"io"
)

type Head struct {
	schedulerContext *entrypoint.ServiceContext
	si.UnimplementedSchedulerServer
}

func Run(endpoint string) {
	// Create gRPC servers
	ctx := entrypoint.StartAllServices()
	ss := newHeadService(ctx)
	s := common.NewNonBlockingGRPCServer()
	s.Start(endpoint, ss)
	s.Wait()
}

func newHeadService(ctx *entrypoint.ServiceContext) si.SchedulerServer {
	return &Head{
		schedulerContext: ctx,
	}
}

// RegisterResourceManager handles  the member registration
func (h *Head) RegisterResourceManager(ctx context.Context, in *si.RegisterResourceManagerRequest) (*si.RegisterResourceManagerResponse, error) {
	log.Log(log.Head).Info("received Member registration request",
		zap.String("request", in.String()))

	response, err := h.schedulerContext.RMProxy.RegisterResourceManager(in, &NoOptCallback{})
	if err != nil {
		panic(err)
	}

	return response, nil
}

func (h *Head) UpdateAllocation(conn si.Scheduler_UpdateAllocationServer) error {
	// Intended left not implemented
	panic("Not implemented")
}

func (h *Head) UpdateApplication(conn si.Scheduler_UpdateApplicationServer) error {
	// Intended left not implemented
	panic("Not implemented")
}

func (h *Head) UpdateNode(conn si.Scheduler_UpdateNodeServer) error {
	ctx := conn.Context()

	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request, err := conn.Recv()
		if err == io.EOF {
			log.Log(log.Head).Info("stream closed")
			return nil
		}
		if err != nil {
			log.Log(log.Head).Error("received error", zap.Error(err))
			continue
		}

		log.Log(log.Head).Info("received UpdateNode request",
			zap.String("request", request.String()))

		err = h.schedulerContext.RMProxy.UpdateNode(request)
		if err != nil {
			// FIXME
			panic(err.Error())
		}

		acceptedNodes := make([]*si.AcceptedNode, 0)
		for _, n := range request.Nodes {
			acceptedNodes = append(acceptedNodes, &si.AcceptedNode{
				NodeID: n.NodeID,
			})
		}

		// Send response to stream
		response := si.NodeResponse{
			Rejected: nil,
			Accepted: acceptedNodes,
		}

		if err := conn.Send(&response); err != nil {
			return err
		}

	}
}
