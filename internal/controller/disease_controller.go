package controller

import (
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/response"
	"graph-med/internal/service"
)

type DiseaseController struct {
	kgService *service.KGService
}

func NewDiseaseController(kgService *service.KGService) *DiseaseController {
	return &DiseaseController{
		kgService: kgService,
	}
}

// GetLabels 获取疾病库所有标签
func (d *DiseaseController) GetLabels(ctx *gin.Context) {
	labels, err := d.kgService.GetLabels()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, labels)
}

// Search 搜索疾病库
func (d *DiseaseController) Search(ctx *gin.Context) {
	var req struct {
		Type  string `json:"type" binding:"required"`
		Label string `json:"label" binding:"required"`
		Query string `json:"query" binding:"required"`
		Limit int    `json:"limit"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, err)
		return
	}

	// 限制 limit
	if req.Limit > 200 || req.Limit <= 0 {
		req.Limit = 200
	}

	switch req.Type {
	case "subgraph":
		subgraph, err := d.kgService.GetSubgraph(req.Label, req.Query, req.Limit)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, subgraph)
		break
	case "nodes":
		diseases, err := d.kgService.GetNodes(req.Label, req.Query, req.Limit)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, diseases)
		break
	default:
		response.Success(ctx, nil)
		break
	}

}
