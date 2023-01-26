package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/medium-project/medium_post_service/storage/repo"
)

type commentRepo struct {
	db *sqlx.DB
}

func NewComment(db *sqlx.DB) repo.CommentStorageI {
	return &commentRepo{
		db: db,
	}
}

func (pr *commentRepo) Create(c *repo.Comment) (*repo.Comment, error) {
	var (
		updatedAt sql.NullTime
	)
	query := `
		INSERT INTO comments (
			post_id,
			user_id,
			description
		) VALUES ($1, $2, $3)
		RETURNING 
		id, 
		created_at,
		updated_at
	`

	err := pr.db.QueryRow(
		query,
		c.PostID,
		c.UserID,
		c.Description,
	).Scan(
		&c.ID,
		&c.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	c.UpdatedAt = updatedAt.Time

	return c, nil
}

func (cr *commentRepo) Update(c *repo.Comment) (*repo.Comment, error) {
	var (
		res       repo.Comment
		updatedAt sql.NullTime
	)
	query := `
		UPDATE comments SET
			description = $1,
			updated_at = $2
	    WHERE id = $3
		RETURNING 
			id,
			description,
			post_id,
			user_id,
			created_at,
			updated_at
	`

	err := cr.db.QueryRow(
		query,
		c.Description,
		time.Now(),
		c.ID,
	).Scan(
		&res.ID,
		&res.Description,
		&res.PostID,
		&res.UserID,
		&res.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	res.UpdatedAt = updatedAt.Time

	return &res, nil
}

func (cr *commentRepo) Delete(commentId int64) error {
	query := `
		DELETE FROM comments WHERE id = $1
	`

	result, err := cr.db.Exec(
		query,
		commentId,
	)
	if err != nil {
		return err
	}

	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pr *commentRepo) GetAll(params *repo.GetCommentsParams) (*repo.GetAllCommentsResult, error) {
	result := repo.GetAllCommentsResult{
		Comments: make([]*repo.Comment, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	orderBy := " ORDER BY created_at DESC "
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s", params.SortBy)
	}

	filter := ""

	if params.PostID != 0 {
		filter += fmt.Sprintf(" WHERE post_id = %d", params.PostID)
	}

	query := `
		SELECT
			id,
			post_id,
			user_id,
			description,
			created_at,
			updated_at
		FROM comments 
	` + filter + orderBy + limit

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			comment   repo.Comment
			updatedAt sql.NullTime
		)
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Description,
			&comment.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		comment.UpdatedAt = updatedAt.Time
		result.Comments = append(result.Comments, &comment)
	}

	queryCount := "SELECT count(1) FROM comments " + filter

	err = pr.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
