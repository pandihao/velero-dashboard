package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterConfig holds a single cluster configuration
type ClusterConfig struct {
	Name       string
	KubeConfig string
	Client     client.Client
}

// Manager manages multiple Kubernetes clusters
type Manager struct {
	clusters map[string]*ClusterConfig
	mu       sync.RWMutex
}

// NewManager creates a new cluster manager
func NewManager() *Manager {
	return &Manager{
		clusters: make(map[string]*ClusterConfig),
	}
}

// LoadFromDirectory loads all kubeconfig files from a directory
func (m *Manager) LoadFromDirectory(dir string, insecureSkipTLS bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Skip hidden files and non-config files
		name := entry.Name()
		if name[0] == '.' {
			continue
		}

		configPath := filepath.Join(dir, name)
		clusterName := name // Use filename as cluster name

		if err := m.AddCluster(clusterName, configPath, insecureSkipTLS); err != nil {
			return fmt.Errorf("failed to add cluster %s: %v", clusterName, err)
		}
	}

	if len(m.clusters) == 0 {
		return fmt.Errorf("no valid kubeconfig files found in %s", dir)
	}

	return nil
}

// AddCluster adds a single cluster configuration
func (m *Manager) AddCluster(name, kubeconfigPath string, insecureSkipTLS bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create Kubernetes client
	cfg, err := m.buildConfig(kubeconfigPath)
	if err != nil {
		return err
	}

	if insecureSkipTLS {
		cfg.TLSClientConfig.Insecure = true
	}

	scheme := runtime.NewScheme()
	if err := velerov1.AddToScheme(scheme); err != nil {
		return err
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return err
	}

	m.clusters[name] = &ClusterConfig{
		Name:       name,
		KubeConfig: kubeconfigPath,
		Client:     k8sClient,
	}

	return nil
}

// GetClient returns the Kubernetes client for a cluster
func (m *Manager) GetClient(clusterName string) (client.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cluster, exists := m.clusters[clusterName]
	if !exists {
		return nil, fmt.Errorf("cluster %s not found", clusterName)
	}

	return cluster.Client, nil
}

// ListClusters returns all available cluster names
func (m *Manager) ListClusters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clusters))
	for name := range m.clusters {
		names = append(names, name)
	}
	return names
}

// GetDefaultCluster returns the first cluster (used as default)
func (m *Manager) GetDefaultCluster() string {
	clusters := m.ListClusters()
	if len(clusters) > 0 {
		return clusters[0]
	}
	return ""
}

// buildConfig creates rest.Config from kubeconfig path
func (m *Manager) buildConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	return rest.InClusterConfig()
}
