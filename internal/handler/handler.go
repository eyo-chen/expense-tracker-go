package handler

type Handler struct {
	User      userHandler
	MainCateg mainCategHandler
}

func New(u UserUC, m MainCategUC) *Handler {
	return &Handler{
		User:      *newUserHandler(u),
		MainCateg: *newMainCategHandler(m),
	}
}
