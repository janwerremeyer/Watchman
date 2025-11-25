package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"watchman/config"
	"watchman/internal/aws"
	"watchman/internal/container"
	"watchman/internal/gofy"
	"watchman/internal/orchestrator"
)

func main() {
	//Load Configs
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	//Initialize AWS - We need it to check tags and pull images
	awsClient, err := aws.NewClient(cfg.AWS)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	docker, err := container.NewDockerClient()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	orchestratorClient := orchestrator.NewOrchestrator(docker, awsClient, cfg.Wanted)

	scheduler := gofy.NewScheduler()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	scheduler.Start(10*time.Second, orchestratorClient.Run)

	<-stop
	scheduler.Stop()

	//tags, err := aws.ListImages(awsCfg, "short-rib-admin")
	//if err != nil {
	//	slog.Error(err.Error())
	//}
	//
	//slog.Info("Image Tags", "Tags", tags)
	//
	//err = aws.PullECRImageWithDockerClient(awsCfg, "short-rib-admin", tags[0])
	//if err != nil {
	//	slog.Error(err.Error())
	//}
}
