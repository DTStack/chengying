# swagger

[Swagger Guide参考地址](https://promacanthus.netlify.app/tools/swagger/swagger-guide/)
[swagger github注释说明](https://github.com/swaggo/swag/blob/master/README_zh-CN.md)

- 生成doc
```
cd matrix
swag init  -o swagger/doc
```

- 访问

```
http://localhost:8864/swagger/
```

- swagger 注释模板

参数说明
```
// @Description  	操作行为的详细说明
// @Summary      	该操作的简短摘要
// @Tags         	每个API操作的标签列表，以逗号分隔。
// @Produce      	API可以生成的MIME类型的列表     // https://github.com/swaggo/swag/blob/master/README_zh-CN.md#mime%E7%B1%BB%E5%9E%8B
// @Accept          API可以生成的MIME类型的列表
// @Param           enumstring  query     string     false  "string enums"       Enums(A, B, C)
// @Param           enumint     query     int        false  "int enums"          Enums(1, 2, 3)
// @Param           enumnumber  query     number     false  "int enums"          Enums(1.1, 1.2, 1.3)
// @Param           string      query     string     false  "string valid"       minlength(5)  maxlength(10)
// @Param           int         query     int        false  "int valid"          minimum(1)    maximum(10)
// @Param           default     query     string     false  "string default"     default(A)
// @Param           collection  query     []string   false  "string collection"  collectionFormat(multi)
// @Param           extensions  query     []string   false  "string collection"  extensions(x-example=test,x-nullable)
// @Success         以空格分隔的成功响应。return code,{param type},data type,comment
// @Failure         以空格分隔的成功响应。return code,{param type},data type,comment
// @Router          以空格分隔的路径定义。 path,[httpMethod]
```

案例

get


```
// CheckPwdConnect  通过password密码检查ssh连通性检查
// @Summary      	ssh 密码连通性测试
// @Description  	通过password密码检查ssh连通性检查
// @Tags         	agent
// @Produce      	json
// @Accept          application/json
// @Param			agent_id query string true "Agent ID"
// @Success         200 {string} string "{"msg": "hello Razeen"}"    // code,{param type},data type,comment
// @Failure         400 {string} string "{"msg": "who are you"}"
// @Router          /agent/install/pwdconnect [post]

// Auth godoc
// @Summary      Auth admin
// @Description  get admin info
// @Tags         accounts,admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Admin
// @Failure      400  {object}  httputil.HTTPError
// @Failure      401  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /admin/auth [post]


// AttributeExample godoc
// @Summary      attribute example
// @Description  attribute
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        enumstring  query     string  false  "string enums"    Enums(A, B, C)
// @Param        enumint     query     int     false  "int enums"       Enums(1, 2, 3)
// @Param        enumnumber  query     number  false  "int enums"       Enums(1.1, 1.2, 1.3)
// @Param        string      query     string  false  "string valid"    minlength(5)  maxlength(10)
// @Param        int         query     int     false  "int valid"       minimum(1)    maximum(10)
// @Param        default     query     string  false  "string default"  default(A)
// @Success      200         {string}  string  "answer"
// @Failure      400         {string}  string  "ok"
// @Failure      404         {string}  string  "ok"
// @Failure      500         {string}  string  "ok"
// @Router       /examples/attribute [get]
```

post
```
// PostExample godoc
// @Summary      post request example
// @Description  post request example
// @Accept       json
// @Produce      plain
// @Param        message  body      model.Account  true  "Account Info"
// @Success      200      {string}  string         "success"
// @Failure      500      {string}  string         "fail"
// @Router       /examples/post [post]
```
