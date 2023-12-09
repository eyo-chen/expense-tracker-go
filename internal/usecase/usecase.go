package usecase

type Usecase struct {
	User userUC
}

func New(UserModel UserModel) *Usecase {
	return &Usecase{
		User: *newUserUC(UserModel),
	}
}
