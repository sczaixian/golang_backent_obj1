package V1

import (
	"github.com/gin-gonic/gin"

	"github.com/ProjectsTask/Backend/service/svc"
	"github.com/ProjectsTask/SwapBase/errcode"
	"github.com/ProjectsTask/SwapBase/xhttp"
)

// user.GET("/:address/login-message", v1.GetLoginMessageHandler(svcCtx))
func UserLoginHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Params.ByName("address")
		if address == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}
	}
}
