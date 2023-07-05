package middlewares

import (
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

func Authorization(log *zap.Logger, config *internal.AppConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("Authorization middleware")

		// skip authorization for login and register
		if strings.Contains(ctx.Request.URL.Path, "/auth/user/login") || strings.Contains(ctx.Request.URL.Path, "/auth/user/signup") || strings.Contains(ctx.Request.URL.Path, "/auth/user/refresh-token") || strings.Contains(ctx.Request.URL.Path, "/health") {
			ctx.Next()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Error("Authorization header is empty")

			apiError := system.NewHttpResponse(false, "Authorization header is empty", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		accessToken := strings.Split(authHeader, " ")[1]
		if accessToken == "" {
			log.Error("Access token is empty")

			apiError := system.NewHttpResponse(false, "Access token is empty", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		keycloakServer := config.KeycloakServer
		keycloakRealm := config.KeycloakRealm
		keycloakClientId := config.KeycloakClientId
		keycloakClientSecret := config.KeycloakClientSecret

		client := gocloak.NewClient(keycloakServer)

		rtResult, err := client.RetrospectToken(ctx, accessToken, keycloakClientId, keycloakClientSecret, keycloakRealm)
		if err != nil {
			log.Error("Retrospect token error: ", zap.Error(err))

			apiError := system.NewHttpResponse(false, "Retrospect token error", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		isValidToken := *rtResult.Active
		if !isValidToken {
			log.Error("Token is not valid")

			apiError := system.NewHttpResponse(false, "Token is not valid", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		// decode access token
		user, err := client.GetUserInfo(ctx, accessToken, keycloakRealm)
		if err != nil {
			log.Error("Decode access token error: ", zap.Error(err))

			apiError := system.NewHttpResponse(false, "Decode access token error", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		// Requesting Part Token
		rpt, err := client.GetRequestingPartyToken(ctx, accessToken, keycloakRealm, gocloak.RequestingPartyTokenOptions{
			Audience: &keycloakClientId,
		})
		if err != nil {
			log.Error("Get requesting party token error: ", zap.Error(err))

			apiError := system.NewHttpResponse(false, "Get requesting party token error", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		rptResult, err := client.RetrospectToken(ctx, rpt.AccessToken, keycloakClientId, keycloakClientSecret, keycloakRealm)
		if err != nil {
			log.Error("Retrospect token error: ", zap.Error(err))

			apiError := system.NewHttpResponse(false, "Retrospect token error", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		// check if permissions are there in the token and if the user has the permission to access the resource
		if rptResult.Permissions == nil {
			log.Error("No permissions found in the token")

			apiError := system.NewHttpResponse(false, "No permissions found in the token", nil)
			ctx.AbortWithStatusJSON(401, apiError)
			return
		}

		serviceName := config.ServiceName
		method := getScopeFromReq(ctx.Request)

		log.Debug("serviceName: ", zap.String("serviceName", serviceName))
		log.Debug("method: ", zap.String("method", method))

		for _, permission := range *rptResult.Permissions {
			resourceName := *permission.RSName
			log.Debug("resourceName: ", zap.String("resourceName", resourceName))

			if resourceName == serviceName {
				log.Info("Permission granted for user: ", zap.String("userId", *user.Sub), zap.String("method", method), zap.String("resourceName", resourceName))
				for _, scope := range *permission.Scopes {
					if scope == method {

						ctx.Set("user", user)
						ctx.Set("userId", *user.Sub)
						ctx.Set("accessToken", accessToken)
						ctx.Set("permissions", *rptResult.Permissions)

						ctx.Next()
						return
					}
				}

				log.Error("Permission denied")
				apiError := system.NewHttpResponse(false, "Permission denied", nil)
				ctx.AbortWithStatusJSON(401, apiError)
				return
			}
		}

		log.Error("Permission denied")
		apiError := system.NewHttpResponse(false, "Permission denied", nil)
		ctx.AbortWithStatusJSON(401, apiError)
		return
	}
}
