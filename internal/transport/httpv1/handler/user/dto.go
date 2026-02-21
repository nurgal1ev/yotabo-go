package user

type RegisterInput struct {
	Body struct {
		FirstName string `json:"firstname" minLength:"1" maxLength:"50" pattern:"^[\\p{L}]+$"`
		LastName  string `json:"lastname"  minLength:"1" maxLength:"50" pattern:"^[\\p{L}]+$"`
		Username  string `json:"username"  minLength:"3" maxLength:"12"`
		Email     string `json:"email"     format:"email"`
		Password  string `json:"password"  minLength:"7" maxLength:"12"`
	}
}

type RegisterOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

type LoginInput struct {
	Body struct {
		Username string `json:"username" minLength:"3" maxLength:"12"`
		Password string `json:"password" minLength:"7" maxLength:"12"`
	}
}

type LoginOutput struct {
	Body struct {
		AccessToken string `json:"accessToken"`
	}
}
