package repository

type Repository interface {
	Migrate() error
	GetJobs() ([]*Job, error)
	GetJob(id string) (*Job, error)
	CreateJob(job *Job) error
	DeleteJob(job *Job) error
	CreateArticle(article *Article) error
	GetArticles(q string, page, rowsPerPage int64) ([]*Article, error)
}
