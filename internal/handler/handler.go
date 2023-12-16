package handler

type Handler struct {
	User      userHandler
	MainCateg mainCategHandler
	SubCateg  subCategHandler
}

func New(u UserUC, m MainCategUC, s SubCategUC) *Handler {
	return &Handler{
		User:      *newUserHandler(u),
		MainCateg: *newMainCategHandler(m),
		SubCateg:  *newSubCategHandler(s),
	}
}
