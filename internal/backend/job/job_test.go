package job_test

import (
	"grabber/internal/backend"
	"grabber/internal/config"
	"grabber/internal/grabber"
	"grabber/internal/helper"
	"grabber/internal/repository"
	"grabber/internal/rest/client"
	"grabber/internal/rest/client/jobs"
	"grabber/internal/rest/models"
	"grabber/internal/rest/restapi/operations"
	"log"
	"net/http"
	"testing"

	"github.com/go-openapi/loads"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/suite"
)

type JobsSuite struct {
	suite.Suite
	repo       repository.Repository
	jobsClient *client.Grabber
}

func (js *JobsSuite) SetupSuite() {
	swaggerSpec, err := loads.Spec("../../../api/grabber.swagger.yml")
	if err != nil {
		log.Fatal("failed to load spec: ", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	repo, err := repository.NewPostgresRepository(cfg.DbURL)
	if err != nil {
		log.Fatal("failed to create repository: ", err)
	}
	js.repo = repo
	api := operations.NewGrabberAPI(swaggerSpec)

	handler := backend.InitAndBindToAPI(repo, grabber.NewHtmlGrabber(), api)

	httpClient := &http.Client{Transport: helper.NewTransport(handler)}
	c := runtimeClient.NewWithClient(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes, httpClient)
	c.Debug = true
	js.jobsClient = client.New(c, nil)
}

func TestJobs(t *testing.T) {
	suite.Run(t, new(JobsSuite))
}

func (js *JobsSuite) TestJobCycle() {
	createParams := &jobs.CreateJobParams{
		Job: &models.Job{
			URL:       helper.StringPtr("https://lenta.ru"),
			PeriodSec: helper.Int64Ptr(60),
			Query: &models.Query{
				DescriptionSelector: &models.Selector{
					Value: "inner_text",
					Xpath: "//div[@itemprop='articleBody']",
				},
				ItemsSelector: &models.Selector{
					Xpath: "//div[@class='span4']/div[@class='item']/a",
				},
				LinkSelector: &models.Selector{
					Attr:  "href",
					Value: "attr",
				},
				TitleSelector: &models.Selector{
					Value: "text",
				},
				FollowLinkForDescription: true,
			},
		},
	}
	createdJob, err := js.jobsClient.Jobs.CreateJob(createParams)
	js.NoError(err)
	js.NotEmpty(createdJob.Payload.ID)
	foundJob, err := js.repo.GetJob(createdJob.Payload.ID)
	js.NoError(err)
	js.NotNil(foundJob)

	foundJobs, err := js.jobsClient.Jobs.GetJobs(jobs.NewGetJobsParams())
	js.NoError(err)
	js.NotEmpty(foundJobs.Payload)
	found := false
	for _, job := range foundJobs.Payload {
		if job.ID == foundJob.ID {
			found = true
		}
	}
	js.True(found)

	_, err = js.jobsClient.Jobs.DeleteJob(jobs.NewDeleteJobParams().WithID(foundJob.ID))
	js.NoError(err)
	foundJob, err = js.repo.GetJob(foundJob.ID)
	js.NoError(err)
	js.Nil(foundJob)
}
