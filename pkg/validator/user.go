package validator

// Signup validates email, password and name for signup
func (v *Validator) Signup(email, password, name string) bool {
	v.checkEmail(email)
	v.checkPassword(password)
	v.checkName(name)
	return v.Valid()
}

// Login validates email and password for login
func (v *Validator) Login(email, password string) bool {
	v.checkEmail(email)
	v.checkPassword(password)
	return v.Valid()
}

// Token validates refresh token for token
func (v *Validator) Token(refreshToken string) bool {
	v.checkRefreshToken(refreshToken)
	return v.Valid()
}

func (v *Validator) checkRefreshToken(refreshToken string) {
	v.Check(len(refreshToken) > 0, "refresh_token", "Refresh token can't be empty")
}

func (v *Validator) checkEmail(email string) {
	v.Check(Matches(email, EmailRX), "email", "Invalid email address")
}

func (v *Validator) checkPassword(password string) {
	v.Check(len(password) >= 8, "password", "Password must be at least 8 characters long")
}

func (v *Validator) checkName(name string) {
	v.Check(len(name) > 0, "name", "Name can't be empty")
}
