package config

import "github.com/patrickhuber/pavement/types"

// VirtualMachineType defines a virtual machine
type VirtualMachineType struct {
	Name string `hcl:"name,label"`
}

// VirtualMachine defines a virtual machine
type VirtualMachine struct {
	Name                   string `hcl:"name,label"`
	VirtualMachineTypeName string `hcl:"virtual_machine_type_name"`
	VirtualMachineType     *VirtualMachineType
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (vm *VirtualMachine) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitVirtualMachine(visitContext, vm)
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (vmt *VirtualMachineType) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitVirtualMachineType(visitContext, vmt)
}
