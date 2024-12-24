package controller

import (
	"StressTesting/entity"
	"StressTesting/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var db *sqlx.DB
var Trans ut.Translator

// ConnectMysql 连接数据库中间件
func ConnectMysql(c *gin.Context) {
	db = utils.InitDB()
	c.Next()
	//执行完毕关闭数据库
	db.Close()
}

// InitTranslator 初始化翻译器
func InitTranslator() {
	var err error
	Trans, err = utils.InitTrans("zh")
	if err != nil {
		fmt.Println("获取翻译器失败")
		return
	}
}

func SaveUserHandler(c *gin.Context) {
	InitTranslator()
	var user entity.User
	err := c.ShouldBind(&user)
	if err != nil {
		//断言错误类型是否属于validator.ValidationErrors类型，并获取错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok { //不是validator类型的错误，无法解析，直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(Trans),
		})
		return
	}
	sqlxStr := "insert into user(name,age,password,re_password) values (:name,:age,:password,:repassword)"
	result, err := db.NamedExec(sqlxStr, user)
	if err != nil {
		fmt.Printf("执行插入操作失败：%v\n", err)
		return
	}
	n, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("获取插入数据id错误：%v\n", err)
		return
	}
	fmt.Printf("插入数据成功，最新的数据id为%v\n", n)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": user,
	})
}
