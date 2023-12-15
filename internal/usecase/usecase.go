package usecase

type Usecase struct {
	User      userUC
	MainCateg mainCategUC
}

func New(u UserModel, m MainCategModel) *Usecase {
	return &Usecase{
		User:      *newUserUC(u),
		MainCateg: *newMainCategUC(m),
	}
}
