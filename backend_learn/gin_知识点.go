







通过 URL 路径传递的参数（如 /users/123），使用 c.Param("key")获取：
/users/:id   -----  c.Param("id")
支持多参数（如 /users/:id/:name


查询参数
通过 URL ?后传递的参数（如 /search?keyword=gin），常用方法：
c.Query("key")：获取参数，不存在返回空字符串。
c.DefaultQuery("key", "default")：参数不存在时返回默认值

r.GET("/search", func(c *gin.Context) {
    keyword := c.Query("keyword")
    page := c.DefaultQuery("page", "1") // 默认值 "1"
    c.JSON(200, gin.H{"keyword": keyword, "page": page})
})


获取数组参数：c.QueryArray("tags")。
获取 Map 参数：c.QueryMap("filters")


请求体数据（Request Body）

1. 表单数据
    适用于 application/x-www-form-urlencoded或 multipart/form-data：

r.POST("/submit", func(c *gin.Context) {
    name := c.PostForm("name")
    email := c.DefaultPostForm("email", "default@example.com")
    // 获取数组或 Map
    hobbies := c.PostFormArray("hobbies")
    c.JSON(200, gin.H{"name": name, "email": email})
})。    





直接读取请求体字节流
body, _ := c.GetRawData()
var data map[string]interface{}
json.Unmarshal(body, &data)。


JSON 数据
使用 ShouldBindJSON绑定到结构体
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}


文件上传（File Upload）
处理 multipart/form-data类型的文件
多文件上传：使用 c.MultipartForm()解析多个文件

r.POST("/upload", func(c *gin.Context) {
    file, header, err := c.Request.FormFile("file") // "file" 为表单字段名
    if err != nil {
        c.String(400, "上传失败")
        return
    }
    defer file.Close()
    // 保存文件（示例）
    out, _ := os.Create(header.Filename)
    io.Copy(out, file)
    c.String(200, "上传成功")
})。



Header 参数
通过 c.GetHeader("Key")或绑定到结构体
// 单参数获取
token := c.GetHeader("Authorization")

// 绑定到结构体
type Headers struct {
    Rate   int    `header:"Rate"`
    Domain string `header:"Domain"`
}
var h Headers
c.ShouldBindHeader(&h)




Gin 内置绑定方法 根据请求的 Content-Type 自动选择解析器（JSON/XML/Form 等
ShouldBind(&obj)：通用绑定，自动推断类型 自动根据 Content-Type绑定（支持 JSON/XML/表单等）
ShouldBindJSON(&obj)：仅绑定 JSON 数据
ShouldBindXML(&obj)：绑定 XML 数据
ShouldBindQuery(&obj)：仅绑定 URL Query 参数（适用于 GET 请求） 仅绑定查询参数
ShouldBindYAML(&obj)：绑定 YAML 数据

Gin 返回的均为字符串，需手动转换类型







c *gin.Context
DB.WithContext(ctx)
创建带有上下文的数据库实例，方便管理数据库操作的生命周期，上下文可以控制数据库操作：超时或取消，避免请求结束但数据库操作还在继续
链路追踪 - 在分布式系统中传递追踪信息，就类似传参