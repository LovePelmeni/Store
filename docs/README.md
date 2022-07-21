# *Store Docs* 

#### API Documentation ~ [Swagger UI Documentation]("http:localhost:8000")

--- 

## *Introduction* 

Store Microservice Is a Core Engine of the Project, that is Responsible for managing the Whole Application itself, it has access to

`Email Microservice` ~ Repo : "https://github.com/LovePelmeni/EmailService.git" 
`Payment Microservice` ~ Repo : "https://github.com/LovePelmeni/Payment-Service.git" 

Also Communicating via `Firebase Real Time Database` with

`Order Microservice` ~ Repo : "https://github.com/LovePelmeni/OrderService.git"

--- 

## *Usage*

Clone this Repo 
```
    $ git clone https://github.com/LovePelmeni/Store.git
```

--- 

## *Deploy* 


## Settings Up Databases.


In this Stage the requirements should be: 

`PostgreSQL Database` - `13.3` 

`Firebase Project with Real Time Database` 


### *Using Kubernetes*

1. Go and Setup Volumes and Claims + Secrets for Kubernetes `PostgreSQL` Manifest at `root/kubernetes/volumes/` and `root/kubernetes/config_maps` according to your requirements.

2. Go to the ./kubernetes directory (if you are not there) Deploy the Kubernetes Volume & VolumeClaim & ConfigMap Using `kubectl apply -f ./config_maps/postgresMaps.yaml && kubectl apply -f ./volumes/postgresClaim.yaml && kubectl apply -f ./volumes/postgresVolume.yaml` 

After that your Postgres database for the microservice should be ready.. Go and Check it out by running... 

```
    $ kubectl get statefulset store-postgres-database --namespace=store-namespace 
```
### *Using Docker*


Go to the `./docker-compose.yaml` file at root directory and Change the Environment Variables for `postgres-store-database` service, according to your requirements.


Then, go to the `project_env.env` file and replace the default `postgres` envinronment variables with that, you specified in the previous step. 


Also you might need to change the `firebase` env vars, depends on the project that you have.

--- 

## Building Stage For the Project.

To get the Initial Project Image `without Database and Integrated Services`
Go to the Root Directory and execute 
```
    $ docker build . --name=store_application
```


### *Using Docker-Compose*

To run the Full Microservice Using Docker-compose.yaml 

```
    $ docker-compose up -d 
```


### *Using Kubernetes*


Depends on where you are going to run this application, 
you need to check one of this Kubernetes Docs for this Project 

~ "https://github.com/LovePelmeni/Store/docs/kubernetes_docs/GCLOUD.md" ~ *Deploying `Store Microservice` On Google Cloid*

~ "https://github.com/LovePelmeni/Store/docs/kubernetes_docs/LOCAL_ENV.md" ~ *Deploying `Store Microservice` On Local Dev Machine with Minikube*

--- 

## Jenkins Integration 

Go Checkout Documentation Guide for Jenkings Continious Delivery / Integration Pipeline.

~ "https://github.com/LovePelmeni/Store/docs/jenkins/LOCAL_ENV.md" ~ *Setting Up Jenkins Pipeline for `Store MicroService` in Docker in Local Environment*

~ "https://github.com/LovePelmeni/Store/docs/jenkings/GCLOUD.md" ~ *Setting Up Jenkins Pipeline for `Store MicroService` on the Google Cloud using Kubernetes*

--- 

## *External Links* 

#### *Email For Contributions* ~ `kirklimushin@gmail.com`