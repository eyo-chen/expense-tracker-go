package handler

type Handler struct {
	User        userHandler
	MainCateg   mainCategHandler
	SubCateg    subCategHandler
	Transaction transactionHandler
}

func New(u UserUC, m MainCategUC, s SubCategUC, t TransactionUC) *Handler {
	return &Handler{
		User:        *newUserHandler(u),
		MainCateg:   *newMainCategHandler(m),
		SubCateg:    *newSubCategHandler(s),
		Transaction: *newTransactionHandler(t),
	}
}
