package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dragon3/pdb-monitor-controller/config"
	"github.com/dragon3/pdb-monitor-controller/controller"
	"github.com/dragon3/pdb-monitor-controller/log"
)

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run())
}

func run() int {
	config, err := config.NewFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize config from env vars: %s\n", err)
		return exitError
	}

	logger, err := log.New(config.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %s\n", err)
		return exitError
	}
	defer logger.Sync()

	var kubeConfig *rest.Config
	if config.KubeConfigPath != "" {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to initialize kube config: %s\n", err)
			return exitError
		}

	} else {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to initialize in-cluster config: %s\n", err)
			return exitError
		}
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize clientSet: %s\n", err)
		return exitError
	}

	controller, err := controller.New(
		logger,
		clientSet,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize controller: %s\n", err)
		return exitError
	}

	stopChan := make(chan struct{}, 1)
	go handleSigterm(stopChan)

	controller.Run(stopChan)

	return exitOK
}

func handleSigterm(stopChan chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	<-sigCh
	close(stopChan)
}
