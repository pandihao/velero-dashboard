package handler

import (
	"fmt"
	"net/http"

	"velero-api-server/internal/model"
	"velero-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createRestore(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		var req model.CreateRestoreRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errResponse(ctx, http.StatusBadRequest, err)
			return
		}
		resp, err := svc.CreateRestore(ctx.Request.Context(), k8sClient.(client.Client), req)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		created(ctx, resp)
	}
}

func listRestores(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		backupName := ctx.Query("backupName")
		resp, err := svc.ListRestores(ctx.Request.Context(), k8sClient.(client.Client), backupName)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}

func getRestore(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		name := ctx.Param("name")
		if name == "" {
			errResponse(ctx, http.StatusBadRequest, fmt.Errorf("name parameter is required"))
			return
		}
		resp, err := svc.GetRestore(ctx.Request.Context(), k8sClient.(client.Client), name)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}
