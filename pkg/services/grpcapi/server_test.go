package grpcapi

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/genproto/pluginv2"
	"google.golang.org/grpc"
)

func TestRunServer(t *testing.T) {
	t.Skip()
	//
	//server := &GRPCApiServer{}
	//server.RunAPIServer(":50051")

}

// printHealth lists all the features within the given bounding Rectangle.
func printHealth(client pluginv2.DiagnosticsClient, plugin *pluginv2.PluginContext) {
	log.Printf("Looking for features within %v", plugin)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rsp, err := client.CheckHealth(ctx, &pluginv2.CheckHealthRequest{
		PluginContext: plugin,
	})
	if err != nil {

		return
	}

	log.Println("GOT", rsp, err)
}

// Requires grafana core to be running!
func TestRunClient(t *testing.T) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // for now?
	// anythign else?
	opts = append(opts, grpc.WithBlock())

	address := "localhost:3001"

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("unable to connect to %s (%s)", address, err.Error())
		t.FailNow()
	}
	defer conn.Close()

	client := pluginv2.NewDiagnosticsClient(conn)

	printHealth(client, &pluginv2.PluginContext{
		PluginId: "aaa",
	})
}
