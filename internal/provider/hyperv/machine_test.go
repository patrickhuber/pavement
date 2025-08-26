package hyperv_test

import (
	"context"
	"testing"

	"github.com/patrickhuber/pavement/internal/provider/hyperv"
)

func TestCreateAndDestroy(t *testing.T) {
	machine := hyperv.Machine{Name: "test-vm"}

	err := machine.Create(context.Background())
	if err != nil {
		t.Fatalf("Failed to create machine: %v", err)
	}

	err = machine.Delete(context.Background())
	if err != nil {
		t.Fatalf("Failed to delete machine: %v", err)
	}
}
