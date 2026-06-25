package handler

import (
	"fmt"
	"net/http"

	"velero-api-server/internal/model"
	"velero-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createBackup(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		var req model.CreateBackupRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errResponse(ctx, http.StatusBadRequest, err)
			return
		}
		resp, err := svc.CreateBackup(ctx.Request.Context(), k8sClient.(client.Client), req)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		created(ctx, resp)
	}
}

func listBackups(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		resp, err := svc.ListBackups(ctx.Request.Context(), k8sClient.(client.Client))
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}

func getBackup(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		name := ctx.Param("name")
		if name == "" {
			errResponse(ctx, http.StatusBadRequest, fmt.Errorf("name parameter is required"))
			return
		}
		resp, err := svc.GetBackup(ctx.Request.Context(), k8sClient.(client.Client), name)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}

func deleteBackup(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		name := ctx.Param("name")
		if name == "" {
			errResponse(ctx, http.StatusBadRequest, fmt.Errorf("name parameter is required"))
			return
		}
		if err := svc.DeleteBackup(ctx.Request.Context(), k8sClient.(client.Client), name); err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, nil)
	}
}
