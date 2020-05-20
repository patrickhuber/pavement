# Fault Zones

Each cloud provider and Vsphere provide mechanisms for fault tolerance. In practice, many organizations also implement additional mechanisms for fault tolerance. Zones can be formed for risk management, physical location or fault tolerance.

## Regions

A region is a physical location typically Wide Area Network (WAN) connected zone to the company's network. Regions are often large distances away from other regions to ensure disaster events are isolated to a single region and do not impact other regions. For this reason, regions also must maintain some level of autonomy and critical systems must be distributed among regions for redundancy.

### Region

A region is implemented with the "region" resource. This resource maps a pavement availbility zone to a physical region in the cloud provider.

There is no default region map and these must be defined by the user. Here are some examples for the major providers.

> azure

```
resource "azure_region" "pavement" {
    location    = "eastus"    
}
```

> aws

```
resource "aws_region" "pavement" {
    code    = "us-east-2"    
}
```

> gcp

```
```

> vsphere

```
```

## Availibility Zones

An availability zone is usually some segmentation of a region. Cloud providers ofen implement availibility zones as distinct physical buildings or partitions of existing buildings with distinct power, network and storage. 

On prem data centers have a different concept of availibility zones, for example vsphere, uses clusters and resource pools to create distinct availibilty zones. 

An Availbiity zone is implemented with the "availibility zone" resource. This resource maps a pavement availbility zone to a physical availibilty zone in the cloud provider.

## Environments

An environment is a method of risk management that enables lower risk tiers of a platform or application to receive updates before higher risk tiers. This ensures that issues are "shifted left" to the lower risk environments and hopefully resolved before moving to higher risk tiers. Testing is often the mechanism for discovering these bugs and automation is often utilized to detect bugs at scale. 