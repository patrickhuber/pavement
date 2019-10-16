variable "username" {
    type = "string"
}

variable "password" {
    type = "string"    
}

resource "google_compute_instance" "pavement" {

    name="pavement"
    machine_type = "n1-standard-1"    

    boot_disk {
        initialize_params {
            image = "gce-uefi-images/ubuntu-1804-lts"
        }
    }

    scratch_disk{
    }

    network_interface {
        subnetwork = "${google_compute_subnetwork.pavement.name}"

        // assigns a dynamic public ip address
        access_config{
        }
    }

    metadata_startup_script="useradd ${var.username}; chpasswd <<<\"${var.username}:${var.password}\""
}

resource "google_compute_firewall" "pavement" {
    name = "pavement"
    network = "${google_compute_network.pavement.name}"

    allow{
        protocol = "icmp"
    }

    allow {
        protocol = "tcp"
        ports = ["22"]        
    }
}