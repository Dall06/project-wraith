package rules

type User struct {
	ID       string
	Username string
	Email    string
	Name     string
	Phone    string
	Password string
}

type Reset struct {
	ID          string
	Username    string
	Email       string
	Phone       string
	NewPassword string
	Token       string
}
