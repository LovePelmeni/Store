// link to the manifest of the project that has been versioned...
env.DEPLOYMENT_MANIFEST_LINK="https://github.com/<YOUR-NICKNAME>/<YOUR-REPOSITORY>/<YOUR-KUBERNETES-WHATEVER-YOU-DEPLOY-MANIFEST-FILE-PATH>"
 // link to the file 
// example: if you changed some stuff related to PostgreSQL and decided to pull it into production.

// after you created new Docker Image out of that,
// you need to specify full path to the `Kubernetes Yaml File` where your `PostgreSQL Docker Image` is going to be used,

env.MANIFEST_NAMESPACE=""