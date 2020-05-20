package config

import "github.com/patrickhuber/pavement/types"

// Walker registers itself as a visitor and calls observers upon each node visitation
type Walker interface {
	Walk(p *Pavement) types.AggregateError
}

type walker struct {
	observers []Observer
}

func (w *walker) Walk(p *Pavement) types.AggregateError {
	visitContext := &VisitContext{
		Parent: nil,
		Root:   p,
	}
	return p.Accept(visitContext, w)
}

func (w *walker) VisitPavement(visitContext *VisitContext, pavement *Pavement) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnPavementVisited(visitContext, pavement)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitImage(visitContext *VisitContext, image *Image) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnImageVisited(visitContext, image)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitImageType(visitContext *VisitContext, imageType *ImageType) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnImageTypeVisited(visitContext, imageType)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitNetwork(visitContext *VisitContext, network *Network) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnNetworkVisited(visitContext, network)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitNetworkType(visitContext *VisitContext, networkType *NetworkType) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnNetworkTypeVisited(visitContext, networkType)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitSubnet(visitContext *VisitContext, subnet *Subnet) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnSubnetVisited(visitContext, subnet)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitVirtualMachine(visitContext *VisitContext, virtualMachine *VirtualMachine) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnVirtualMachineVisited(visitContext, virtualMachine)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

func (w *walker) VisitVirtualMachineType(visitContext *VisitContext, virtualMachineType *VirtualMachineType) types.AggregateError {
	ae := types.NewAggregateError()
	for _, o := range w.observers {
		err := o.OnVirtualMachineTypeVisited(visitContext, virtualMachineType)
		if err != nil {
			ae = ae.Join(err)
		}
	}

	if ae.IsEmpty() {
		return nil
	}
	return ae
}

// NewWalker creates a new walker for the given visitors
func NewWalker(observers []Observer) Walker {
	return &walker{
		observers: observers,
	}
}
