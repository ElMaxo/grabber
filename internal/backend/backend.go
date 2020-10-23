package backend

import (
	"grabber/internal/backend/article"
	"grabber/internal/backend/job"
	"grabber/internal/grabber"
	"grabber/internal/repository"
	"grabber/internal/rest/restapi"
	"grabber/internal/rest/restapi/operations"
	"grabber/internal/rest/restapi/operations/articles"
	"grabber/internal/rest/restapi/operations/description"
	"grabber/internal/rest/restapi/operations/jobs"
	"net/http"

	"github.com/dre1080/recovr"

	"github.com/go-openapi/runtime/middleware"
)

// InitAndBindToAPI binds the handlers to the API
func InitAndBindToAPI(repo repository.Repository, grab grabber.Grabber, api *operations.GrabberAPI) http.Handler {
	jobsHandler := job.NewHandler(job.NewService(repo, grab))
	articlesHandler := article.NewHandler(article.NewService(repo))

	api.JobsGetJobsHandler = jobs.GetJobsHandlerFunc(jobsHandler.GetJobs)
	api.JobsCreateJobHandler = jobs.CreateJobHandlerFunc(jobsHandler.AddJob)
	api.JobsDeleteJobHandler = jobs.DeleteJobHandlerFunc(jobsHandler.DeleteJob)

	api.ArticlesGetArticlesHandler = articles.GetArticlesHandlerFunc(articlesHandler.GetArticles)

	api.DescriptionGetAPIHandler = description.GetAPIHandlerFunc(func(params description.GetAPIParams) middleware.Responder {
		return description.NewGetAPIOK().WithPayload(restapi.SwaggerJSON)
	})

	return recovr.New()(api.Serve(nil))
}
