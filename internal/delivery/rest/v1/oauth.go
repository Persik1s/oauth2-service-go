package v1

import (
	"github.com/Persik1s/oauth2-service-go/internal/usecase"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type OAuthDelivery struct {
	logger *slog.Logger

	oauthUsecase usecase.OAuthUsecaseIn
}

func NewOAuthDelivery(gr *echo.Group, oauthUsecase usecase.OAuthUsecaseIn, logger *slog.Logger) *OAuthDelivery {
	oauth := &OAuthDelivery{
		logger:       logger,
		oauthUsecase: oauthUsecase,
	}

	gr.GET("/oauth/callback", oauth.OCallBack)

	return oauth
}

func (o *OAuthDelivery) OCallBack(ctx echo.Context) error {
	code := ctx.QueryParam("code")

	token, err := o.oauthUsecase.OAuthorization(ctx.Request().Context(), code)
	if err != nil {
		o.logger.Debug("o.oauthUsecase.OAuthorization", "err", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, token)
}
