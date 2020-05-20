package config

import (
	"fmt"

	"github.com/patrickhuber/pavement/types"
)

func DefaultObservers() []Observer {
	return []Observer{
		NewNetworkObserver(networkSubnetsDoNotSpecifyNetworkIds),
		NewPavementObserver(topLevelSubnetsMustSpecifyNetworkNames),
	}
}

func networkSubnetsDoNotSpecifyNetworkIds(visitContext *VisitContext, network *Network) types.AggregateError {
	if network.Subnets == nil {
		return nil
	}

	ae := types.NewAggregateError()
	for _, subnet := range network.Subnets {
		if subnet.NetworkName != nil {
			err := fmt.Errorf("Network '%s' contains subnet '%s' with NetworkName specified", network.Name, subnet.Name)
			ae.Append(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func topLevelSubnetsMustSpecifyNetworkNames(visitContext *VisitContext, pavement *Pavement) types.AggregateError {
	if pavement.Subnets == nil {
		return nil
	}

	ae := types.NewAggregateError()
	for _, subnet := range pavement.Subnets {
		if subnet.NetworkName == nil {
			err := fmt.Errorf("Subnet '%s' must specify a NetworkName", subnet.Name)
			ae.Append(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}
