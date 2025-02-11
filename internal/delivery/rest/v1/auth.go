package v1

import (
	"errors"
	"github.com/Persik1s/oauth2-service-go/internal/dto"
	"github.com/Persik1s/oauth2-service-go/internal/errorz"
	"github.com/Persik1s/oauth2-service-go/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type AuthDelivery struct {
	logger    *slog.Logger
	validator *validator.Validate

	userUsecase usecase.AuthUsecaseIn
}

func NewAuthDelivery(gr *echo.Group, userUsecase usecase.AuthUsecaseIn, validator *validator.Validate, logger *slog.Logger) *AuthDelivery {
	auth := &AuthDelivery{
		logger:      logger,
		validator:   validator,
		userUsecase: userUsecase,
	}

	gr.POST("/auth/signup", auth.SignUp)
	gr.POST("/auth/signin", auth.SignIn)

	return auth
}

/*
Status 500 - StatusInternalServerError
Status 400 - StatusBadRequest
Status 409 - StatusConflict
Status 201 - StatusCreated
*/
func (o *AuthDelivery) SignUp(ctx echo.Context) error {
	data := dto.SignUpRequestDto{}
	if err := ctx.Bind(&data); err != nil {
		o.logger.Debug("ctx bind", "err", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if err := o.validator.Struct(data); err != nil {
		o.logger.Debug("o.validator.Struct", "err", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	respone, err := o.userUsecase.Registration(ctx.Request().Context(), data)
	if err != nil {
		if errors.Is(err, errorz.ErrUserAlreadyExists) {
			return ctx.NoContent(http.StatusConflict)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusCreated, respone)
}

/*
Status 500 - StatusInternalServerError
Status 400 - StatusBadRequest
Status 403 - StatusForbidden
Status 200 - StatusOK
*/
func (a *AuthDelivery) SignIn(ctx echo.Context) error {
	data := dto.SignInRequestDto{}
	if err := ctx.Bind(&data); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if err := a.validator.Struct(data); err != nil {
		a.logger.Debug("ctx bind", "err", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	token, err := a.userUsecase.Authorization(ctx.Request().Context(), data)
	if err != nil {
		if errors.Is(err, errorz.ErrPasswordNotValid) {
			return ctx.NoContent(http.StatusForbidden)
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, token)
}

//func (o *OAuthDelivery) Authorization(ctx echo.Context) error {
//	// scope
//	// client_id
//	// redirect_url
//	// respone_type
//
//	// name
//	// password
//
//	// redis database code
//	return ctx.NoContent(http.StatusOK)
//}

// scope
// respone_type
// client_id
// redirect_url

// name
// password

// /authorization code - 1
// scope: admin, user

// /token
