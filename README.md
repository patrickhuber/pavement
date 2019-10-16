# Pavement

Pavement provides a cli to scaffold and manage generic infrastructure across multiple clouds and private clouds. 

## Compute

Compute resources are created with a Dockerfile like syntax called Pavementfile. 

### Pavementfile

Pavementfile is cross provider implementation that compiles down to Packer scripts. 

Pavementfile is primarly concerned with what goes inside the VM and the VM inputs/outputs not what goes outside the VM. When a Pavementfile is done executing a vm will be produced.

> example

```
FROM ubuntu:18.04
```

#### FROM command

#### RUN command

### Images

Pavement images map cloud resources to pavement image name and version.

> azure

```
# creates an azure image map mapping an azure base image to a pavement image and alias
# `FROM ubuntu:18.04`
resource "azure_image" "ubuntu_18_04" {
    name        = "ubuntu"
    version     = "18.04"

    # azure specific details
    publisher   = "Canonical"
    offer       = "UbuntuServer"
    sku         = "18.04-LTS"
    version     = "latest"
}
```

> aws

```
# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:18.04`
resource "aws_image" "ubuntu_18_04"{
    name                = "ubuntu"
    version             = "18.04"

    # aws specific details
    image_name          = "ubuntu/images/hvm-ssd/ubuntu-xenial-18.04-amd64-server-*"
    virtualization_type = "hvm"
    owner               = "099720109477" # canonical
}
```

> gcp

```
# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:18.04`
resource "gcp_image" "ubuntu_18_04"{
    name        = "ubuntu"
    version     = "18.04"

    # gcp specific details
    # image_name  = "" # mutually exclusive with `family`
    family      = "ubuntu-1804-lts" # mutually exclusive with `image_name`
    project     = "ubuntu-os-cloud"
}
```

> vsphere

```
# creates a vsphere map mapping an iso or clone template to pavement image 
# `FROM ubuntu:18.04`
resource "vsphere_image" "ubuntu_18_04" {
    name        = "ubuntu"
    version     = "18.04"

    # vsphere specific details
    iso_path        = ""    # Either url or datastore path. Mutually exclusive with `template_path`.
    iso_datasotre   = ""    # Specify the iso datastore. Not needed if using clone templates
    template_path   = ""    # Speicify the template vm path. Mutually exclusive with `iso_path`.
}
```

### Image Registry

An image registry is a registration of provided images. It is meant to abstract image maps from the discovery of those maps to avoid coding vendor
specific image names. 

#### Default Registry

The defualt registry provides the following cross platform maps

> azure

| name   | version | os    | publisher | offer        | sku       | version |
| ------ | ------- | ----- | --------- | ------------ | --------- | ------- |
| ubuntu | 18.04   | linux | Canonical | UbuntuServer | 18.04-LTS | latest  |

> aws

| name   | version | os    | image_name                                               | virtualization_type |
| ------ | ------- | ----- | -------------------------------------------------------- | ------------------- |
| ubuntu | 18.04   | linux | ubuntu/images/hvm-ssd/ubuntu-xenial-18.04-amd64-server-* | hvm                 |

> gcp

| name   | version | os    | family          | project         |
| ------ | ------- | ----- | --------------- | --------------- |
| ubuntu | 18.04   | linux | ubuntu-1804-lts | ubuntu-os-cloud |

> vsphere

No defaults provided. See "vsphere" under images on how to create image_maps for vsphere.

> docker

| name   | version | os    | image  | tag   |
| ------ | ------- | ----- | ------ | ----- |
| ubuntu | 18.04   | linux | ubuntu | 18.04 |

#### Override default registry values

You can override default registry values by specifying an image map for the specified image and removing the "registry" attribute.

#### Registry Image Lookup

You can look up an image from the registry by using the registry name and type. The "default" image_registry name is reserved.

#### Custom registry

To create a custom registry you can use the `image_registry` resource

```
resource "image_registry" "pavement" {
    name = "pavement"    
}
```

You can then specify your images and map them to the image registry

```
resource "azure_image" "ubuntu_18_04" { 
    name        = "ubuntu"
    version     = "18.04"

    # specify custom registry
    registry_id = image_registry.pavement.id

    # azure specific details
    publisher   = "Canonical"
    offer       = "UbuntuServer"
    sku         = "18.04-LTS"
    version     = "latest"    
}
```

## Compute Types

Pavement compute types map vm sizes to pavement compute types. 

> azure

```
resource "azure_compute_type" "default" {
    name            = "default"
    instance_type   = "Standard_D1_v2"
}
```

> aws

```
resource "aws_compute_type" "default" {
    name                = "default"
    instance_type       = "m4.large"
    ephemeral_disk_size = "25GB"
}
```

> gcp

```
resource "gcp_compute_type" "default" {
    name            = "default"
    machine_type    = "n1-standard-2"
    root_disk_size  = "50GB"
    root_disk_type  = "pd-ssd"
}
```

> vsphere

```
resource "vsphere_compute_type" "default" {
    name    = "default"
    cpu     = "2"
    ram     = 1GB
    disk    = 30GB
}
```

## Compute Type Registry

### Default Registry

The defualt registry provides the following cross platform maps

> azure

| name    | instance_type  |
| ------- | -------------  |
| default | Standard_D1_v2 |

> aws

| name    | instance_type  | ephemeral_disk_size |
| ------- | -------------  | ------------------- |
| default | m4.large       | 25MB                |

> gcp

| name    | machine_type  | root_disk_size      | root_disk_type |
| ------- | ------------- | ------------------- | -------------- |
| default | n1-standard-2 | 50GB                | pd-ssd         |

> vsphere

| name    | cpu | ram    | disk   |
| ------- | --- | ------ | ------ |
| default | 2   | 1024MB | 30GB   |

> docker

| name    |
| ------- |
| default |

## Packages

Packages are versioned instances of software. Packages contains assets in the form of compiled bits and/or source code. If the compiled bits are not specified the package will be compiled upon consumption.

A package is a .tar.gz, .tgz, .tar or .zip file that contains a metadata package.yml manifest file at its root

```
.
├── package.yml
├── Dockerfile
└── platforms    
    ├── darwin
    │   └── Dockerfile
    ├── linux
    │   └── Dockerfile
    └── windows
        ├── Dockerfile
        └── amd64
            └── Dockerfile
```

A package should be self contained. 

### Package Manifest

A minimal package manifest includes the name and version of the package.

```
name: package-name
version: 1.0.0
```

A package can also specify platforms to build for and architectures to build for. The default Dockerfile in the package root is used to build the package. If overrides are needed, they are placed in the platforms/{os}/{architecture} and platforms/{os} folder by convention and that Dockerfile will be used instead. All package implementors should provide at a minimum a posix Dockerfile. A windows Dockerfile is desired as well to make sure pavement stays cross platform. 

os is a valid GOOS value
architecture is a valid GOARCH value

The docker file in each of these folders should assume all assets are inside the container or in container mounts. There are some default container mounts used to. 

## Package Registry

A package registry is a centrallized, compiled list of software packages for pavement. All packages are organized under an organization and name. All packages have versions which are semver compliant.

### Local Cache

A package cache exists locally that caches packge lookups to avoid multiple round trips. 

Packages are cached based on the following criteria:

- package name
- package version
- os
- arch

By default pre-compiled package assets are dropped in the package cache upon use.

If a package cache miss occurs, the package will be build from source and the result assets will be dropped in the package cache.

## Deployments

A deployment is a manifest for deploying software across clusters. A deployment will contain the following assets:

- images
- packages
- vm types
- configuration

This is a sample deployment

```
resource "deployment" "rabbitmq" {

    name    = "rabbitmq"
    version = "3.7.19"

    set {
        name        = "rabbitmq"
        instances   = 3
        
        service {
            name        = "rabbitmq"
            image       = "ubuntu:18.04"
        }
        
        package {
            name    = "rabbitmq"
            version = "3.7.19"
        }
    }
}
```