package handler

type Handler struct {
	User userHandler
}

func New(userUC UserUC) *Handler {
	return &Handler{
		User: *newUserHandler(userUC),
	}
}
