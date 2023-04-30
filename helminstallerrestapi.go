package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Configuration for Kubernetes cluster
type ClusterConfig struct {
	Kubeconfig string `json:"kubeconfig"`
	Context    string `json:"context"`
	Namespace  string `json:"namespace"`
}

// Configuration for Helm chart installation
type InstallConfig struct {
	Chart      string            `json:"chart"`
	Release    string            `json:"release"`
	Version    string            `json:"version"`
	Values     map[string]string `json:"values"`
	Wait       bool              `json:"wait"`
	Timeout    int               `json:"timeout"`
	KubeConfig ClusterConfig     `json:"kube_config"`
}

func main() {
	r := gin.Default()

	// Route for installing Helm chart
	r.POST("/install", func(c *gin.Context) {
		var installConfig InstallConfig
		if err := c.ShouldBindJSON(&installConfig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create Kubernetes REST config
		var kubeConfig *rest.Config
		if installConfig.KubeConfig.Kubeconfig != "" {
			// Use provided kubeconfig file
			config, err := clientcmd.BuildConfigFromFlags("", installConfig.KubeConfig.Kubeconfig)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			kubeConfig = config
		} else {
			// Use default kubeconfig
			config, err := rest.InClusterConfig()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			kubeConfig = config
		}

		// Set Kubernetes context and namespace
		if installConfig.KubeConfig.Context != "" {
			kubeConfig.CurrentContext = installConfig.KubeConfig.Context
		}
		if installConfig.KubeConfig.Namespace != "" {
			kubeConfig.Namespace = installConfig.KubeConfig.Namespace
		}

		// Create Helm action configuration
		actionConfig := new(action.Configuration)
		if err := actionConfig.Init(kubeConfig, installConfig.KubeConfig.Namespace, "secrets", log.Printf); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Install Helm chart
		client := action.NewInstall(actionConfig)
		client.ReleaseName = installConfig.Release
		client.Namespace = installConfig.KubeConfig.Namespace
		client.Timeout = installConfig.Timeout
		client.Wait = installConfig.Wait
		if installConfig.Version != "" {
			client.Chart.Version = installConfig.Version
		}
		if err := client.Run(installConfig.Chart, installConfig.Values); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Helm chart %s installed successfully", installConfig.Chart)})
	})

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
