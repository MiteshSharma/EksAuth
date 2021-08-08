package main

import (
	"fmt"
	"io/ioutil"

	"github.com/MiteshSharma/EksAuth/k8s"
)

var certificateAuthorityDataFile = "cert file path"
var clusterName = ""
var clusterServerUrl = ""
var awsKeyId = ""
var awsSecretKey = ""
var awsRegion = ""
var clusterId = ""

func main() {
	certificateAuthorityData, err := ioutil.ReadFile(certificateAuthorityDataFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cluster := &k8s.EksCluster{
		Name:                     clusterName,
		Server:                   clusterServerUrl,
		CertificateAuthorityData: []byte(certificateAuthorityData),
		AwsKeyId:                 awsKeyId,
		AwsSecretKey:             awsSecretKey,
		AwsRegion:                awsRegion,
		ClusterID:                clusterId,
	}
	fmt.Println("Testing eks cluster auth by fetching all namespaces")
	cluster.TestCluster()
	fmt.Println("Testing eks cluster completed")
}
