package user

type User struct {
	ID             string
	Login          string
	HashedPassword string
}

func New(login, hashedPassword string) User {
	return User{
		Login:          login,
		HashedPassword: hashedPassword,
	}
}
