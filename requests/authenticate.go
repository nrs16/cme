package requests

type RegisterBody struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterBody) Validate() error {
	return nil
}
