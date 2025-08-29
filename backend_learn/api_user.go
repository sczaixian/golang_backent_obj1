package backendlearn


/*
1 它允许在闭包创建时捕获外部变量（如svcCtx），
	这样可以在多个请求间共享这些变量而  《不需要每次请求都重新初始化》。
2 这种方式提供了更好的封装性，可以将依赖项（如服务上下文）注入到处理器中。
3 它支持更灵活的函数签名，允许在闭包创建时执行一些初始化逻辑
*/
//  注册路由
user.GET("/:address/login-message", v1.GetLoginMessageHandler(svcCtx)) // 生成login签名信息
func GetLoginMessageHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Params.ByName("address")
		if address == "" {
			xhttp.Error(c, errcode.NewCustomErr("user addr is null"))
			return
		}

		res, err := service.GetUserLoginMsg(c.Request.Context(), svcCtx, address)
		if err != nil {
			xhttp.Error(c, errcode.NewCustomErr(err.Error()))
			return
		}

		xhttp.OkJson(c, res)
	}
}

依赖注入​​：
需传递配置、数据库连接、服务对象等​​非请求级依赖​​时，闭包是标准解决方案

通过闭包​​隔离依赖与业务逻辑​​，适合需要动态注入配置、服务对象等场景，提升代码可维护性和可测试性
返回HandlerFunc支持依赖注入、闭包和代码复用，符合软件设计原则如单一职责和开放封闭

返回HandlerFunc是一种函数式编程的模式，它支持闭包和配置封装。这使得代码更模块化、可测试和可重用。
例如，我们可以将依赖注入（如svcCtx）到处理程序中，而不需要全局变量


但在Gin中，路由处理程序必须匹配HandlerFunc类型，即func(*gin.Context)。
所以，我们不能直接传递svcCtx给处理程序，除非通过闭包或全局状态。
因此，返回HandlerFunc是必要的，以注入依赖