# Cloud Provider

The cloud provider is an abstraction layer where cloud (aws, gcp, azure, docker, vmware) specific details are defined. These then are applied to instances of the absractions.


## Schema

## Example

> Azure

```hcl
provider {
    name                = "azure"
    alias               = "default"
    version             = "v1"
    tenant_id           = "a5f5bd35-2d84-495b-86bb-e55327dcb06e"        
    subscription_id     = var.subscription_id
    subscription_name   = var.subscription_name
    resource_group_name = var.resource_group_name
}

identity {
    provider        = provider.default
    type            = "client_id"
    name            = "pavement"
    alias           = "default"
    client_id       = secret.client_id
    client_secret   = secret.client_secret
}

region {
    name     = "eastus"
    provider = provider.default
    alias    = "primary"
}

zone {    
    name     = "1"
    alias    = "z1"
    provider = provider.default
    region   = region.primary    
}

zone {
    name     = "2"
    alias    = "z2"
    provider = provider.default
    region   = region.primary        
}

zone {
    name     = "3"
    alias    = "z3"
    provider = provider.default
    region   = region.primary
}
```

> AWS

```
provider {
    name    = "aws"
    alias   = "default"
    version = "v1"
}

identity {
    type            = "client_id"
    name            = "pavement"
    alias           = "default"
    client_id       = secrets.client_id
    clennt_secret   = secrets.client_secret
}

region {
    name    = "us-west-1"
    alias   = "primary"
}

zone {
    name     = "us-west-1a"
    alias    = "z1"
    provider = provider.default
    region   = region.primary    
}

zone {
    name     = "us-west-1b"
    alias    = "z2"
    provider = provider.default
    region   = region.primary    
}

zone {
    name     = "us-west-1c"
    alias    = "z3"
    provider = provider.default
    region   = region.primary
}
```

> Google