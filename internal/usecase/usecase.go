package usecase

type Usecase struct {
	User        userUC
	MainCateg   mainCategUC
	SubCateg    subCategUC
	Transaction transactionUC
}

func New(u UserModel, m MainCategModel, s SubCategModel, i IconModel, t TransactionModel) *Usecase {
	return &Usecase{
		User:        *newUserUC(u),
		MainCateg:   *newMainCategUC(m, i),
		SubCateg:    *newSubCategUC(s, m),
		Transaction: *newTransactionUC(t, m, s),
	}
}
