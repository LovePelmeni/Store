
variable CREDENTIALS_FILE {
    type = "string"
    description = "Filename, that provides credentials for your Google Cloud Project."
    validation {
        error_message = "Invalid Credentials File"
    }
}

variable GOOGLE_PROJECT_NAME {
    type = "string"
    description = "Google Cloud Project Name."
}

variable GOOGLE_PROJECT_TIMEZONE {
    type = "string"
    description = "Google Project Timezone."
}

variable GOOGLE_BOOT_DISK {
    type = "string"
    description = ""
    validation {
        error_message = "Invalid Boot Disk."
    }
}

variable GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE {
    type = "string"
    description = "Google Comput Instance Machine Type, That it is going to be run on. Example - `f1-micro`"
    validation {
        error_message = "Invalid Compute Instance Machine Type"
        output "" {
            value = "f1-micro"
        }
    }
}

variable GOOGLE_CLOUD_BOOT_DISK {
    type = "string"
    description = "Google Cloud Boot Disk Image."
    validation {

    }
}

provider "google" {

    version = "3.5.0"
    credentials = "${CREDENTIALS_FILE}.json"
    project = "${GOOGLE_PROJECT_NAME}"
    zone = "${GOOGLE_PROJECT_TIMEZONE}"
    
}

resource "google_compute_instance" "appserver" {

    name = "store-service" 
    machine_type = "${GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE}"

    boot_disk {
      initialize_params {
        image = "${GOOGLE_CLOUD_BOOT_DISK}"
      }
    }
}