provider "kubernetes" {
    config_context_cluster = "minikube"
} 

// Runs 4 different namespaces to initialize abstract edges for the microservices, where they will be located at.
resource "kubernetes_namespace" "store-namespace" {
    name = "store-namespace"
    metadata {
        name = "store-namespace"
    }
}

resource "kubernetes_namespace" "email-namespace" {
    name = "email-namespace" 
    metadata {
      name = "email-namespace"
    }
} 

resource "kubernetes_namespace" "payment-namespace" {
    name = "payment-namespace" 
    metadata {
      name = "payment-namespace"
    }
}