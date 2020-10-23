package job

import (
	"encoding/json"
	"grabber/internal/grabber"
	"grabber/internal/httperror"
	"grabber/internal/repository"
	"log"
	"sync"
	"time"
)

type Service interface {
	GetJobs() ([]*repository.Job, error)
	AddJob(job *repository.Job) error
	DeleteJob(id string) error
}

type service struct {
	grab       grabber.Grabber
	repo       repository.Repository
	mu         sync.Mutex
	jobWorkers map[string]chan struct{}
}

// NewService creates new jobs service instance with specified repository
func NewService(repo repository.Repository, grab grabber.Grabber) Service {
	return &service{repo: repo, grab: grab, jobWorkers: map[string]chan struct{}{}}
}

// GetJobs returns all articles grabbing jobs from the database
func (s *service) GetJobs() ([]*repository.Job, error) {
	return s.repo.GetJobs()
}

// AddJob adds new articles grabbing job to the database and starts its execution
func (s *service) AddJob(job *repository.Job) error {
	if err := s.repo.CreateJob(job); err != nil {
		return err
	}
	var query grabber.Query
	if err := json.Unmarshal(job.Query.RawMessage, &query); err != nil {
		return err
	}
	doneChan := make(chan struct{})
	go s.jobWorker(time.Duration(job.Period)*time.Second, job.Url, query, doneChan)
	s.mu.Lock()
	s.jobWorkers[job.ID] = doneChan
	s.mu.Unlock()
	return nil
}

// DeleteJob deletes job with specified ID from the database and stops its execution
func (s *service) DeleteJob(id string) error {
	foundJob, err := s.repo.GetJob(id)
	if err != nil {
		return httperror.NewNotFoundError("job not found")
	}
	s.mu.Lock()
	jobCh, ok := s.jobWorkers[foundJob.ID]
	s.mu.Unlock()
	if ok {
		jobCh <- struct{}{}
	}
	return s.repo.DeleteJob(foundJob)
}

func (s *service) jobWorker(period time.Duration, url string, query grabber.Query, doneChan <-chan struct{}) {
	log.Print("job started")
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			articles, err := s.grab.GrabNews(url, query)
			if err != nil {
				log.Print("failed to grab articles: ", err)
				continue
			}
			for _, article := range articles {
				repoArticle := &repository.Article{
					Link:        article.Link,
					Title:       article.Title,
					Description: article.Description,
				}
				if err := s.repo.CreateArticle(repoArticle); err != nil {
					log.Print("failed to save article to DB: ", err)
				}
			}
		case <-doneChan:
			log.Print("worker done")
			ticker.Stop()
			return
		}
	}
}
