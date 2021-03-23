package live

import (
	"github.com/centrifugal/centrifuge"
	"github.com/grafana/grafana/pkg/models"
)

// LiveDataHandler will manages `data/${name}/*` channels
// Each channel maps to a strict schema
// Each channel should (eventually) have saved config stored that describes
// This can accept POST requests to /api/live/${name} and will be converted based on some strategy
// each ${name} will (eventually) have saved configuraiton that specifies what to do with values once
type LiveDataHandler struct {
	Name string
}

// GetHandlerForPath called on init
func (h *LiveDataHandler) GetHandlerForPath(path string) (models.ChannelHandler, error) {
	return h, nil // all dashboards share the same handler
}

// OnSubscribe for now allows anyone to subscribe
func (h *LiveDataHandler) OnSubscribe(c *centrifuge.Client, e centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	return centrifuge.SubscribeReply{}, nil
}

// OnPublish checks if a message from the websocket can be broadcast on this channel
func (h *LiveDataHandler) OnPublish(c *centrifuge.Client, e centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	return centrifuge.PublishReply{}, nil // broadcast any event
}
