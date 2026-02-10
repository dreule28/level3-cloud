package kube

import (
	"os"

	"github.com/dreule28/Week_4/paas-api/internal/config"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	K8sClient client.Client
}

func buildRestConfig() (*rest.Config, error) {
	//when deployed
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	//local kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func	NewClient(cfg config.Config) (*Client, error) {
	restCfg, err := buildRestConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()

	//get basic k8s types
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	//get CRDs
	utilruntime.Must(cnpgv1.AddToScheme(scheme))


	k8scllient, err := client.New(restCfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return &Client{K8sClient: k8scllient}, nil
}
