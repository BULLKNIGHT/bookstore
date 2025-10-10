package models

// User represents a user in the system
// @Description User information for authentication and authorization
type User struct {
	Name string `json:"name" example:"john_doe"`
	Role string `json:"role" example:"guest"`
}

func (user *User) IsValid() bool {
	return user.Name != "" && user.Role != ""
}
