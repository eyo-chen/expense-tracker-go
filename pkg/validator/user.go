package validator

// Signup validates email, password and name for signup
func (v *Validator) Signup(email, password, name string) bool {
	v.checkEmail(email)
	v.checkPassword(password)
	v.checkName(name)
	return v.Valid()
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
