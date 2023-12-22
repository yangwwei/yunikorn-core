package head

import (
	"context"
	"github.com/apache/yunikorn-core/pkg/entrypoint"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FleetManager implements the FleetServer gRPC interface
type FleetManager struct {
	schedulerContext *entrypoint.ServiceContext
	si.UnimplementedFleetServer
}

func newFleetManager(ctx *entrypoint.ServiceContext) *FleetManager {
	return &FleetManager{
		schedulerContext: ctx,
	}
}

func (f *FleetManager) RegisterMember(context context.Context, request *si.RegistrationRequest) (*si.RegistrationResponse, error) {
	log.Log(log.Head).Info("received member registration request",
		zap.String("request", request.String()))

	return nil, nil
}

func (f *FleetManager) Heartbeat(context context.Context, request *si.HeartbeatRequest) (*si.HeartbeatResponse, error) {
	log.Log(log.Head).Info("received heartbeat request",
		zap.String("request", request.String()))

	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
