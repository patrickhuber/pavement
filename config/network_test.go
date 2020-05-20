package config_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/patrickhuber/pavement/config"
)

var _ = Describe("Network", func() {
	It("renders entity subnet", func() {
		networkType := &config.NetworkType{
			Name: "test_network_type",
		}
		network := &config.Network{
			Name: "test_network",
			Cidr: "10.0.0.0/16",
		}
		dmzSubnet := &config.Subnet{
			Name: "dmz",
			Cidr: "10.0.1.0/24",
		}
		intranetSubnet := &config.Subnet{
			Name: "intranet",
			Cidr: "10.0.2.0/24",
		}
		p := &config.Pavement{
			Networks: []*config.Network{
				network,
			},
			NetworkTypes: []*config.NetworkType{
				networkType,
			},

			Subnets: []*config.Subnet{
				dmzSubnet,
				intranetSubnet,
			},
		}
		NewTester().validate(p, "fixtures/network_subnet_entity.txt")
	})
	It("renders embedded subnet", func() {
		networkType := &config.NetworkType{
			Name: "test_network_type",
		}

		dmzSubnet := &config.Subnet{
			Name: "dmz",
			Cidr: "10.0.1.0/24",
		}
		intranetSubnet := &config.Subnet{
			Name: "intranet",
			Cidr: "10.0.2.0/24",
		}
		network := &config.Network{
			Name: "test_network",
			Cidr: "10.0.0.0/16",
			Subnets: []*config.Subnet{
				dmzSubnet,
				intranetSubnet,
			},
		}
		p := &config.Pavement{
			Networks: []*config.Network{
				network,
			},
			NetworkTypes: []*config.NetworkType{
				networkType,
			},
		}
		NewTester().validate(p, "fixtures/network_subnet_embedded.txt")
	})
})
