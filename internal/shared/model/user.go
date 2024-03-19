package model

type User struct {
	ID             string `json:"id"`
	Login          string `json:"login"`
	HashedPassword string `json:"password"`
}

func NewUser(login, hashedPassword string) User {
	return User{
		Login:          login,
		HashedPassword: hashedPassword,
	}
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
