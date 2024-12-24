package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	_ "github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	//解析模版
	t, err := template.ParseFiles("./template/hello.tmpl")
	if err != nil {
		fmt.Printf("解析模版失败，错误：%v\n", err)
		return
	}
	//渲染模版
	t.Execute(w, "lsl")
}

func structPractice(w http.ResponseWriter, r *http.Request) {
	//解析模版
	t, err := template.ParseFiles("./template/structPractice.tmpl")
	if err != nil {
		fmt.Printf("解析模版失败，错误：%v\n", err)
		return
	}
	users := []struct {
		Name string
		Age  int
	}{{"lsl", 18},
		{"gsy", 18},
		{"fym", 18}}
	//渲染模版
	t.Execute(w, users)
}

// 自定义函数
func customFunc(w http.ResponseWriter, r *http.Request) {
	//在解析模版之前将自定义函数添加进去
	kua := func(arg string) (string, error) {
		return arg + "真帅", nil
	}
	//解析模版
	t, err := template.New("customFunc.tmpl").Funcs(template.FuncMap{"kua": kua}).ParseFiles("./template/customFunc.tmpl")
	if err != nil {
		fmt.Printf("解析模版失败，错误：%v\n", err)
		return
	}
	user := struct {
		Name string
		Age  int
	}{"lsl", 18}
	//渲染模版
	t.Execute(w, user)
}

// 上面使用的是传统的net/http包创建服务
func main() {
	fileUpload()
}

/*
下面的函数使用Gin框架创建服务
*/

func HelloGin() {
	//创建一个默认的路由引擎
	r := gin.Default()
	//Get:请求方式 /hello：请求的路径
	r.GET("/hello", func(c *gin.Context) {
		//返回JSON格式的数据
		c.JSON(200, gin.H{"msg": "Hello world!"})
	})
	//启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()
}

func RESTfulAPI() {
	r := gin.Default()
	r.GET("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GET",
		})
	})

	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST",
		})
	})

	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT",
		})
	})

	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DELETE",
		})
	})
	r.Run()
}

// HTML渲染
func htmlRendering() {
	r := gin.Default()
	//解析模版前加载静态文件
	r.Static("/static", "./statics")
	//解析模版
	r.LoadHTMLGlob("./template/**/*")
	r.GET("/post/index", func(c *gin.Context) {
		//name指定了要渲染HTML模版文件的名称
		//gin.H是map[string]interface{}的别名，用于向模版传递数据
		c.HTML(http.StatusOK, "post/index.html", gin.H{
			"title": "post/index",
		})
	})
	r.GET("user/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/index.html", gin.H{
			"title": "user/index",
		})
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.Run(":8080")
}

// 自定义模版函数
func customFunction() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"safe": func(arg string) template.HTML {
			return template.HTML(arg)
		},
	})
	router.LoadHTMLFiles("./index.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", "<a href='https://liwenzhou.com'>李文周的博客</a>")
	})
	router.Run(":8080")
}

// 模版继承
// 加载模版文件
func loadTemplate(templateDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	layouts, err := filepath.Glob("/layouts/*.tmpl")
	if err != nil {
		panic(err)
	}
	includes, err := filepath.Glob("/includes/*.tmpl")
	if err != nil {
		panic(err)
	}
	//为layouts和includes目录生成template map
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
func indexFunc(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}
func homeFunc(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}
func inheritFunc() {
	r := gin.Default()
	//解析模版
	r.HTMLRender = loadTemplate("./template")
	//渲染模版
	r.GET("/index", indexFunc)
	r.GET("/home", homeFunc)
	r.Run()
}

// json渲染
func jsonRendering() {
	r := gin.Default()
	//没有模版直接进行渲染
	r.GET("/json", func(c *gin.Context) {
		user := struct {
			Name string
			Age  int
		}{"lsl", 18}
		c.JSON(http.StatusOK, user)
	})
	r.Run(":8080")
}

// 获取querystring参数
func getQuerystring() {
	r := gin.Default()
	r.GET("user/search", func(c *gin.Context) {
		username := c.DefaultQuery("username", "lsl")
		address := c.Query("address")
		c.JSON(http.StatusOK, gin.H{
			"msg":      "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

// 从form表单获取数据
func getParameterFromForm() {
	r := gin.Default()
	r.POST("/user/search", func(c *gin.Context) {
		username := c.PostForm("username")
		address := c.PostForm("address")
		c.JSON(http.StatusOK, gin.H{
			"msg":      "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

// 获取json格式的数据
func getJSONParameter() {
	r := gin.Default()
	r.POST("/user/search", func(c *gin.Context) {
		byte, _ := c.GetRawData() //从c.Request.Body读取请求数据
		//定义map
		var m map[string]interface{}
		// 反序列化
		_ = json.Unmarshal(byte, &m)
		c.JSON(http.StatusOK, m)
	})
	r.Run(":8080")
}

// ParameterBind 使用ShouldBind 进行参数绑定
func ParameterBind() {
	type Login struct {
		User     string `form:"user" json:"user" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	r := gin.Default()
	r.POST("/bind", func(c *gin.Context) {
		var para Login
		if err := c.ShouldBind(&para); err == nil {
			fmt.Printf("Login info:%v\n", para)
			c.JSON(http.StatusOK, gin.H{
				"user":     para.User,
				"password": para.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	// 绑定form表单示例 (user=q1mi&password=123456)
	r.POST("/form", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	// 绑定QueryString示例 (/loginQuery?user=q1mi&password=123456)
	r.GET("/loginForm", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.Run(":8080")
}

// 文件上传
func fileUpload() {
	r := gin.Default()
	//解析模版
	r.LoadHTMLFiles("./template/upload.html")
	//渲染模版
	r.GET("/upload.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})
	//上传文件请求
	r.POST("upload", func(c *gin.Context) {
		//单个文件
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Println(file.Filename)
		dir := fmt.Sprintf("E:/桌面文件/Goland/GinFramework/%s", file.Filename)
		//将文件上传到指定的目录
		c.SaveUploadedFile(file, dir)
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	r.Run(":8080")
}
