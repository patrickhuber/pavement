package config

import "github.com/patrickhuber/pavement/types"

// Pavement is the wrapper type for config
type Pavement struct {
	Networks            []*Network            `hcl:"network,block"`
	NetworkTypes        []*NetworkType        `hcl:"network_type,block"`
	Subnets             []*Subnet             `hcl:"subnet,block"`
	VirtualMachines     []*VirtualMachine     `hcl:"virtual_machine,block"`
	VirtualMachineTypes []*VirtualMachineType `hcl:"virtual_machine_type,block"`
	ImageTypes          []*ImageType          `hcl:"image_type,block"`
	Images              []*Image              `hcl:"image,block"`
	Properties          *map[string]string    `hcl:"properties"`
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (p *Pavement) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	ae := visitor.VisitPavement(visitContext, p)
	if ae == nil {
		ae = types.NewAggregateError()
	}
	newVisitContext := &VisitContext{
		Parent: p,
		Root:   p,
	}

	if p.NetworkTypes != nil {
		for _, nt := range p.NetworkTypes {
			err := nt.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if p.Networks != nil {
		for _, n := range p.Networks {
			err := n.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if p.Subnets != nil {
		for _, s := range p.Subnets {
			err := s.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if p.VirtualMachineTypes != nil {
		for _, vt := range p.VirtualMachineTypes {
			err := vt.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if p.VirtualMachines != nil {
		for _, vm := range p.VirtualMachines {
			err := vm.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if p.ImageTypes != nil {
		for _, im := range p.ImageTypes {
			err := im.Accept(newVisitContext, visitor)
			if err != nil {
				ae = ae.Join(err)
			}
		}
	}

	if !ae.IsEmpty() {
		return ae
	}

	return nil
}

// GetVirtualMachineTypeByName returns the given virtual machine type by name
func (p *Pavement) GetVirtualMachineTypeByName(name string) *VirtualMachineType {
	for _, vmt := range p.VirtualMachineTypes {
		if vmt.Name == name {
			return vmt
		}
	}
	return nil
}
