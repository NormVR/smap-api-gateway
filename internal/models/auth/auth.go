package auth

type UserData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
