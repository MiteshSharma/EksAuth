# EksAuth

In this project, we are making API request to EKS cluster with information provided in kubeconfig file and AWS credentials using aws-iam-authenticator.

Needed variables:
1. certificateAuthorityDataFile: This file contains certificate data which is stored as base64 in kubeconfig file as variable certificate-authority-data. 
2. clusterName: Name of cluster 
3. clusterServerUrl: Cluster server url to make request to server
4. AWS details: Need aws access key, secret and region to authenticate with help of aws-iam-authenticator to authenticate with EKS cluster
5. clusterId: This is unique cluster identifier. Detail: https://github.com/kubernetes-sigs/aws-iam-authenticator#what-is-a-cluster-id

Once all information is updated in main.go, command to run main.go file:

```
go run main.go
```
