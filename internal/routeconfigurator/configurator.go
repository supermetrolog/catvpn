package routeconfigurator

import "net"

type TrafficRouteConfigurator struct {
}

func New() *TrafficRouteConfigurator {
	return &TrafficRouteConfigurator{}
}

func (t *TrafficRouteConfigurator) RouteToSubnet(subnet net.IPNet) error {

}
