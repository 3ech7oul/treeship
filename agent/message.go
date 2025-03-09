package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"treeship/kube"
)

func KubeQuery(ctx context.Context, kubeCtl *kube.KubeClient, msgType, namespace string) ([]byte, error) {
	switch msgType {
	case "GetAllFromNamespace":
		data, err := kube.GetAllFromNamespace(ctx, kubeCtl, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get all resources from namespace %s: %v", namespace, err)
		}

		response, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
		}

		return response, nil
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", msgType)
	}
}
