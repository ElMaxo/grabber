package article_test

import (
	"grabber/internal/backend"
	"grabber/internal/config"
	"grabber/internal/grabber"
	"grabber/internal/helper"
	"grabber/internal/repository"
	"grabber/internal/rest/client"
	"grabber/internal/rest/client/articles"
	"grabber/internal/rest/restapi/operations"
	"log"
	"net/http"
	"testing"

	"github.com/go-openapi/loads"
	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/suite"
)

type ArticlesSuite struct {
	suite.Suite
	repo           repository.Repository
	articlesClient *client.Grabber
}

func (as *ArticlesSuite) SetupSuite() {
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
	as.repo = repo
	api := operations.NewGrabberAPI(swaggerSpec)

	handler := backend.InitAndBindToAPI(repo, grabber.NewHtmlGrabber(), api)

	httpClient := &http.Client{Transport: helper.NewTransport(handler)}
	c := runtimeClient.NewWithClient(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes, httpClient)
	c.Debug = true
	as.articlesClient = client.New(c, nil)
}

func TestArticles(t *testing.T) {
	suite.Run(t, new(ArticlesSuite))
}

func (as *ArticlesSuite) TestGetArticles() {
	article := &repository.Article{
		ID:    "",
		Link:  "https://lenta.ru",
		Title: "Коронавирусом заразился каждый сотый россиянин",
		Description: `Каждый сотый россиянин заразился коронавирусом с начала пандемии, подсчитало РИА Новости на основе данных Росстата.

По информации агентства, число пациентов с COVID-19 в России превысило один процент населения.

Заместитель директора по клинико-аналитической работе ЦНИИ Эпидемиологии Роспотребнадзора Наталья Пшеничная в разговоре с РИА Новости отметила, что 75 процентов инфицированных уже поправились.`,
	}
	as.NoError(as.repo.CreateArticle(article))
	findParams := &articles.GetArticlesParams{}
	response, err := as.articlesClient.Articles.GetArticles(findParams)
	as.NoError(err)
	as.NotEmpty(response.Payload)
	as.NotEmpty(response.Payload[0].Title)
	as.NotEmpty(response.Payload[0].Link)

	findParams.Q = helper.StringPtr("коронавирус зараза")
	response, err = as.articlesClient.Articles.GetArticles(findParams)
	as.NoError(err)
	as.NotEmpty(response.Payload)
	as.NotEmpty(response.Payload[0].Title)
	as.NotEmpty(response.Payload[0].Link)
}
