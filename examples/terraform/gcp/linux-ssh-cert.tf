resource "google_compute_instance" "default" {
  name         = "test"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    ssh-keys = "${var.username}:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0QAW7FdVV8fmcwm9Wmnfay9SQbnmElpWrhEbpdPJuboJYE5j6C1PUcAwYOmswAATb4idT8bdl4EvTME4hpiSZZVzblkfojVX9FiNh7Q+tTPjUyBstBOZk6vZ/IVB+wLi2StYFDj4plqAJEAinpAqCMT42dYwRzDBve30jsveH7DPQjRWrQi/mzsihsVkQHxL66vA5zfRJl3gMqJ0BTB1pbEqy4gxLDdxu+OsV+r5s80FU1ASiVPR5REkIdHiFn/ZxbMYIco3UGTWl9xBlD17f+rrWc7voxciqluNBd8BlD6m5QdSY3IkJVGqRmvV6Uu2rZOfk2HRrnr0FoUzDMHxR"
  }
}
