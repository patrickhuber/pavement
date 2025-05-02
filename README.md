# Pavement

Pavement provides a cli to scaffold and manage generic infrastructure across multiple clouds and private clouds. 

# Getting Started

## Create Pavement File

> Pavementfile

```
FROM unbuntu:latest
```

## Create resource definitions

> pavement.hcl

```
image file_image {
    file = "./Pavementfile"
}

machine "demo" {
    image = image.file_image    
}
```

## Run Pavement Apply

```

```

```bash
pavement apply .
```

# Cloud Providers

# Resources

Pavement uses generic resources to create images on a specific cloud provider