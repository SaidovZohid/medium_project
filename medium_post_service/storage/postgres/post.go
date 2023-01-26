package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/medium-project/medium_post_service/pkg/utils"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPost(db *sqlx.DB) repo.PostStorageI {
	return &postRepo{
		db: db,
	}
}

func (pr *postRepo) Create(p *repo.Post) (*repo.Post, error) {
	var (
		updatedAt sql.NullTime
	)
	query := `
		INSERT INTO posts(
			title,
			description,
			image_url,
			user_id,
			category_id
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, views_count
	`

	err := pr.db.QueryRow(
		query,
		p.Title,
		p.Description,
		utils.NullString(p.ImageUrl),
		p.UserID,
		p.CategoryID,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&updatedAt,
		&p.ViewsCount,
	)

	if err != nil {
		return nil, err
	}
	p.UpdatedAt = updatedAt.Time

	return p, nil
}

func (pr *postRepo) Get(post_id int64) (*repo.Post, error) {
	var (
		res       repo.Post
		updatedAt sql.NullTime
		imageUrl  sql.NullString
	)
	queryViews := "UPDATE posts SET views_count = views_count + 1 WHERE id = $1"
	_, err := pr.db.Exec(queryViews, post_id)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			p.id,
			p.title,
			p.description,
			p.image_url,
			p.user_id,
			p.category_id,
			p.created_at,
			p.updated_at,
			p.views_count
		FROM posts p 
		WHERE p.id = $1 
	`

	err = pr.db.QueryRow(
		query,
		post_id,
	).Scan(
		&res.ID,
		&res.Title,
		&res.Description,
		&imageUrl,
		&res.UserID,
		&res.CategoryID,
		&res.CreatedAt,
		&updatedAt,
		&res.ViewsCount,
	)

	if err != nil {
		return nil, err
	}
	res.ImageUrl = imageUrl.String
	res.UpdatedAt = updatedAt.Time

	return &res, nil
}

func (pr *postRepo) Update(p *repo.Post) (*repo.Post, error) {
	var (
		res      repo.Post
		imageUrl sql.NullString
	)
	query := `
		UPDATE posts SET
			title = $1,
			description = $2,
			image_url = $3,
			user_id = $4,
			category_id = $5,
			updated_at = $6
	    WHERE id = $7
		RETURNING 
			id,
			title,
			description,
			image_url,
			user_id,
			category_id,
			created_at,
			updated_at,
			views_count
	`

	err := pr.db.QueryRow(
		query,
		p.Title,
		p.Description,
		utils.NullString(p.ImageUrl),
		p.UserID,
		p.CategoryID,
		time.Now(),
		p.ID,
	).Scan(
		&res.ID,
		&res.Title,
		&res.Description,
		&imageUrl,
		&res.UserID,
		&res.CategoryID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.ViewsCount,
	)

	if err != nil {
		return nil, err
	}
	res.ImageUrl = imageUrl.String

	return &res, nil
}

func (pr *postRepo) Delete(post_id int64) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`

	res, err := pr.db.Exec(
		query,
		post_id,
	)
	if res, _ := res.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	return nil
}

func (pr *postRepo) GetAll(params *repo.GetPostsParams) (*repo.GetAllPostResult, error) {
	result := repo.GetAllPostResult{
		Posts: make([]*repo.Post, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := " WHERE true"
	if params.Search != "" {
		filter += " AND title ILIKE '%s'" + "%" + params.Search + "%"
	}

	if params.UserID != 0 {
		filter += fmt.Sprintf(" AND user_id = %d", params.UserID)
	}

	if params.CategoryID != 0 {
		filter += fmt.Sprintf(" AND category_id = %d", params.CategoryID)
	}

	orderBy := " ORDER BY created_at DESC"
	if params.SortByDate != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s", params.SortByDate)
	}

	query := `
		SELECT
			id,
			title,
			description,
			image_url,
			user_id,
			category_id,
			created_at,
			updated_at,
			views_count
		FROM posts
	` + filter + orderBy + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			post      repo.Post
			updatedAt sql.NullTime
			imageUrl  sql.NullString
		)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&imageUrl,
			&post.UserID,
			&post.CategoryID,
			&post.CreatedAt,
			&updatedAt,
			&post.ViewsCount,
		)
		if err != nil {
			return nil, err
		}
		post.ImageUrl = imageUrl.String
		post.UpdatedAt = updatedAt.Time
		result.Posts = append(result.Posts, &post)
	}

	queryCount := "SELECT count(1) FROM posts " + filter

	err = pr.db.QueryRow(queryCount).Scan(&result.Count)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
