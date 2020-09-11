package grpcapi

import (
	"context"
	"net"

	"github.com/grafana/grafana-plugin-sdk-go/genproto/pluginv2"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/services/live"
	"google.golang.org/grpc"
)

var (
	logger = log.New("grpcapi")
)

// GRPCApiServer exposes grafana API over GRPC
type GRPCApiServer struct {
	// Services
	pluginv2.DiagnosticsServer

	// Dependencies
	BackendPluginManager *backendplugin.Manager
	PluginManager        *plugins.PluginManager
	Live                 *live.GrafanaLive
}

func main() {
	logger.Info("Server running ...")

	sssss := &GRPCApiServer{}
	sssss.RunAPIServer(":50051")
}

// RunAPIServer runs the api server
func (s *GRPCApiServer) RunAPIServer(addr string) error {

	logger.Info("****************************************************")
	logger.Info("****************************************************")
	logger.Info("****************************************************")
	logger.Info("Starting GRPC api server", "addr", addr)
	logger.Info("****************************************************")
	logger.Info("****************************************************")
	logger.Info("****************************************************")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pluginv2.RegisterDiagnosticsServer(srv, s)

	return srv.Serve(lis)
}

// CheckHealth checks health.
func (s *GRPCApiServer) CheckHealth(ctx context.Context, req *pluginv2.CheckHealthRequest) (*pluginv2.CheckHealthResponse, error) {
	logger.Debug("check health", "req", req)

	return &pluginv2.CheckHealthResponse{
		Status:  pluginv2.CheckHealthResponse_OK,
		Message: "hello",
	}, nil
}

// CollectMetrics gets metrics from a plugin instance
func (s *GRPCApiServer) CollectMetrics(ctx context.Context, req *pluginv2.CollectMetricsRequest) (*pluginv2.CollectMetricsResponse, error) {
	logger.Debug("collect metrics", "req", req)

	// ???
	return nil, nil
}
