package grpcapi

import (
	"context"
	"fmt"
	"net"

	"github.com/grafana/grafana-plugin-sdk-go/genproto/pluginv2"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
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
	pluginv2.DataServer

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
	pluginv2.RegisterDataServer(srv, s)

	// 	https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go#L112

	return srv.Serve(lis)
}

// CheckHealth checks health.
func (s *GRPCApiServer) CheckHealth(ctx context.Context, req *pluginv2.CheckHealthRequest) (*pluginv2.CheckHealthResponse, error) {
	logger.Debug("check health", "req", req)

	// Check a specific datasource config
	if req.PluginContext.DataSourceInstanceSettings != nil {
		datasourceID := req.PluginContext.DataSourceInstanceSettings.Id
		getDsInfo := &models.GetDataSourceByIdQuery{
			OrgId: req.PluginContext.OrgId,
			Id:    datasourceID,
		}

		// ds, err = hs.DatasourceCache.GetDatasource(datasourceID, c.SignedInUser, c.SkipCache)
		// if err != nil {
		// 	hs.log.Debug("Encountered error getting data source", "err", err)
		// 	if errors.Is(err, models.ErrDataSourceAccessDenied) {
		// 		return Error(403, "Access denied to data source", err)
		// 	}
		// 	if errors.Is(err, models.ErrDataSourceNotFound) {
		// 		return Error(400, "Invalid data source ID", err)
		// 	}
		// 	return Error(500, "Unable to load data source metadata", err)
		// }

		if err := bus.Dispatch(getDsInfo); err != nil {
			return nil, fmt.Errorf("Could not find datasource %v", err)
		}

		plugin, exist := s.PluginManager.GetDatasource(getDsInfo.Result.Type)
		if !exist {
			return nil, fmt.Errorf("Unknown plugin type %s", getDsInfo.Result.Type)
		}

		logger.Info("****************************************************", "plugin", plugin)
	}

	return &pluginv2.CheckHealthResponse{
		Status:  pluginv2.CheckHealthResponse_OK,
		Message: "hello",
	}, nil
}

// CollectMetrics gets metrics from a plugin instance
func (s *GRPCApiServer) CollectMetrics(ctx context.Context, req *pluginv2.CollectMetricsRequest) (*pluginv2.CollectMetricsResponse, error) {
	logger.Debug("collect metrics", "req", req)
	return nil, fmt.Errorf("Unsupported call")
}

// QueryData gets metrics from a plugin instance
func (s *GRPCApiServer) QueryData(ctx context.Context, req *pluginv2.QueryDataRequest) (*pluginv2.QueryDataResponse, error) {
	logger.Debug("query data", "req", req)

	// for i, query := range req.Queries {
	// 	query.datasourceID

	// 	sj, err := simplejson.NewJson(query.Json)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	queries[i] = &tsdb.Query{
	// 		RefId:         query.RefId,
	// 		IntervalMs:    query.IntervalMS,
	// 		MaxDataPoints: query.MaxDataPoints,
	// 		QueryType:     query.QueryType,
	// 		DataSource:    getDsInfo.Result,
	// 		Model:         sj,
	// 	}
	// }

	// ???
	return nil, nil
}
