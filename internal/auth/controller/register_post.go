package controller

import (
	"fmt"
	tModels "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/auth/controller/models"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/global"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/logger"
	utils2 "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/utils"
	"net/http"
)

// AuthRegisterPost - ручка регистрации
func (handler *Handler) AuthRegisterPost(w http.ResponseWriter, r *http.Request) {
	op := "auth.controller.api_auth.AuthRegisterPost"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	ctx := r.Context()

	// Парсинг фронтовой модели данных регистрации
	var p tModels.Reg
	if err := utils2.ParseModels(r, &p, op); err != nil {
		// проставляем http.StatusBadRequest
		logger.StandardResponse(err.Error(), global.GetCodeError(err), r.Host, op)
		w.WriteHeader(global.GetCodeError(err))
		// отправляем структуру ошибки
		utils2.SendStringModel(&tModels.ModelError{Message: global.GetMsgError(err)}, w, op)
		return
	}

	// Валидация полей модели
	if _, err := p.Validate(); err != nil {
		logger.StandardWarnF(op, "Received validation error {%v}", err.Error())
		// проставляем http.StatusBadRequest
		logger.StandardResponse(err.Error(), global.GetCodeError(err), r.Host, op)
		w.WriteHeader(global.GetCodeError(err))
		// отправляем структуру ошибки
		utils2.SendStringModel(&tModels.ModelError{Message: global.GetMsgError(err)}, w, op)
		return
	}

	// Создание пользователя и генерация токена
	tokenString, err := handler.b.RegisterNewUser(ctx, p.Username, p.Password) // передаем username и password
	if err != nil {
		logger.StandardDebugF(op, "Received register error {%v}", err)
		// проставляем http.StatusBadRequest
		logger.StandardResponse(err.Error(), global.GetCodeError(err), r.Host, op)
		w.WriteHeader(global.GetCodeError(err))
		// отправляем структуру ошибки
		utils2.SendStringModel(&tModels.ModelError{Message: global.GetMsgError(err)}, w, op)
		return
	}

	// Сохранение токена в куки
	cookie := utils2.CreateCookie(tokenString)

	// Устанавливаем токен в куку
	http.SetCookie(w, &cookie)

	logger.StandardResponse(
		fmt.Sprintf("Successful register user=%v with token='%v'", p.Username, tokenString),
		http.StatusOK, r.Host, op)

	w.WriteHeader(http.StatusOK)
}
