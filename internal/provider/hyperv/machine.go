package hyperv

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

type Machine struct {
	Name string
}

func (m *Machine) Create(ctx context.Context) error {
	cmd := exec.Command("powershell", "-Command", "New-VM", "-Name", m.Name, "-MemoryStartupBytes", "512MB")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing PowerShell command: %v\nOutput: %s", err, output)
	}
	fmt.Printf("Created virtual machine: %s\n", output)
	return nil
}

func (m *Machine) Delete(ctx context.Context) error {
	cmd := exec.Command("powershell", "-Command", "Remove-VM", "-Name", m.Name, "-Force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing PowerShell command: %v\nOutput: %s", err, output)
	}
	fmt.Printf("Deleted virtual machine: %s\n", output)
	return nil
}
