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

type UserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}
type UpdateUserInput struct {
	Body struct {
		FirstName *string `json:"firstname"`
		LastName  *string `json:"lastname"`
		Username  *string `json:"username"`
		Email     *string `json:"email"`
	}
}

type ChangePasswordInput struct {
	Body struct {
		CurrentPassword string `json:"currentPassword" minLength:"7" maxLength:"12"`
		NewPassword     string `json:"newPassword" minLength:"7" maxLength:"12"`
	}
}
