terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
  }
}
variable "GOOGLE_CREDENTIALS_FILE" {
  type        = string
  description = "Filename, that provides credentials for your Google Cloud Project."
}

variable "GOOGLE_PROJECT_NAME" {
  type        = string
  description = "Google Cloud Project Name."
}

variable "GOOGLE_PROJECT_REGION_ZONE" {
  type        = string
  description = "Google Project Region Zone."
}

variable "GOOGLE_PROJECT_TIMEZONE" {
  type        = string
  description = "Google Project Timezone."
}

variable "GOOGLE_CLOUD_BOOT_DISK" {
  type        = string
  description = "Google Boot Disk for the Compute Instance."
}

variable "GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE" {
  type        = string
  description = "Google Comput Instance Machine Type, That it is going to be run on. Example - `f1-micro`"
}

provider "google" {

  version     = "3.5.0"
  credentials = file("${CREDENTIALS_FILE}.json")
  project     = GOOGLE_PROJECT_NAME
  region      = GOOGLE_PROJECT_REGION_ZONE
  zone        = GOOGLE_PROJECT_TIMEZONE
}

resource "google_compute_instance" "appserver" {

  name         = "store-service"
  machine_type = GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE

  boot_disk {
    initialize_params {
      image = GOOGLE_CLOUD_BOOT_DISK
    }
  }
  tags {
    Name = "Store Service Google Compute Instance"
  }
}

resource "google_compute_network" "vpc_network" {
  name = "ProjectVpcNetwork"
}