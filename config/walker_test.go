package config_test

import (
	"errors"
"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/pavement/config"
	"github.com/patrickhuber/pavement/types"
)

var _ = Describe("Walker", func() {
	Describe("Walk", func() {
		Context("WhenEmptyPavementStructAndNoError", func() {
			It("should return nil error", func() {
				observer := config.NewNetworkObserver(func(visitContext *config.VisitContext, network *config.Network) types.AggregateError {
					return nil
				})
				walker := config.NewWalker([]config.Observer{
					observer,
				})
				p := &config.Pavement{}
				err := walker.Walk(p)
				Expect(err).To(BeNil())
			})
		})
		Context("WhenErrorAtEachLevel", func() {
			It("should return all errors", func() {
				p := &config.Pavement{
					Networks: []*config.Network{
						{
							Name: "Network",
						},
					},
				}
				walker := config.NewWalker([]config.Observer{
					config.NewNetworkObserver(func(visitContext *config.VisitContext, network *config.Network) types.AggregateError {
						return types.NewAggregateError(fmt.Errorf("Error Network %s", network.Name))
					}),
					config.NewPavementObserver(func(visitContext *config.VisitContext, pavement *config.Pavement) types.AggregateError {
						return types.NewAggregateError(errors.New("Error Pavement"))
					}),
				})
				e := walker.Walk(p)
				Expect(e).ToNot(BeNil())
				Expect(e.Length()).To(Equal(2))
			})
		})
	})
})
