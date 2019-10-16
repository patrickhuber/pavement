resource "google_compute_network" "pavement" {
    name = "pavement-network"
}

resource "google_compute_subnetwork" "pavement" {
    name = "pavement-network"    
    ip_cidr_range = "10.0.0.0/24"
    network = "${google_compute_network.pavement.id}"
}