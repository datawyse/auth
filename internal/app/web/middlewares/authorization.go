package middlewares

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
	"strings"

	"auth/core/domain/system"
	"auth/internal"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getScopeFromReq(req *http.Request) string {
	// map GET / read, POST / create, PUT / update, DELETE / delete
	switch req.Method {
	case http.MethodGet:
		return "read"
	case http.MethodPost:
		return "create"
	case http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	}
	return "read"
}

func Authorization(log *otelzap.Logger, config *internal.AppConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("Authorization middleware")

		// skip authorization for login and register
		if strings.Contains(ctx.Request.URL.Path, "/test") || strings.Contains(ctx.Request.URL.Path, "/user/login") || strings.Contains(ctx.Request.URL.Path, "/user/signup") || strings.Contains(ctx.Request.URL.Path, "/user/refresh-token") || strings.Contains(ctx.Request.URL.Path, "/health") {
			ctx.Next()
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			_ = ctx.Error(system.ErrInvalidAuthorizationToken)
			ctx.Next()
			return
		}

		accessToken := strings.Split(authHeader, " ")[1]
		if accessToken == "" {
			_ = ctx.Error(system.ErrInvalidToken)
			ctx.Next()
			return
		}

		keycloakServer := config.KeycloakServer
		keycloakRealm := config.KeycloakRealm
		keycloakClientId := config.KeycloakClientId
		keycloakClientSecret := config.KeycloakClientSecret

		client := gocloak.NewClient(keycloakServer)

		rtResult, err := client.RetrospectToken(ctx, accessToken, keycloakClientId, keycloakClientSecret, keycloakRealm)
		if err != nil {
			_ = ctx.Error(system.ErrRetrospectToken)
			ctx.Next()
			return
		}

		isValidToken := *rtResult.Active
		if !isValidToken {
			_ = ctx.Error(system.ErrInvalidToken)
			ctx.Next()
			return
		}

		// decode access token
		user, err := client.GetUserInfo(ctx, accessToken, keycloakRealm)
		if err != nil {
			_ = ctx.Error(system.ErrInvalidToken)
			ctx.Next()
			return
		}

		// Requesting Part Token
		rpt, err := client.GetRequestingPartyToken(ctx, accessToken, keycloakRealm, gocloak.RequestingPartyTokenOptions{
			Audience: &keycloakClientId,
		})
		if err != nil {
			_ = ctx.Error(system.ErrInvalidToken)
			ctx.Next()
			return
		}

		rptResult, err := client.RetrospectToken(ctx, rpt.AccessToken, keycloakClientId, keycloakClientSecret, keycloakRealm)
		if err != nil {
			_ = ctx.Error(system.ErrRetrospectToken)
			ctx.Next()
			return
		}

		// check if permissions are there in the token and if the user has the permission to access the resource
		if rptResult.Permissions == nil {
			_ = ctx.Error(system.ErrPermissionDenied)
			ctx.Next()
			return
		}

		serviceName := config.ServiceName
		method := getScopeFromReq(ctx.Request)

		log.Debug("serviceName: ", zap.String("serviceName", serviceName), zap.String("method", method))

		for _, permission := range *rptResult.Permissions {
			resourceName := *permission.RSName
			log.Debug("resourceName: ", zap.String("resourceName", resourceName))

			if resourceName == serviceName {
				log.Debug("Permission granted for user: ", zap.String("userId", *user.Sub), zap.String("method", method), zap.String("resourceName", resourceName))
				for _, scope := range *permission.Scopes {
					if scope == method {

						ctx.Set("user", user)
						ctx.Set("userId", *user.Sub)
						ctx.Set("accessToken", accessToken)
						ctx.Set("permissions", *rptResult.Permissions)

						ctx.Next()
					}
				}

				_ = ctx.Error(system.ErrPermissionDenied)
				ctx.Next()
				return
			}
		}

		_ = ctx.Error(system.ErrPermissionDenied)
		ctx.Next()
		return
	}
}
