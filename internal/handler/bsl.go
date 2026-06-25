package handler

import (
	"fmt"
	"net/http"

	"velero-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func listBSLs(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		resp, err := svc.ListBSLs(ctx.Request.Context(), k8sClient.(client.Client))
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}

func getBSL(svc *service.VeleroService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		k8sClient, _ := ctx.Get("k8sClient")
		name := ctx.Param("name")
		if name == "" {
			errResponse(ctx, http.StatusBadRequest, fmt.Errorf("name parameter is required"))
			return
		}
		resp, err := svc.GetBSL(ctx.Request.Context(), k8sClient.(client.Client), name)
		if err != nil {
			errResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		success(ctx, resp)
	}
}
