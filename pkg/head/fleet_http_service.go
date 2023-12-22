/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package head

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/common"
	"github.com/apache/yunikorn-core/pkg/entrypoint"
	"golang.org/x/net/context"
	"net"
	"os"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
)

func NewFleetHttpService(endpoint string) *FleetHttpService {
	return &FleetHttpService{
		endpoint: endpoint,
		wg:       sync.WaitGroup{},
	}
}

type FleetHttpService struct {
	wg       sync.WaitGroup
	endpoint string
	server   *grpc.Server
	fleetMgr *FleetManager
}

func (s *FleetHttpService) Start() {
	s.wg.Add(1)

	schedulerContext := entrypoint.StartAllServices()
	s.fleetMgr = newFleetManager(schedulerContext)

	go s.serve(s.endpoint, s.fleetMgr)
}

func (s *FleetHttpService) Wait() {
	s.wg.Wait()
}

func (s *FleetHttpService) Stop() {
	s.server.GracefulStop()
}

func (s *FleetHttpService) ForceStop() {
	s.server.Stop()
}

func (s *FleetHttpService) serve(endpoint string, ss si.FleetServer) {
	proto, addr, err := common.ParseEndpoint(endpoint)
	if err != nil {
		log.Log(log.RPC).Fatal("fatal error", zap.Error(err))
	}

	if proto == "unix" {
		addr = "/" + addr
		if err = os.Remove(addr); err != nil && !os.IsNotExist(err) {
			log.Log(log.RPC).Fatal("failed to remove unix domain socket",
				zap.String("uds", addr),
				zap.Error(err))
		}
	}

	var listener net.Listener
	listener, err = net.Listen(proto, addr)
	if err != nil {
		log.Log(log.RPC).Fatal("failed to listen to address",
			zap.Error(err))
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(logGRPC))
	s.server = server

	if ss != nil {
		si.RegisterFleetServer(server, ss)
	}

	log.Log(log.RPC).Info("listening for connections",
		zap.Stringer("address", listener.Addr()))

	if err = server.Serve(listener); err != nil {
		log.Log(log.RPC).Fatal("failed to serve", zap.Error(err))
	}
}

// Logging unary interceptor function to log every RPC call
func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Log(log.RPC).Debug("GPRC call",
		zap.String("method", info.FullMethod))
	log.Log(log.RPC).Debug("GPRC request",
		zap.String("request", fmt.Sprintf("%+v", req)))
	resp, err := handler(ctx, req)
	if err != nil {
		log.Log(log.RPC).Debug("GPRC error", zap.Error(err))
	} else {
		log.Log(log.RPC).Debug("GPRC response",
			zap.String("response", fmt.Sprintf("%+v", resp)))
	}
	return resp, err
}
