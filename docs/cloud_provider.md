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

subscription "default" {
    provider        = provider.default
    id              = var.subscription_id
    name            = var.subscription_name
}
```