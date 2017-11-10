package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubeClient() (*kubernetes.Clientset, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	home := homeDir()
	if home == "" {
		return nil, fmt.Errorf("Home dir not found")
	}

	path := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", path)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
