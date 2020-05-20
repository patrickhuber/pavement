package config

import "github.com/patrickhuber/pavement/types"

// NetworkType defines the network type. This type encapsulates specific iaas properties and exposes abstract iaas properties.
type NetworkType struct {
	Name string `hcl:"name,label"`
}

// Network defines a network structure that is translated to the particular IaaS vendor network
type Network struct {
	Name            string    `hcl:"name,label"`
	NetworkTypeName string    `hcl:"network_type_name"`
	Cidr            string    `hcl:"cidr"`
	Subnets         []*Subnet `hcl:"subnets,block"`
	NetworkType     *NetworkType
}

// Subnet defines a network subnet
type Subnet struct {
	Name        string  `hcl:"name,label"`
	NetworkName *string `hcl:"network_name,optional"`
	Cidr        string  `hcl:"cidr"`
	Network     *Network
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (nt *NetworkType) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitNetworkType(visitContext, nt)
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (n *Network) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {

	ae := visitor.VisitNetwork(visitContext, n)
	if ae == nil {
		ae = types.NewAggregateError()
	}

	newVisitContext := &VisitContext{
		Parent: n,
		Root:   visitContext.Root,
	}
	if n.Subnets == nil {
		if !ae.IsEmpty() {
			return ae
		}
		return nil
	}

	for _, s := range n.Subnets {
		err := s.Accept(newVisitContext, visitor)
		if err != nil {
			ae.Append(err)
		}
	}
	if !ae.IsEmpty() {
		return ae
	}
	return nil
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (s *Subnet) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitSubnet(visitContext, s)
}
