package config_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/patrickhuber/pavement/config"
)

type tester struct {
}

func NewTester() *tester {
	return &tester{}
}

func (t *tester) validate(p *config.Pavement, file string) {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(p, f.Body())
	var writer bytes.Buffer
	fmt.Fprintf(&writer, "%s", f.Bytes())
	b, err := ioutil.ReadFile(file)
	Expect(err).To(BeNil())
	actual := writer.String()
	expected := string(b)
	if strings.Compare(actual, expected) != 0 {
		fmt.Println("actual:")
		fmt.Println(actual)
		fmt.Println("expected:")
		fmt.Println(expected)
		Fail("strings do not match")
	}
}
