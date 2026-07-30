package dto

type UserRes struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
