package usecase

type Usecase struct {
	User      userUC
	MainCateg mainCategUC
	SubCateg  subCategUC
}

func New(u UserModel, m MainCategModel, s SubCategModel, i IconModel) *Usecase {
	return &Usecase{
		User:      *newUserUC(u),
		MainCateg: *newMainCategUC(m, i),
		SubCateg:  *newSubCategUC(s),
	}
}
