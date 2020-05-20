package config

import "github.com/patrickhuber/pavement/types"

// VisitContext defines the context for the visit
type VisitContext struct {
	Parent interface{}
	Root   *Pavement
}

// PavementVisitor defines an interface for visitin the top level pavement node
type PavementVisitor interface {
	VisitPavement(visitContext *VisitContext, pavement *Pavement) types.AggregateError
}

// NetworkVisitor defines an interface for visiting network nodes
type NetworkVisitor interface {
	VisitNetwork(visitContext *VisitContext, network *Network) types.AggregateError
}

// NetworkTypeVisitor defines an interface for visiting network nodes
type NetworkTypeVisitor interface {
	VisitNetworkType(visitContext *VisitContext, networkType *NetworkType) types.AggregateError
}

// SubnetVisitor defines an interface for visiting subnet nodes
type SubnetVisitor interface {
	VisitSubnet(visitContext *VisitContext, subnet *Subnet) types.AggregateError
}

// VirtualMachineVisitor defines an interface for visiting virtual machine nodes
type VirtualMachineVisitor interface {
	VisitVirtualMachine(visitContext *VisitContext, virtualMachine *VirtualMachine) types.AggregateError
}

// VirtualMachineTypeVisitor defines an interface for visiting virtual machine nodes
type VirtualMachineTypeVisitor interface {
	VisitVirtualMachineType(visitContext *VisitContext, virtualMachineType *VirtualMachineType) types.AggregateError
}

// ImageTypeVisitor defines an interface for visiting image type nodes
type ImageTypeVisitor interface {
	VisitImageType(visitContext *VisitContext, imageType *ImageType) types.AggregateError
}

// ImageVisitor defines an interface for visiting an image node
type ImageVisitor interface {
	VisitImage(visitContext *VisitContext, image *Image) types.AggregateError
}

// Visitor defines an aggregate visitor interface for visiting nodes
type Visitor interface {
	NetworkVisitor
	NetworkTypeVisitor
	SubnetVisitor
	VirtualMachineVisitor
	VirtualMachineTypeVisitor
	ImageTypeVisitor
	ImageVisitor
	PavementVisitor
}

// Acceptor defines an interface for nodes in a graph to be visited by a visitor
type Acceptor interface {
	Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError
}

// PavementFunc defines a function interface for visiting a pavement node
type PavementFunc func(visitContext *VisitContext, pavement *Pavement) types.AggregateError

// NetworkFunc defines a function interface for visiting a network node
type NetworkFunc func(visitContext *VisitContext, network *Network) types.AggregateError

// NetworkTypeFunc defines a function interface for visiting a network type node
type NetworkTypeFunc func(visitContext *VisitContext, networkType *NetworkType) types.AggregateError

// SubnetFunc defines a function interface for visiting a subnet node
type SubnetFunc func(visitContext *VisitContext, subnet *Subnet) types.AggregateError

// VirtualMachineFunc defines a function interface for visiting a virtual machine node
type VirtualMachineFunc func(visitContext *VisitContext, virtualMachine *VirtualMachine) types.AggregateError

// VirtualMachineTypeFunc defines a function interface for visiting a virtual machine type node
type VirtualMachineTypeFunc func(visitContext *VisitContext, virtualMachine *VirtualMachineType) types.AggregateError

// ImageTypeFunc defines a function interface for visiting an image type node
type ImageTypeFunc func(visitContext *VisitContext, imageType *ImageType) types.AggregateError

// ImageFunc defines a function for visiting an image node
type ImageFunc func(visitContext *VisitContext, image *Image) types.AggregateError
