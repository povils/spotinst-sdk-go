package mcs

import (
	"context"
	"github.com/povils/spotinst-sdk-go/spotinst"
	"github.com/povils/spotinst-sdk-go/spotinst/client"
	"github.com/povils/spotinst-sdk-go/spotinst/session"
)

// Service provides the API operation methods for making requests to
// endpoints of the Spotinst API. See this package's package overview docs
// for details on the service.
type Service interface {
	GetClusterCosts(context.Context, *ClusterCostInput) (*ClusterCostOutput, error)
}

type ServiceOp struct {
	Client *client.Client
}

var _ Service = &ServiceOp{}

func New(sess *session.Session, cfgs ...*spotinst.Config) *ServiceOp {
	cfg := &spotinst.Config{}
	cfg.Merge(sess.Config)
	cfg.Merge(cfgs...)

	return &ServiceOp{
		Client: client.New(sess.Config),
	}
}
