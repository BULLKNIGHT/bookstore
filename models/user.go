package models

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func (user *User) IsValid() bool {
	return user.Name != "" && user.Role != ""
}
