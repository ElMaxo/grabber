package repository

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	gormigrate "gopkg.in/gormigrate.v1"

	"github.com/jinzhu/gorm"
)

type repo struct {
	db *gorm.DB
}

// NewPostgresRepository creates new postgres repository instance with specified connection string
func NewPostgresRepository(url string) (Repository, error) {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(30 * time.Second)
	repo := &repo{db}
	return repo, nil
}

// Migrate performs schema migration on specified connection
func (r *repo) Migrate() error {
	m := gormigrate.New(r.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				type Job struct {
					CreatedAt time.Time `json:"-"`
					ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
					Url       string    `gorm:"not null"`
					Query     postgres.Jsonb
					Period    int64 `gorm:"not null"`
				}

				type Article struct {
					CreatedAt   time.Time `json:"-"`
					ID          string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
					Link        string
					Title       string
					Description string
				}

				if err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
					return err
				}

				if err := createTableIfNotExists("jobs", &Job{}, tx); err != nil {
					return err
				}

				if err := createTableIfNotExists("articles", &Article{}, tx); err != nil {
					return err
				}

				if err := tx.Exec("ALTER TABLE articles ADD column IF NOT EXISTS fts_vector tsvector").Error; err != nil {
					return err
				}

				return tx.Exec("CREATE INDEX IF NOT EXISTS fts_articles_idx ON articles USING gin(fts_vector)").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("jobs", "articles").Error
			},
		},
	})

	return m.Migrate()
}

// GetJobs returns all jobs records from the database
func (r *repo) GetJobs() ([]*Job, error) {
	var jobs []*Job
	if err := r.db.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJob returns job record with specified ID from the database
func (r *repo) GetJob(id string) (*Job, error) {
	job := &Job{ID: id}
	if err := r.db.First(job).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return job, nil
}

// CreateJob creates job record in the database
func (r *repo) CreateJob(job *Job) error {
	return r.db.Create(job).Error
}

// DeleteJob deletes job record from the database
func (r *repo) DeleteJob(job *Job) error {
	return r.db.Delete(job).Error
}

// CreateArticle creates new article record in the database
func (r *repo) CreateArticle(article *Article) error {
	if err := r.db.Create(article).Error; err != nil {
		return err
	}
	return r.db.Exec("UPDATE articles SET fts_vector=(setweight(to_tsvector('russian', title), 'A') || setweight(to_tsvector('russian', description), 'B')) WHERE id=?",
		article.ID).Error
}

// GetArticles returns articles record matches search phrase with paging
func (r *repo) GetArticles(q string, page, rowsPerPage int64) ([]*Article, error) {
	var articles []*Article
	offset := (page - 1) * rowsPerPage
	query := r.db.Table("articles").Select(` 
			count(*) OVER() as rows_count,
			articles.title, 
			articles.link,
			articles.description`)
	if q != "" {
		query = query.Where("articles.fts_vector @@ plainto_tsquery('russian', ?)", q).
			Order(gorm.Expr("ts_rank(articles.fts_vector, plainto_tsquery('russian', ?)), created_at DESC", q))
	} else {
		query = query.Order("created_at", true)
	}
	query = query.Limit(rowsPerPage).Offset(offset)
	if err := query.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *repo) dropTables(tables []interface{}) error {
	return r.db.DropTableIfExists(tables...).Error
}

func createTableIfNotExists(table string, schema interface{}, db *gorm.DB) error {
	if !db.HasTable(table) {
		return db.Table(table).CreateTable(schema).Error
	}

	return nil
}
