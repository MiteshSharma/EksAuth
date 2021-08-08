package k8s

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	token "sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EksCluster struct {
	Name                     string
	Server                   string
	CertificateAuthorityData []byte
	ClusterID                string
	AwsKeyId                 string
	AwsSecretKey             string
	AwsRegion                string
}

func (cluster *EksCluster) TestCluster() {
	restConf, err := cluster.GetRESTConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	clientset, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	pods, err := clientset.CoreV1().Namespaces().List(
		context.TODO(),
		metav1.ListOptions{},
	)
	fmt.Println(pods)
	fmt.Println(err)
}

func (cluster *EksCluster) GetDynamicClient() (dynamic.Interface, error) {
	fmt.Println("GetDynamicClient start")
	restConf, err := cluster.GetRESTConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client, err := dynamic.NewForConfig(restConf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("GetDynamicClient end")

	return client, nil
}

func (cluster *EksCluster) GetRESTConfig() (*rest.Config, error) {
	fmt.Println("GetRESTConfig start")
	cmdConf, err := cluster.GetClientConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	restConf, err := cmdConf.ClientConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rest.SetKubernetesDefaults(restConf)
	fmt.Println("GetRESTConfig end")

	return restConf, nil
}

func (cluster *EksCluster) GetClientConfig() (clientcmd.ClientConfig, error) {
	fmt.Println("GetClientConfig start")
	apiConfig, err := cluster.CreateRawConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	overrides := &clientcmd.ConfigOverrides{}
	overrides.Context = api.Context{
		Namespace: "default",
	}

	config := clientcmd.NewDefaultClientConfig(*apiConfig, overrides)

	fmt.Println("GetClientConfig end")

	return config, nil
}

func (cluster *EksCluster) CreateRawConfig() (*api.Config, error) {
	fmt.Println("CreateRawConfig start")
	apiConfig := &api.Config{}

	clusterMap := make(map[string]*api.Cluster)

	clusterMap[cluster.Name] = &api.Cluster{
		Server:                   cluster.Server,
		InsecureSkipTLSVerify:    false,
		CertificateAuthorityData: cluster.CertificateAuthorityData,
	}

	awsSession, err := cluster.getAwsSession()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	clusterToken, err := cluster.getClusterToken(awsSession)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	authInfoName := cluster.Name
	authInfoMap := make(map[string]*api.AuthInfo)
	authInfo := &api.AuthInfo{}
	authInfo.Token = clusterToken
	authInfoMap[authInfoName] = authInfo

	contextMap := make(map[string]*api.Context)

	contextMap[cluster.Name] = &api.Context{
		Cluster:  cluster.Name,
		AuthInfo: authInfoName,
	}

	apiConfig.Clusters = clusterMap
	apiConfig.AuthInfos = authInfoMap
	apiConfig.Contexts = contextMap
	apiConfig.CurrentContext = cluster.Name

	fmt.Println("CreateRawConfigFromCluster end")
	return apiConfig, nil
}

func (cluster *EksCluster) getAwsSession() (*session.Session, error) {
	awsConf := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			cluster.AwsKeyId,
			cluster.AwsSecretKey,
			"",
		),
	}
	awsConf.Region = aws.String(cluster.AwsRegion)

	awsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *awsConf,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return awsSession, nil
}

func (cluster *EksCluster) getClusterToken(awsSession *session.Session) (string, error) {
	generator, err := token.NewGenerator(false, false)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	token, err := generator.GetWithOptions(&token.GetTokenOptions{
		Session:   awsSession,
		ClusterID: cluster.ClusterID,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return token.Token, nil
}
