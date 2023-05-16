package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/pavement/config"
)

type ruleTester struct {
	walker config.Walker
}

func NewRuleTester() *ruleTester {
	observers := config.DefaultObservers()
	walker := config.NewWalker(observers)
	return &ruleTester{
		walker: walker,
	}
}
func (t *ruleTester) ExpectErrorCountToMatch(p *config.Pavement, count int) {
	err := t.walker.Walk(p)
	Expect(err).ToNot(BeNil())
	Expect(err.Length()).To(Equal(count))
}

var _ = Describe("DefaultObservers", func() {
	var (
		tester *ruleTester
	)
	BeforeEach(func() {
		tester = NewRuleTester()
	})
	Describe("Networks", func() {
		It("embeddedSubnetShouldNotSpecifyNetworkName", func() {
			networkId := "networkId"
			p := &config.Pavement{
				Networks: []*config.Network{
					{
						Subnets: []*config.Subnet{
							{
								NetworkName: &networkId,
							},
						},
					},
				},
			}
			tester.ExpectErrorCountToMatch(p, 1)
		})
		It("topLevelSubnetShouldSpecifyNetworkName", func() {
			p := &config.Pavement{
				Subnets: []*config.Subnet{
					{
						Name: "Subnet",
					},
				},
			}
			tester.ExpectErrorCountToMatch(p, 1)
		})
	})
})
