package domain

type User struct {
	ID             string
	Login          string
	HashedPassword string
}

func NewUser(login, hashedPassword string) User {
	return User{
		Login:          login,
		HashedPassword: hashedPassword,
	}
}
