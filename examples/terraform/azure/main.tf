provider "azurerm"{    
    version = "=1.27.1"
}

variable "username" {
    default = "pavement"
}

variable "password" {}

resource "azurerm_resource_group" "pavement" {
    name = "pavement"
    location = "East US 2"
}
