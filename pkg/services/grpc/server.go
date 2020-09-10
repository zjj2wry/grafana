package grpc

import (
	"context"
	"log"
	"net"

	"github.com/grafana/grafana-plugin-sdk-go/genproto/pluginv2"
	"google.golang.org/grpc"
)

type server struct {
	pluginv2.DiagnosticsServer
}

func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	sssss := &server{}

	srv := grpc.NewServer()
	pluginv2.RegisterDiagnosticsServer(srv, *sssss)

	log.Fatalln(srv.Serve(lis))
}

// CheckHealth checks health.
func (s *server) CheckHealth(ctx context.Context, req *pluginv2.CheckHealthRequest) (*pluginv2.CheckHealthResponse, error) {
	return &pluginv2.CheckHealthResponse{
		Status:  pluginv2.CheckHealthResponse_OK,
		Message: "hello",
	}, nil
}

func (s *server) CollectMetrics(ctx context.Context, in *pluginv2.CollectMetricsRequest, opts ...grpc.CallOption) (*pluginv2.CollectMetricsResponse, error) {
	// ???
	return nil, nil
}
