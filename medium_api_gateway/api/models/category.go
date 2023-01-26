package models

type Category struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type CreateCategoryRequest struct {
	Title string `json:"title" binding:"required,max=100"`
}

type GetCategoriesResponse struct {
	Categories []*Category `json:"categories"`
	Count      int64       `json:"count"`
}
