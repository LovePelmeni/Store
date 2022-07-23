terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
    kubernetes {
      source = "hashicorp/kubernetes"
      version = ">=2.0.1"
    }
  }
}
// Google Cloud Credentials.... 
data "google_prodiver_credentials" "credentials" {

    googleProjectCredentials = "googleCredentials.json"
    googleProjectName = "google-project-id"
    googleProjectRegionName = "us-central-1" 

    googleProjectTimeZone = "us-central"
    googleProjectCloudBootDisk = "cloud-debian/debian-9"

    GoogleComputeInstanceMachineType = "f1.micro"
}

data "google_container_cluster" "kubernetes_cluster" {
  name = "cluster-name"
  location = "us-central1"
}

data "google_default_config" "default" { // contains generated access token for google cloud management,
//that will be active for specific period of time.
  access_token = "access-token"
}

data "terraform_remote_state" "gke" { // Parses google Cloud Project terraform state..
    backend = "local"
    config {
      path = "../config-path"
    }
}


data "kubernetes_manager_credentials" { // Credentials for Managing Cluster
  KubernetesUser = "kube-cluster-manager"  // username of the cluster role 
  KubernetesPassword = "kube-cluster-password" // password of the cluster role.
}

// Providers goes there... 

provider "google" {
  
  credentials = file("${data.google_provider_credentials.credentials.googleProjectCredentials}.json")
  project     = data.google_provider_credentials.credentials.googleProjectName
  region      = data.google_provider_credentials.credentials.googleRegionName 
  zone        = data.google_provider_credentials.credentials.googleTimeZone 
}

provider "kubernetes" {

  username = data.kubernetes_manager_credentials.KubernetesUser 
  password = data.kubernetes_manager_credentials.KubernetesPassword 

  host = data.kubernetes_cluster_credentials.credentials.ClusterHost 
  token = data.google_default_config.default.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.kubernetes_cluster.master_auth[0].cluster_ca_certificate)
}
