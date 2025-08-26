package azure

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var subscriptionId string

var (
	resourcesClientFactory *armresources.ClientFactory
	computeClientFactory   *armcompute.ClientFactory
	networkClientFactory   *armnetwork.ClientFactory
)

var (
	resourceGroupClient *armresources.ResourceGroupsClient

	virtualNetworksClient   *armnetwork.VirtualNetworksClient
	subnetsClient           *armnetwork.SubnetsClient
	securityGroupsClient    *armnetwork.SecurityGroupsClient
	publicIPAddressesClient *armnetwork.PublicIPAddressesClient
	interfacesClient        *armnetwork.InterfacesClient

	virtualMachinesClient *armcompute.VirtualMachinesClient
	disksClient           *armcompute.DisksClient
)

func main() {
	machine := &Machine{
		ResourceGroupName: "sample-resource-group",
		Name:              "sample-vm",
		VnetName:          "sample-vnet",
		SubnetName:        "sample-subnet",
		NsgName:           "sample-nsg",
		NicName:           "sample-nic",
		DiskName:          "sample-disk",
		PublicIPName:      "sample-public-ip",
		Location:          "westus2",
	}
	err := machine.Create(context.Background())
	if err != nil {
		log.Fatalf("cannot create virtual machine:%+v", err)
	}
}

type Machine struct {
	ResourceGroupName string
	SubscriptionID    string
	Name              string
	VnetName          string
	SubnetName        string
	NsgName           string
	NicName           string
	DiskName          string
	PublicIPName      string
	Location          string
}

func (m *Machine) Create(ctx context.Context) error {
	conn, err := connectionAzure()
	if err != nil {
		return err
	}

	resourcesClientFactory, err := armresources.NewClientFactory(m.SubscriptionID, conn, nil)
	if err != nil {
		return err
	}
	resourceGroupClient = resourcesClientFactory.NewResourceGroupsClient()
	networkClientFactory, err = armnetwork.NewClientFactory(m.SubscriptionID, conn, nil)
	if err != nil {
		return err
	}
	virtualNetworksClient = networkClientFactory.NewVirtualNetworksClient()
	subnetsClient = networkClientFactory.NewSubnetsClient()
	securityGroupsClient = networkClientFactory.NewSecurityGroupsClient()
	publicIPAddressesClient = networkClientFactory.NewPublicIPAddressesClient()
	interfacesClient = networkClientFactory.NewInterfacesClient()
	computeClientFactory, err = armcompute.NewClientFactory(subscriptionId, conn, nil)
	if err != nil {
		return fmt.Errorf("cannot create compute client factory:%+v", err)
	}
	virtualMachinesClient = computeClientFactory.NewVirtualMachinesClient()
	disksClient = computeClientFactory.NewDisksClient()

	log.Println("start creating virtual machine...")
	resourceGroup, err := m.createResourceGroup(ctx)
	if err != nil {
		return fmt.Errorf("cannot create resource group:%+v", err)
	}
	log.Printf("Created resource group: %s", *resourceGroup.ID)

	virtualNetwork, err := m.createVirtualNetwork(ctx)
	if err != nil {
		return fmt.Errorf("cannot create virtual network:%+v", err)
	}
	log.Printf("Created virtual network: %s", *virtualNetwork.ID)

	subnet, err := m.createSubnets(ctx)
	if err != nil {
		return fmt.Errorf("cannot create subnet:%+v", err)
	}
	log.Printf("Created subnet: %s", *subnet.ID)

	publicIP, err := m.createPublicIP(ctx)
	if err != nil {
		return fmt.Errorf("cannot create public IP address:%+v", err)
	}
	log.Printf("Created public IP address: %s", *publicIP.ID)

	// network security group
	nsg, err := m.createNetworkSecurityGroup(ctx)
	if err != nil {
		return fmt.Errorf("cannot create network security group:%+v", err)
	}
	log.Printf("Created network security group: %s", *nsg.ID)

	netWorkInterface, err := m.createNetworkInterface(ctx, *subnet.ID, *publicIP.ID, *nsg.ID)
	if err != nil {
		return fmt.Errorf("cannot create network interface:%+v", err)
	}
	log.Printf("Created network interface: %s", *netWorkInterface.ID)

	networkInterfaceID := netWorkInterface.ID
	virtualMachine, err := m.createVirtualMachine(ctx, *networkInterfaceID)
	if err != nil {
		return fmt.Errorf("cannot create virtual machine:%+v", err)
	}

	log.Printf("Created network virual machine: %s", *virtualMachine.ID)
	log.Println("Virtual machine created successfully")

	return nil
}

func (m *Machine) Delete(ctx context.Context) error {

	log.Println("start deleting virtual machine...")
	err := m.deleteVirtualMachine(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete virtual machine:%+v", err)
	}
	log.Println("deleted virtual machine")

	err = m.deleteDisk(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete disk:%+v", err)
	}
	log.Println("deleted disk")

	err = m.deleteNetworkInterface(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete network interface:%+v", err)
	}
	log.Println("deleted network interface")

	err = m.deleteNetworkSecurityGroup(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete network security group:%+v", err)
	}
	log.Println("deleted network security group")

	err = m.deletePublicIP(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete public IP address:%+v", err)
	}
	log.Println("deleted public IP address")

	err = m.deleteSubnets(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete subnet:%+v", err)
	}
	log.Println("deleted subnet")

	err = m.deleteVirtualNetwork(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete virtual network:%+v", err)
	}
	log.Println("deleted virtual network")

	err = m.deleteResourceGroup(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete resource group:%+v", err)
	}
	log.Println("deleted resource group")
	log.Println("success deleted virtual machine.")
	return nil
}

func connectionAzure() (azcore.TokenCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

func (m *Machine) createResourceGroup(ctx context.Context) (*armresources.ResourceGroup, error) {

	parameters := armresources.ResourceGroup{
		Location: to.Ptr(m.Location),
		Tags:     map[string]*string{"sample-rs-tag": to.Ptr("sample-tag")}, // resource group update tags
	}

	resp, err := resourceGroupClient.CreateOrUpdate(ctx, m.ResourceGroupName, parameters, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ResourceGroup, nil
}

func (m *Machine) deleteResourceGroup(ctx context.Context) error {

	pollerResponse, err := resourceGroupClient.BeginDelete(ctx, m.ResourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) createVirtualNetwork(ctx context.Context) (*armnetwork.VirtualNetwork, error) {

	parameters := armnetwork.VirtualNetwork{
		Location: to.Ptr(m.Location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16"), // example 10.1.0.0/16
				},
			},
			//Subnets: []*armnetwork.Subnet{
			//	{
			//		Name: to.Ptr(subnetName+"3"),
			//		Properties: &armnetwork.SubnetPropertiesFormat{
			//			AddressPrefix: to.Ptr("10.1.0.0/24"),
			//		},
			//	},
			//},
		},
	}

	pollerResponse, err := virtualNetworksClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.VnetName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualNetwork, nil
}

func (m *Machine) deleteVirtualNetwork(ctx context.Context) error {

	pollerResponse, err := virtualNetworksClient.BeginDelete(ctx, m.ResourceGroupName, m.VnetName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) createSubnets(ctx context.Context) (*armnetwork.Subnet, error) {

	parameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.Ptr("10.1.10.0/24"),
		},
	}

	pollerResponse, err := subnetsClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.VnetName, m.SubnetName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Subnet, nil
}

func (m *Machine) deleteSubnets(ctx context.Context) error {

	pollerResponse, err := subnetsClient.BeginDelete(ctx, m.ResourceGroupName, m.VnetName, m.SubnetName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) createNetworkSecurityGroup(ctx context.Context) (*armnetwork.SecurityGroup, error) {

	parameters := armnetwork.SecurityGroup{
		Location: to.Ptr(m.Location),
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			SecurityRules: []*armnetwork.SecurityRule{
				// Windows connection to virtual machine needs to open port 3389,RDP
				// inbound
				{
					Name: to.Ptr("sample_inbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description:              to.Ptr("sample network security group inbound port 22"),
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
					},
				},
				// outbound
				{
					Name: to.Ptr("sample_outbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description:              to.Ptr("sample network security group outbound port 22"),
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
					},
				},
			},
		},
	}

	pollerResponse, err := securityGroupsClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.NsgName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.SecurityGroup, nil
}

func (m *Machine) deleteNetworkSecurityGroup(ctx context.Context) error {

	pollerResponse, err := securityGroupsClient.BeginDelete(ctx, m.ResourceGroupName, m.NsgName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) createPublicIP(ctx context.Context) (*armnetwork.PublicIPAddress, error) {

	parameters := armnetwork.PublicIPAddress{
		Location: to.Ptr(m.Location),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic), // Static or Dynamic
		},
	}

	pollerResponse, err := publicIPAddressesClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.PublicIPName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.PublicIPAddress, err
}

func (m *Machine) deletePublicIP(ctx context.Context) error {

	pollerResponse, err := publicIPAddressesClient.BeginDelete(ctx, m.ResourceGroupName, m.PublicIPName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) createNetworkInterface(ctx context.Context, subnetID string, publicIPID string, networkSecurityGroupID string) (*armnetwork.Interface, error) {

	parameters := armnetwork.Interface{
		Location: to.Ptr(m.Location),
		Properties: &armnetwork.InterfacePropertiesFormat{
			//NetworkSecurityGroup:
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr("ipConfig"),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						Subnet: &armnetwork.Subnet{
							ID: to.Ptr(subnetID),
						},
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: to.Ptr(publicIPID),
						},
					},
				},
			},
			NetworkSecurityGroup: &armnetwork.SecurityGroup{
				ID: to.Ptr(networkSecurityGroupID),
			},
		},
	}

	pollerResponse, err := interfacesClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.NicName, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Interface, err
}

func (m *Machine) deleteNetworkInterface(ctx context.Context) error {

	pollerResponse, err := interfacesClient.BeginDelete(ctx, m.ResourceGroupName, m.NicName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) createVirtualMachine(ctx context.Context, networkInterfaceID string) (*armcompute.VirtualMachine, error) {
	//require ssh key for authentication on linux
	//sshPublicKeyPath := "/home/user/.ssh/id_rsa.pub"
	//var sshBytes []byte
	//_,err := os.Stat(sshPublicKeyPath)
	//if err == nil {
	//	sshBytes,err = ioutil.ReadFile(sshPublicKeyPath)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(m.Location),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					// search image reference
					// az vm image list --output table
					Offer:     to.Ptr("WindowsServer"),
					Publisher: to.Ptr("MicrosoftWindowsServer"),
					SKU:       to.Ptr("2019-Datacenter"),
					Version:   to.Ptr("latest"),
					//require ssh key for authentication on linux
					//Offer:     to.Ptr("UbuntuServer"),
					//Publisher: to.Ptr("Canonical"),
					//SKU:       to.Ptr("18.04-LTS"),
					//Version:   to.Ptr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(m.DiskName),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					//DiskSizeGB: to.Ptr[int32](100), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes("Standard_F2s")), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: &armcompute.OSProfile{ //
				ComputerName:  to.Ptr("sample-compute"),
				AdminUsername: to.Ptr("sample-user"),
				AdminPassword: to.Ptr("Password01!@#"),
				//require ssh key for authentication on linux
				//LinuxConfiguration: &armcompute.LinuxConfiguration{
				//	DisablePasswordAuthentication: to.Ptr(true),
				//	SSH: &armcompute.SSHConfiguration{
				//		PublicKeys: []*armcompute.SSHPublicKey{
				//			{
				//				Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "sample-user")),
				//				KeyData: to.Ptr(string(sshBytes)),
				//			},
				//		},
				//	},
				//},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(networkInterfaceID),
					},
				},
			},
		},
	}

	pollerResponse, err := virtualMachinesClient.BeginCreateOrUpdate(ctx, m.ResourceGroupName, m.Name, parameters, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualMachine, nil
}

func (m *Machine) deleteVirtualMachine(ctx context.Context) error {

	pollerResponse, err := virtualMachinesClient.BeginDelete(ctx, m.ResourceGroupName, m.Name, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) deleteDisk(ctx context.Context) error {

	pollerResponse, err := disksClient.BeginDelete(ctx, m.ResourceGroupName, m.DiskName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
