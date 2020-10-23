package article

import (
	"grabber/internal/helper"
	"grabber/internal/httperror"
	"grabber/internal/repository"
	"grabber/internal/rest/models"
	"grabber/internal/rest/restapi/operations/articles"

	"github.com/go-openapi/runtime/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetArticles(params articles.GetArticlesParams) middleware.Responder {
	art, err := h.service.GetArticles(helper.SafeStringGet(params.Q), helper.SafeInt64Get(params.Page), helper.SafeInt64Get(params.RowsPerPage))
	if err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}

	return articles.NewGetArticlesOK().WithPayload(convertArticlesToResponse(art))
}

func convertArticleToResponse(article *repository.Article) *models.Article {
	return &models.Article{
		Description: article.Description,
		Link:        helper.StringPtr(article.Link),
		Title:       helper.StringPtr(article.Title),
	}
}

func convertArticlesToResponse(articles []*repository.Article) []*models.Article {
	response := make([]*models.Article, 0, len(articles))
	for _, article := range articles {
		response = append(response, convertArticleToResponse(article))
	}
	return response
}
