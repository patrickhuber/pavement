package config

import "github.com/patrickhuber/pavement/types"

// Observer defines an interface for observing visitation to a node. It enables a walker to walk the tree and fire events.
type Observer interface {
	OnPavementVisited(visitContext *VisitContext, pavement *Pavement) types.AggregateError
	OnNetworkVisited(visitContext *VisitContext, network *Network) types.AggregateError
	OnNetworkTypeVisited(visitContext *VisitContext, networkType *NetworkType) types.AggregateError
	OnVirtualMachineTypeVisited(visitContext *VisitContext, virtualMachineType *VirtualMachineType) types.AggregateError
	OnVirtualMachineVisited(visitContext *VisitContext, virtualMachine *VirtualMachine) types.AggregateError
	OnSubnetVisited(visitContext *VisitContext, subnet *Subnet) types.AggregateError
	OnImageVisited(visitContext *VisitContext, image *Image) types.AggregateError
	OnImageTypeVisited(visitContext *VisitContext, imageType *ImageType) types.AggregateError
}

func NewPavementObserver(f PavementFunc) Observer {
	return &observer{
		pavementFunc: f,
	}
}

func NewNetworkObserver(f NetworkFunc) Observer {
	return &observer{
		networkFunc: f,
	}
}

type observer struct {
	pavementFunc           PavementFunc
	networkFunc            NetworkFunc
	networkTypeFunc        NetworkTypeFunc
	subnetFunc             SubnetFunc
	virtualMachineTypeFunc VirtualMachineTypeFunc
	virtualMachineFunc     VirtualMachineFunc
	imageFunc              ImageFunc
	imageTypeFunc          ImageTypeFunc
}

func (o *observer) OnPavementVisited(visitContext *VisitContext, pavement *Pavement) types.AggregateError {
	if o.pavementFunc == nil {
		return nil
	}
	return o.pavementFunc(visitContext, pavement)
}

func (o *observer) OnNetworkVisited(visitContext *VisitContext, network *Network) types.AggregateError {
	if o.networkFunc == nil {
		return nil
	}
	return o.networkFunc(visitContext, network)
}

func (o *observer) OnNetworkTypeVisited(visitContext *VisitContext, networkType *NetworkType) types.AggregateError {
	if o.networkTypeFunc == nil {
		return nil
	}
	return o.networkTypeFunc(visitContext, networkType)
}

func (o *observer) OnVirtualMachineTypeVisited(visitContext *VisitContext, virtualMachineType *VirtualMachineType) types.AggregateError {
	if o.virtualMachineTypeFunc == nil {
		return nil
	}
	return o.virtualMachineTypeFunc(visitContext, virtualMachineType)
}

func (o *observer) OnVirtualMachineVisited(visitContext *VisitContext, virtualMachine *VirtualMachine) types.AggregateError {
	if o.virtualMachineFunc == nil {
		return nil
	}
	return o.virtualMachineFunc(visitContext, virtualMachine)
}

func (o *observer) OnSubnetVisited(visitContext *VisitContext, subnet *Subnet) types.AggregateError {
	return nil
}

func (o *observer) OnImageVisited(visitContext *VisitContext, image *Image) types.AggregateError {
	return nil
}

func (o *observer) OnImageTypeVisited(visitContext *VisitContext, imageType *ImageType) types.AggregateError {
	return nil
}
