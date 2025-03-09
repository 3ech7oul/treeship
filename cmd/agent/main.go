package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"treeship/agent"
	"treeship/helper"
	"treeship/kube"

	"github.com/go-logr/zapr"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	logLevel       = pflag.StringP("log-level", "l", "info", "log level")
	kubeConfigPath = pflag.StringP("kubeconfig", "k", "", "path to kubeconfig file")
	agentId        = pflag.StringP("id", "i", "", "agent id")
	addr           = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	pflag.Parse()

	logger, err := helper.NewLogger(*logLevel)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	zaprLogger := zapr.NewLogger(logger)
	ctrllog.SetLogger(zaprLogger)

	fmt.Println("agentId", *agentId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := kube.RestConfig(*kubeConfigPath)
	if err != nil {
		logger.Fatal("failed to get kubeconfig", zap.Error(err))
	}

	client := kube.NewKubeClient(config)

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("did not connect", zap.Error(err))
	}
	defer conn.Close()

	manger, err := agent.New(ctx, conn, logger, client)
	if err != nil {
		logger.Fatal("could not create message route", zap.Error(err))
	}

	manger.Send(*agentId, "Hello from agent")
	manger.ReadSteam(ctx)
}
