package kube

import (
	hrv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	scv1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/node/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = runtime.NewScheme()

func init() {
	runtimeutil.Must(appsv1.AddToScheme(scheme))
	runtimeutil.Must(corev1.AddToScheme(scheme))
	runtimeutil.Must(batchv1.AddToScheme(scheme))
	runtimeutil.Must(networkingv1.AddToScheme(scheme))
	runtimeutil.Must(scv1beta2.AddToScheme(scheme))
	runtimeutil.Must(hrv2beta2.AddToScheme(scheme))
	runtimeutil.Must(v1alpha1.AddToScheme(scheme))
}

type KubeClient struct {
	conf *rest.Config
}

func (c *KubeClient) RestConfig() *rest.Config {
	return c.conf
}

func NewKubeClient(conf *rest.Config) *KubeClient {
	return &KubeClient{conf: conf}
}

func RestConfig(kubeConfigPath string) (*rest.Config, error) {
	return getConfig(kubeConfigPath)
}

func NewControllerClient(conf *rest.Config) (ctrlClient.Client, error) {
	return ctrlClient.New(conf, ctrlClient.Options{Scheme: scheme})
}

func getConfig(kubeConfigPath string) (*rest.Config, error) {
	if kubeConfigPath == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		nil).ClientConfig()
}
