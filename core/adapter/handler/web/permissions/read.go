package permissions

import (
	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

func (ctrl Controller) GetPermissions(ctx *gin.Context) {
	ctrl.log.Info("GetPermissions")

	systemPermissions := ctx.MustGet("permissions").([]gocloak.ResourcePermission)

	permissionRes := gin.H{
		"permissions": gin.H{
			"system": systemPermissions,
		},
	}

	message := "User Permissions"
	res := system.NewHttpResponse(true, message, permissionRes)
	ctx.JSON(200, res)
}
