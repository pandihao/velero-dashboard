package handler

import (
	"net/http"

	"velero-api-server/internal/model"
	"velero-api-server/internal/service"
	"velero-api-server/pkg/cluster"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, clusterMgr *cluster.Manager, namespace string) {
	svc := service.NewVeleroService(namespace)

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.APIResponse{Code: 200, Message: "ok"})
	})

	// Cluster management
	r.GET("/api/v1/clusters", listClusters(clusterMgr))

	api := r.Group("/api/v1")
	api.Use(clusterMiddleware(clusterMgr))
	{
		// Backup
		api.POST("/backups", createBackup(svc))
		api.GET("/backups", listBackups(svc))
		api.GET("/backups/:name", getBackup(svc))
		api.DELETE("/backups/:name", deleteBackup(svc))

		// Restore
		api.POST("/restores", createRestore(svc))
		api.GET("/restores", listRestores(svc))
		api.GET("/restores/:name", getRestore(svc))

		// Schedule
		api.POST("/schedules", createSchedule(svc))
		api.GET("/schedules", listSchedules(svc))
		api.GET("/schedules/:name", getSchedule(svc))
		api.DELETE("/schedules/:name", deleteSchedule(svc))

		// BackupStorageLocation
		api.GET("/storage-locations", listBSLs(svc))
		api.GET("/storage-locations/:name", getBSL(svc))
	}
}

// clusterMiddleware extracts cluster name from query/header and sets k8s client in context
func clusterMiddleware(clusterMgr *cluster.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get cluster name from query parameter or header
		clusterName := ctx.Query("cluster")
		if clusterName == "" {
			clusterName = ctx.GetHeader("X-Cluster")
		}
		if clusterName == "" {
			clusterName = clusterMgr.GetDefaultCluster()
		}

		// Get k8s client for the cluster
		k8sClient, err := clusterMgr.GetClient(clusterName)
		if err != nil {
			errResponse(ctx, http.StatusBadRequest, err)
			ctx.Abort()
			return
		}

		// Store client in context
		ctx.Set("k8sClient", k8sClient)
		ctx.Set("clusterName", clusterName)
		ctx.Next()
	}
}

// listClusters returns all available clusters
func listClusters(clusterMgr *cluster.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clusters := clusterMgr.ListClusters()
		success(ctx, gin.H{
			"clusters": clusters,
			"default":  clusterMgr.GetDefaultCluster(),
		})
	}
}

func success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, model.APIResponse{
		Code:    201,
		Message: "created",
		Data:    data,
	})
}

func errResponse(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, model.APIResponse{
		Code:    code,
		Message: err.Error(),
	})
}
