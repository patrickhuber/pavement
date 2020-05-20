package config

import "github.com/patrickhuber/pavement/types"

// ImageType defines an image specific to an IaaS provider
type ImageType struct {
	Name       string                 `hcl:"name,label"`
	Version    string                 `hcl:"version,label"`
	Properties map[string]interface{} `hcl:"properties,block"`
}

// Image defines an image for use in creating virtual machines
type Image struct {
	Name          string `hcl:"name,label"`
	Version       string `hcl:"version,label"`
	ImageTypeName string `hcl:"image_type_name"`
	ImageType     *ImageType
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (it *ImageType) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitImageType(visitContext, it)
}

// Accept implements the Acceptor interface and allows nodes to be visited by a visitor
func (i *Image) Accept(visitContext *VisitContext, visitor Visitor) types.AggregateError {
	return visitor.VisitImage(visitContext, i)
}
