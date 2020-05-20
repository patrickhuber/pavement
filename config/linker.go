package config

// Link takes a freshly serialized model and converts it to a linked model
func Link(pavement *Pavement) {
	linkPavement(pavement)
}

func linkPavement(pavement *Pavement) {
	linkVirtualMachineTyps(pavement)
	linkImageTypes(pavement)
	//linkNetworkTypes(pavement)
}

func linkVirtualMachineTyps(pavement *Pavement) {
	vmTypeMap := map[string]*VirtualMachineType{}
	for _, vmt := range pavement.VirtualMachineTypes {
		vmTypeMap[vmt.Name] = vmt
	}

	// if its not found add to an aggregate error?
	for _, vm := range pavement.VirtualMachines {
		vmt := vmTypeMap[vm.VirtualMachineTypeName]
		vm.VirtualMachineType = vmt
	}
}

func linkImageTypes(pavement *Pavement) {
	imageTypeMap := map[string]*ImageType{}
	if pavement.ImageTypes != nil {
		for _, it := range pavement.ImageTypes {
			imageTypeMap[it.Name] = it
		}
	}

	if pavement.Images == nil {
		return
	}
	for _, i := range pavement.Images {
		it := imageTypeMap[i.ImageTypeName]
		i.ImageType = it
	}
}

// func linkNetworkTypes(pavement *Pavement) {
// 	networkTypeMap := map[string]*NetworkType{}
// 	networkMap := map[string]*Network{}

// 	if pavement.NetworkTypes != nil {
// 		for _, nt := range pavement.NetworkTypes {
// 			networkTypeMap[nt.Name] = nt
// 		}
// 	}

// 	if pavement.Networks != nil {
// 		for _, n := range pavement.Networks {
// 			networkMap[n.Name] = n
// 			nt := networkTypeMap[n.NetworkTypeName]
// 			n.NetworkType = nt

// 			if n.Subnets == nil {
// 				continue
// 			}
// 			for _, s := range n.Subnets {
// 				s.Network = n
// 			}
// 		}
// 	}
// 	if pavement.Subnets != nil {
// 		for _, s := range pavement.Subnets {
// 			network := networkMap[*s.NetworkName]
// 			s.Network = network
// 		}
// 	}
// }
