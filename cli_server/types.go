package cli_server

type UserSignUp struct {
	name     string
	surname  string
	email    string
	password string
}

type CompanySignUp struct {
	title    string
	email    string
	password string
}

type UserSignIn struct {
	email    string
	password string
}

type SignUpResponse struct {
	return_type  string
	user_data    UserSignUp
	company_data CompanySignUp
}
