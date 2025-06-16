package authModel

type User struct {
	ID       string `json:"id"`
	Username string	`json:"username"`
	Password string `json:"password"`
}

type Request struct {
	Username string `json:"username"`
    Password string `json:"password"`
}