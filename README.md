# Pavement

Pavement provides a cli to scaffold and manage generic infrastructure across multiple clouds and private clouds. 

## Compute

Compute resources are created with a Dockerfile like syntax called Pavementfile. 

### Images

Pavement images map cloud resources to pavement image name and version.

### azure

```
# creates an azure image map mapping an azure base image to a pavement image and alias

# `FROM ubuntu:18.04`
image_type "ubuntu" "18_04" {

    provider = "azure"

    properties = {
        "publisher"   = "Canonical"
        "offer"       = "UbuntuServer"
        "sku"         = "18.04-LTS"
        "version"     = "18.04"
    }
}

# `FROM ubuntu:16.04`
image_type "ubuntu" "16_04" {

    provider = "azure" 

    properties = {
        "publisher"   = "Canonical"
        "offer"       = "UbuntuServer"
        "sku"         = "16.04-LTS"
        "version"     = "16.04"
    }
}

image "ubuntu" "16_04" {
    image_type_id = image_type.ubuntu.16_04
}

image "ubuntu" "18_04" {
    image_type_id = image_type.ubuntu.18_04
}
```

### aws

```
# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:18.04`
image_type "ubuntu" "18_04"{

    provider = "aws"

    # aws specific details
    properties = {
        "image_name"          = "ubuntu/images/hvm-ssd/ubuntu-xenial-18.04-amd64-server-*"
        "virtualization_type" = "hvm"
        "owner"               = "099720109477" # canonical
    }    
}

# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:16.04`
image_type "ubuntu" "16_04"{

    provider = "aws"

    # aws specific details
    properties = {
        "image_name"          = "ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*"
        "virtualization_type" = "hvm"
        "owner"               = "099720109477" # canonical
    }    
}

image "ubuntu" "16_04" {
    image_type_id = image_type.ubuntu.16_04
}

image "ubuntu" "18_04" {
    image_type_id = image_type.ubuntu.18_04
}
```

### gcp

```
# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:18.04`
image_type "ubuntu" "18_04" {

    provider = "gcp"

    # gcp specific details
    properties = {
        # "image_name"  = "" # mutually exclusive with `family`
        "family"      = "ubuntu-1804-lts" # mutually exclusive with `image_name`
        "project"     = "ubuntu-os-cloud"
    }    
}

# creates an aws image map mapping an aws base image to a pavement image and alias
# `FROM ubuntu:16.04`
image_type "ubuntu" "16_04" {

    provider = "gcp"

    # gcp specific details
    properties = {
        # "image_name"  = "" # mutually exclusive with `family`
        "family"      = "ubuntu-1604-lts" # mutually exclusive with `image_name`
        "project"     = "ubuntu-os-cloud"
    }    
}

image "ubuntu" "16_04" {
    image_type_id = image_type.ubuntu.16_04
}

image "ubuntu" "18_04" {
    image_type_id = image_type.ubuntu.18_04
}
```

### vsphere

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

## Virtual Machine Types

Pavement virtual machine types map native virtual machine properties to reusable pavement virtual machine types. You can then create instances of the virtual machines using the given types. This allows you to define IaaS specific types but abstract virtual machines.

### azure

```
virtual_machine_type "default" {

    provider    = "azure"
    is_default  = true

    properties = {
        "instance_type" = "Standard_D1_v2"
    }
}

virtual_machine "web" {
    image_id    = image.ubuntu.18_04.id
    vm_type_id  = vm_type.default.id
}
```

### aws

```
virtual_machine_type "default" {

    provider    = "aws"
    is_default  = true

    properties = {
        "instance_type"       = "m4.large"
        "ephemeral_disk_size" = "25GB"
    }
}

virtual_machine "web" {
    image_id    = image.ubuntu.18_04.id
    vm_type_id  = vm_type.default.id
}
```

### gcp

```
virtual_machine_type "default" {

    provider    = "gcp"
    is_default  = true

    properties = {
        "machine_type"    = "n1-standard-2"
        "root_disk_size"  = "50GB"
        "root_disk_type"  = "pd-ssd"
    }
}

virtual_machine "web" {
    image_id    = image.ubuntu.18_04.id
    vm_type_id  = vm_type.default.id
}
```

### vsphere

```
virtual_machine_type "default" {

    provider    = "vsphere"
    is_default  = true

    properties = {
        "name"    = "default"
        "cpu"     = "2"
        "ram"     = 1GB
        "disk"    = 30GB
    }
}

virtual_machine "web" {
    image_type_id   = image_type.ubuntu.18_04.id
    vm_type_id      = vm_type.default.id
}
```