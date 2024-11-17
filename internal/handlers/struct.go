package handlers

type User struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
	Ref   string `json:"ref"`
}