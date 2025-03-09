package kube

import (
	"context"
	"fmt"
	"reflect"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	hrv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceType defines known Kubernetes resource types
type ResourceType string

// ResourceListFetcher defines a function type for fetching lists of resources
type ResourceListFetcher func(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error)

// Known resource types
const (
	ResourcePod            ResourceType = "pod"
	ResourceService        ResourceType = "service"
	ResourceServiceAccount ResourceType = "service_account"
	ResourceHelmRelease    ResourceType = "helm_releases"
	ResourceReplicaSet     ResourceType = "replicaset"
	ResourceStatefulSet    ResourceType = "statefulset"
	ResourceDeployment     ResourceType = "deployment"
	ResourceIngress        ResourceType = "ingress"
	ResourceExternalSecret ResourceType = "external_secrets"
	ResourceCronJob        ResourceType = "cronjob"
)

// GetResource fetches a specific Kubernetes resource by namespace, type and name
func GetResource(ctx context.Context, c *KubeClient, namespace, resource, name string) (ctrlClient.Object, error) {
	client, err := NewControllerClient(c.RestConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create controller client: %w", err)
	}

	obj, err := newResourceObject(ResourceType(resource))
	if err != nil {
		return nil, err
	}

	err = client.Get(ctx, ctrlClient.ObjectKey{Namespace: namespace, Name: name}, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s '%s' in namespace '%s': %w",
			resource, name, namespace, err)
	}

	removeFromObj(obj, "ManagedFields")
	return obj, nil
}

// newResourceObject returns an empty object of the requested resource type
func newResourceObject(resourceType ResourceType) (ctrlClient.Object, error) {
	switch resourceType {
	case ResourcePod:
		return &corev1.Pod{}, nil
	case ResourceService:
		return &corev1.Service{}, nil
	case ResourceServiceAccount:
		return &corev1.ServiceAccount{}, nil
	case ResourceHelmRelease:
		return &hrv2beta2.HelmRelease{}, nil
	case ResourceReplicaSet:
		return &appsv1.ReplicaSet{}, nil
	case ResourceStatefulSet:
		return &appsv1.StatefulSet{}, nil
	case ResourceDeployment:
		return &appsv1.Deployment{}, nil
	case ResourceIngress:
		return &networkingv1.Ingress{}, nil
	case ResourceExternalSecret:
		return &esv1beta1.ExternalSecret{}, nil
	case ResourceCronJob:
		return &batchv1.CronJob{}, nil
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// GetAllFromNamespace fetches all supported resources from a namespace
func GetAllFromNamespace(ctx context.Context, c *KubeClient, namespace string) (map[string][]interface{}, error) {
	client, err := NewControllerClient(c.RestConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime client: %w", err)
	}

	// Map of resource fetcher functions
	resourceFetchers := map[string]ResourceListFetcher{
		"pods":         fetchPods,
		"services":     fetchServices,
		"replicasets":  fetchReplicaSets,
		"statefulsets": fetchStatefulSets,
		"deployments":  fetchDeployments,
		"ingresses":    fetchIngresses,
	}

	items := make(map[string][]interface{})

	// Execute each fetcher and collect results
	for resourceType, fetcher := range resourceFetchers {
		_, objects, err := fetcher(ctx, client, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list %s: %w", resourceType, err)
		}

		for _, obj := range objects {
			// Process common fields
			removeFromObj(obj, "ManagedFields")

			// Add processed object to results
			items[resourceType] = append(items[resourceType], obj)
		}
	}

	return items, nil
}

// Resource list fetcher functions for each resource type
func fetchPods(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &corev1.PodList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

func fetchServices(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &corev1.ServiceList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

func fetchReplicaSets(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &appsv1.ReplicaSetList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		// Remove Spec from ReplicaSets for brevity
		removeFromObj(&list.Items[i], "Spec")
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

func fetchStatefulSets(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &appsv1.StatefulSetList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		// Remove Spec from StatefulSets for brevity
		removeFromObj(&list.Items[i], "Spec")
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

func fetchDeployments(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &appsv1.DeploymentList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		// Remove Spec from Deployments for brevity
		removeFromObj(&list.Items[i], "Spec")
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

func fetchIngresses(ctx context.Context, client ctrlClient.Client, namespace string) (interface{}, []interface{}, error) {
	list := &networkingv1.IngressList{}
	err := client.List(ctx, list, ctrlClient.InNamespace(namespace))
	if err != nil {
		return nil, nil, err
	}

	result := make([]interface{}, len(list.Items))
	for i := range list.Items {
		result[i] = &list.Items[i]
	}

	return list, result, nil
}

// removeFromObj safely removes a field from an object by setting it to its zero value
func removeFromObj(obj interface{}, fieldName string) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	v = v.Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.Zero(field.Type()))
	}
}
