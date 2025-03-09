package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	pb "treeship/api/gen/v1"
	"treeship/kube"

	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type Manger struct {
	kubeCtl *kube.KubeClient
	stream  pb.AgentService_MessageRouteClient
	logger  *zap.Logger
}

func (a *Manger) ReadSteam(ctx context.Context) {
	for {
		req, err := a.stream.Recv()
		if err == io.EOF {
			a.logger.Info("stream closed")
			break
		}

		if err != nil {
			a.logger.Error("failed to receive message", zap.Error(err))
			break
		}

		a.logger.Info("received message", zap.String("agent_id", req.AgentId), zap.String("message", req.Message))

		data, err := KubeQuery(ctx, a.kubeCtl, req.Message, req.Namespace)
		if err != nil {
			a.logger.Error("failed to get data", zap.Error(err))

			err := a.Send(req.AgentId, fmt.Sprintf("failed to get data: %v", err))
			if err != nil {
				a.logger.Error("failed to send message", zap.Error(err))
			}
			continue
		}

		var convertedData map[string]interface{}
		if err := json.Unmarshal(data, &convertedData); err != nil {
			a.logger.Error("failed to unmarshal data", zap.Error(err))
			a.Send(req.AgentId, fmt.Sprintf("failed to unmarshal data: %v", err))
			continue
		}

		d, err := structpb.NewStruct(convertedData)
		if err != nil {
			a.logger.Error("failed to convert data to struct", zap.Error(err))
			a.Send(req.AgentId, fmt.Sprintf("failed to convert data to struct: %v", err))
			continue
		}

		err = a.stream.Send(&pb.MessageRequest{AgentId: req.AgentId, Responce: d})
		if err != nil {
			a.logger.Error("failed to send message", zap.Error(err))
			a.Send(req.AgentId, fmt.Sprintf("failed to send message: %v", err))
		}
	}
}

func (a *Manger) Send(agentID, message string) error {
	err := a.stream.Send(&pb.MessageRequest{AgentId: agentID, Message: message})
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}

	return nil
}

func New(ctx context.Context, cc grpc.ClientConnInterface, logger *zap.Logger, kubeCtl *kube.KubeClient) (*Manger, error) {
	client := pb.NewAgentServiceClient(cc)

	stream, err := client.MessageRoute(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create stream: %v", err)
	}

	return &Manger{
		kubeCtl: kubeCtl,
		logger:  logger,
		stream:  stream,
	}, nil
}
