package requests

type RegisterForm struct {
	Name            string `valid:"required"`
	Password        string `valid:"required,eqfield=ConfirmPassword"`
	PasswordConfirm string `valid:"required"`
	Agree           string `valid:"required"`
}

