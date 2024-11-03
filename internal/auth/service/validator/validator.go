package validator

import (
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/auth/config"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/global"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/logger"
	"github.com/pkg/errors"
	"regexp"
)

func validateUsername(username string) (string, error) {
	op := "internal.auth.service.validator.validateUsername"
	// Длина не менее 4 символов
	if len(username) < config.MinUsernameLen {
		logger.StandardDebugF(op, "Username is too short username=%v", username)
		return username, errors.Wrap(global.ErrNotValidUserAndPassword, op)
	}

	// Длина не более 10 символов
	if len(username) > config.MaxUsernameLen {
		return username, errors.Wrap(global.ErrNotValidUserAndPassword, op)
	}

	return username, nil
}

func validatePassword(password string) (string, error) {
	op := "internal.auth.service.validator.validatePassword"

	// Длина не меньше 8 символов
	if len(password) < config.MinPasswordLen {
		return password, errors.Wrap(global.ErrNotValidUserAndPassword, op)
	}

	// Длина не больше 64 символов
	if len(password) > config.MaxPasswordLen {
		return password, errors.Wrap(global.ErrNotValidUserAndPassword, op)
	}

	re := regexp.MustCompile(config.PatternPasswordMustContainDiffChar)
	if !re.MatchString(password) {
		return password, errors.Wrap(global.ErrNotValidUserAndPassword, op)
	}

	return password, nil
}

func ValidateUsernameAndPassword(username string, password string) (string, string, error) {
	op := "internal.auth.service.validator.ValidateUsernameAndPassword"

	username, err := validateUsername(username)
	if err != nil {
		return username, password, errors.Wrap(err, op)
	}

	password, err = validatePassword(password)
	if err != nil {
		return username, password, errors.Wrap(err, op)
	}
	return username, password, nil
}
