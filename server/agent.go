package server

import (
	"fmt"
	"sync"
	pb "treeship/api/gen/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentManager struct {
	mu      sync.RWMutex
	logger  *zap.Logger
	streams map[string]*pb.AgentService_MessageRouteServer
	pb.UnimplementedAgentServiceServer
}

func (r *AgentManager) AgentsConnnected() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.streams)
}

func (r *AgentManager) AddAgentStream(agentID string, stream *pb.AgentService_MessageRouteServer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.streams[agentID] = stream
	r.logger.Info("added agent stream", zap.String("agent_id", agentID))
}

func (r *AgentManager) RemoveAgentStream(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.streams, agentID)
	r.logger.Info("removed agent stream", zap.String("agent_id", agentID))
}

func (r *AgentManager) SendMessage(agentID string, message string) (*pb.MessageRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stream, exists := r.streams[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not found", agentID)
	}

	err := (*stream).Send(&pb.MessageResponse{AgentId: agentID, Message: message})
	if err != nil {
		return nil, fmt.Errorf("could not send message: %v", err)
	}

	msReg, err := (*stream).Recv()
	if err != nil {
		return nil, fmt.Errorf("could not receive message: %v", err)
	}

	return msReg, nil
}

func (r *AgentManager) MessageRoute(stream pb.AgentService_MessageRouteServer) error {
	ctx := stream.Context()
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to receive initial message: %v", err)
	}

	agentID := req.AgentId

	r.AddAgentStream(agentID, &stream)
	defer r.RemoveAgentStream(agentID)

	for {
		select {
		case <-ctx.Done():
			return status.Error(codes.Canceled, "client disconnected")

			//default:
			/* req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					r.logger.Info("agent stream ended normally",
						zap.String("agent_id", agentID))
					return nil
				}

				// Check for context cancellation in receive errors
				if status.Code(err) == codes.Canceled {
					r.logger.Info("agent disconnected during receive",
						zap.String("agent_id", agentID))
					return status.Error(codes.Canceled, "client disconnected")
				}

				r.logger.Error("error receiving from agent",
					zap.String("agent_id", agentID),
					zap.Error(err))
				return status.Errorf(codes.Unknown, "failed to receive: %v", err)
			} */

		}
	}
}

func NewAgentManager(logger *zap.Logger) *AgentManager {
	return &AgentManager{
		streams: make(map[string]*pb.AgentService_MessageRouteServer),
		logger:  logger,
	}
}
