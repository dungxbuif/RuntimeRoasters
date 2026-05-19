package domain

type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Role     string   `json:"role"`
	OrgID    string   `json:"org_id,omitempty"`
	StoreIDs []string `json:"store_ids,omitempty"`
}

type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Name     string   `json:"name" binding:"required"`
	Role     string   `json:"role" binding:"required"`
	OrgID    string   `json:"org_id"`
	StoreIDs []string `json:"store_ids"`
}

type CreateUserResponse struct {
	User *User `json:"user"`
}

type AcceptLoginRequest struct {
	LoginChallenge string `json:"login_challenge" binding:"required"`
	Subject        string `json:"subject" binding:"required"`
	Remember       bool   `json:"remember"`
	RememberFor    int64  `json:"remember_for"`
}

type AcceptLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}
