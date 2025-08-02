package dto

type CreateAccountRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
