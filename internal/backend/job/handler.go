package job

import (
	"encoding/json"
	"grabber/internal/helper"
	"grabber/internal/httperror"
	"grabber/internal/repository"
	"grabber/internal/rest/models"
	"grabber/internal/rest/restapi/operations/jobs"

	"github.com/jinzhu/gorm/dialects/postgres"

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

func (h *Handler) GetJobs(jobs.GetJobsParams) middleware.Responder {
	foundJobs, err := h.service.GetJobs()
	if err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	response, err := convertJobsToResponse(foundJobs)
	if err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	return jobs.NewGetJobsOK().WithPayload(response)
}

func (h *Handler) AddJob(params jobs.CreateJobParams) middleware.Responder {
	if *params.Job.PeriodSec <= 0 || *params.Job.PeriodSec > 3600 {
		return httperror.NewBadRequestError("invalid period")
	}
	var query postgres.Jsonb
	queryBytes, err := json.Marshal(params.Job.Query)
	if err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	if err := query.Scan(queryBytes); err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	job := &repository.Job{
		Url:    *params.Job.URL,
		Query:  query,
		Period: *params.Job.PeriodSec,
	}
	if err := h.service.AddJob(job); err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	response, err := convertJobToResponse(job)
	if err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	return jobs.NewCreateJobOK().WithPayload(response)
}

func (h *Handler) DeleteJob(params jobs.DeleteJobParams) middleware.Responder {
	if err := h.service.DeleteJob(params.ID); err != nil {
		return httperror.ConvertHTTPErrorToResponse(err)
	}
	return jobs.NewDeleteJobNoContent()
}

func convertJobToResponse(job *repository.Job) (*models.Job, error) {
	query := &models.Query{}
	if err := json.Unmarshal(job.Query.RawMessage, query); err != nil {
		return nil, err
	}
	return &models.Job{
		ID:        job.ID,
		PeriodSec: helper.Int64Ptr(job.Period),
		Query:     query,
		URL:       helper.StringPtr(job.Url),
	}, nil
}

func convertJobsToResponse(jobs []*repository.Job) ([]*models.Job, error) {
	response := make([]*models.Job, 0, len(jobs))
	for _, job := range jobs {
		converted, err := convertJobToResponse(job)
		if err != nil {
			return nil, err
		}
		response = append(response, converted)
	}
	return response, nil
}
