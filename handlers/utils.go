package handlers

import "net/url"

type userRequestBody struct {
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

func (u *userRequestBody) validate() map[string]interface{} {
	errs := url.Values{}
	// check if the nickname empty
	if u.Nickname == "" {
		errs.Add("nickname", "The nickname field is required!")
	}

	// check if the first_name empty
	if u.FirstName == "" {
		errs.Add("first_name", "The first_name field is required!")
	}

	// check if the last_name empty
	if u.LastName == "" {
		errs.Add("last_name", "The last_name field is required!")
	}

	// check if the password empty
	if u.Password == "" {
		errs.Add("password", "The password field is required!")
	}

	// check if the email empty
	if u.Email == "" {
		errs.Add("email", "The email field is required!")
	}

	// check if the title empty
	if u.Country == "" {
		errs.Add("country", "The country field is required!")
	}
	validationErrors := map[string]interface{}{"validationError": errs}

	return validationErrors

}
