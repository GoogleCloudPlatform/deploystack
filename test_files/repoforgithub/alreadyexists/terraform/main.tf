/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */


# Enabling services in your GCP project


data "google_project" "project" {
  project_id = var.project_id
}

resource "google_project_service" "all" {
  for_each                   = toset(var.gcp_service_list)
  project                    = data.google_project.project.number
  service                    = each.key
  disable_dependent_services = false
  disable_on_destroy         = false
}


data "google_compute_network" "default" {
  project    = var.project_id
  name       = "default"
  depends_on = [google_project_service.all]
}

resource "google_compute_network" "main" {
  provider                = google-beta
  name                    = "${var.basename}-network"
  auto_create_subnetworks = true
  project                 = var.project_id
  depends_on              = [google_project_service.all]
}

resource "google_compute_firewall" "default-allow-http" {
  name    = "deploystack-allow-http"
  project = data.google_project.project.number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = ["http-server"]
  depends_on  = [google_project_service.all]
}

resource "google_compute_firewall" "default-allow-internal" {
  name    = "deploystack-allow-internal"
  project = data.google_project.project.number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["10.128.0.0/20"]
  depends_on    = [google_project_service.all]

}

resource "google_compute_firewall" "default-allow-ssh" {
  name    = "deploystack-allow-ssh"
  project = data.google_project.project.number
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = ["ssh-server"]
  depends_on  = [google_project_service.all]
}

# Create Instances
resource "google_compute_instance" "server" {
  name                      = "server"
  zone                      = var.zone
  project                   = var.project_id
  machine_type              = "e2-standard-2"
  tags                      = ["ssh-server", "http-server"]
  allow_stopping_for_update = true


  boot_disk {
    auto_delete = true
    device_name = "server"
    initialize_params {
      image = "family/ubuntu-1804-lts"
      size  = 10
      type  = "pd-standard"
    }
  }

  network_interface {
    network = google_compute_network.main.name
    access_config {
      // Ephemeral public IP
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/logging.write"]
  }

  metadata_startup_script = <<SCRIPT
    apt-get update
    apt-get install -y mongodb
    service mongodb stop
    sed -i 's/bind_ip = 127.0.0.1/bind_ip = 0.0.0.0/' /etc/mongodb.conf
    iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 27017
    service mongodb start
  SCRIPT
  depends_on              = [google_project_service.all]
}

resource "google_compute_instance" "client" {
  name                      = "client"
  zone                      = var.zone
  project                   = var.project_id
  machine_type              = "e2-standard-2"
  tags                      = ["http-server", "https-server", "ssh-server"]
  allow_stopping_for_update = true

  boot_disk {
    auto_delete = true
    device_name = "client"
    initialize_params {
      image = "family/ubuntu-1804-lts"
      size  = 10
      type  = "pd-standard"
    }
  }
  service_account {
    scopes = ["https://www.googleapis.com/auth/logging.write"]
  }

  network_interface {
    network = google_compute_network.main.name

    access_config {
      // Ephemeral public IP
    }
  }

  metadata_startup_script = <<SCRIPT
    add-apt-repository ppa:longsleep/golang-backports -y && \
    apt update -y && \
    apt install golang-go -y
    mkdir /modcache
    mkdir /go
    mkdir /app && cd /app
    curl https://raw.githubusercontent.com/GoogleCloudPlatform/golang-samples/main/compute/quickstart/compute_quickstart_sample.go --output main.go
    go mod init exec
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache go mod tidy
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache go get -u 
    sed -i 's/mongoport = "80"/mongoport = "27017"/' /app/main.go
    echo "GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache HOST=${google_compute_instance.server.network_interface.0.network_ip} go run main.go"
    GOPATH=/go GOMODCACHE=/modcache GOCACHE=/modcache HOST=${google_compute_instance.server.network_interface.0.network_ip} go run main.go & 
  SCRIPT

  depends_on = [google_project_service.all]
}
