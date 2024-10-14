# Cloud Provider

The cloud provider is an abstraction layer where cloud (aws, gcp, azure, docker, vmware) specific details are defined. These then are applied to instances of the absractions.


## Schema

```hcl
provider {
    name      = "azure"
    alias     = "default"
    version   = "v1"
    tenant_id = "a5f5bd35-2d84-495b-86bb-e55327dcb06e"    
}

identity {
    provider        = provider.default
    type            = "client_id"
    name            = "pavement"
    alias           = "default"
    client_id       = secret.client_id
    client_secret   = secret.client_secret
}

subscription {
    provider        = provider.default
    id              = var.subscription_id
    name            = var.subscription_name
    alias           = "default"
}

region {
    name     = "eastus"
    provider = provider.default
    alias    = "primary"
}

availability_zone {    
    name     = "az1"
    alias    = "az1"
    provider = provider.default
    region   = region.primary    
}

availability_zone {
    name     = "az2"
    alias    = "az2"
    provider = provider.default
    region   = region.primary    
}

availability_zone {
    name     = "az3"
    alias    = "az3"
    provider = provider.default
    region   = region.primary    
}
```