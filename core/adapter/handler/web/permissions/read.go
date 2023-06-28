package permissions

import (
	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

func (ctrl Controller) GetPermissions(ctx *gin.Context) {
	ctrl.log.Info("GetPermissions")

	authPermissions := ctx.MustGet("permissions").([]gocloak.ResourcePermission)
	res := system.NewHttpResponse(true, "User Permissions", gin.H{
		"permissions": gin.H{
			"auth": authPermissions,
		},
	})

	ctx.JSON(200, res)
}
