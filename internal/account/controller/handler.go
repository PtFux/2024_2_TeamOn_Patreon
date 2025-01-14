package controller

import "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/account/controller/interfaces"

type Handler struct {
	serv interfaces.AccountService
}

func New(service interfaces.AccountService) *Handler {
	return &Handler{serv: service}
}
