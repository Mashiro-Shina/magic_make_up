package entities

type User struct {
	ID int
	Phone string
	Password string
	Signature string
	Name string
	Avatar string
	Level int
	FollowingNum int
	FollowersNum int
}

func NewUser(phone string, password string, name string, signature string, avatar string) *User {
	return &User{
		Phone:     phone,
		Password:  password,
		Name:      name,
		Signature: signature,
		Avatar: avatar,
	}
}
