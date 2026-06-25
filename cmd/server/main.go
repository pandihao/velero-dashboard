package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"velero-api-server/internal/handler"
	"velero-api-server/pkg/cluster"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	// Initialize controller-runtime logger
	ctrllog.SetLogger(zap.New(zap.UseDevMode(false)))

	var (
		port            int
		kubeconfigDir   string
		kubeconfig      string
		namespace       string
		insecureSkipTLS bool
	)

	flag.IntVar(&port, "port", 8080, "HTTP server listen port")
	flag.StringVar(&kubeconfigDir, "kubeconfig-dir", ".kube", "Directory containing kubeconfig files for multiple clusters")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to single kubeconfig file (overrides kubeconfig-dir)")
	flag.StringVar(&namespace, "namespace", "velero", "Velero namespace")
	flag.BoolVar(&insecureSkipTLS, "insecure-skip-tls", false, "Skip TLS certificate verification")
	flag.Parse()

	if env := os.Getenv("VELERO_NAMESPACE"); env != "" {
		namespace = env
	}
	if env := os.Getenv("KUBECONFIG"); env != "" && kubeconfig == "" {
		kubeconfig = env
	}
	if env := os.Getenv("KUBECONFIG_DIR"); env != "" && kubeconfigDir == ".kube" {
		kubeconfigDir = env
	}

	// Create cluster manager
	clusterMgr := cluster.NewManager()

	// Load clusters
	if kubeconfig != "" {
		// Single cluster mode
		log.Printf("Loading single cluster from kubeconfig: %s", kubeconfig)
		if err := clusterMgr.AddCluster("default", kubeconfig, insecureSkipTLS); err != nil {
			log.Fatalf("Failed to add cluster: %v", err)
		}
	} else {
		// Multi-cluster mode
		log.Printf("Loading clusters from directory: %s", kubeconfigDir)
		if err := clusterMgr.LoadFromDirectory(kubeconfigDir, insecureSkipTLS); err != nil {
			log.Fatalf("Failed to load clusters: %v", err)
		}
	}

	clusters := clusterMgr.ListClusters()
	log.Printf("Loaded %d cluster(s): %v", len(clusters), clusters)

	r := gin.Default()

	// Serve Swagger documentation
	r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	r.StaticFile("/swagger", "./docs/swagger.html")
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger")
	})

	handler.RegisterRoutes(r, clusterMgr, namespace)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting velero-api-server on %s (namespace: %s)", addr, namespace)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
