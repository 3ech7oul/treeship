package server

import (
	"fmt"
	"sync"
	pb "treeship/api/gen/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentRegistry struct {
	mu      sync.RWMutex
	logger  *zap.Logger
	streams map[string]*pb.AgentService_MessageRouteServer
	pb.UnimplementedAgentServiceServer
}

func (r *AgentRegistry) AgentsConnnected() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.streams)
}

func (r *AgentRegistry) AddAgentStream(agentID string, stream *pb.AgentService_MessageRouteServer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.streams[agentID] = stream
	r.logger.Info("added agent stream", zap.String("agent_id", agentID))
}

func (r *AgentRegistry) RemoveAgentStream(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.streams, agentID)
	r.logger.Info("removed agent stream", zap.String("agent_id", agentID))
}

func (r *AgentRegistry) SendMessage(agentID string, message string) (*pb.MessageRequest, error) {
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

func (r *AgentRegistry) MessageRoute(stream pb.AgentService_MessageRouteServer) error {
	ctx := stream.Context()
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to receive initial message: %v", err)
	}

	agentID := req.AgentId

	r.AddAgentStream(agentID, &stream)
	defer r.RemoveAgentStream(agentID)

	<-ctx.Done()
	return status.Error(codes.Canceled, "client disconnected")
}

func NewAgentRegistry(logger *zap.Logger) *AgentRegistry {
	return &AgentRegistry{
		streams: make(map[string]*pb.AgentService_MessageRouteServer),
		logger:  logger,
	}
}
