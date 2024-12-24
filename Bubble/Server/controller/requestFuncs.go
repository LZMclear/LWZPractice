package controller

import (
	"Server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

var db *sqlx.DB

type Todo struct {
	ID     int
	Title  string `form:"title" json:"title"`
	Status bool   //两种状态：待做 0，已做 1，这个应该由后端生成
}

// ConnectSql 连接数据库中间件
func ConnectSql(c *gin.Context) {
	//计算运行时间
	start := time.Now()
	db = utils.InitDB()
	c.Next()
	cost := time.Since(start)
	fmt.Printf("程序执行时间：%v\n", cost)
}

// HandlerPost 增添数据，接收的参数是json格式的
func HandlerPost(c *gin.Context) {
	//获取参数
	var todo Todo
	err := c.ShouldBind(&todo)
	fmt.Printf("%v", todo)
	//增添数据，所以状态都是未做，状态设置为0
	todo.Status = false
	sqlxStr := "insert into todo (title, status) values (?,?)"
	result, err := db.Exec(sqlxStr, todo.Title, todo.Status)
	if err != nil {
		fmt.Printf("insert data failed, err: %v\n")
		return
	}
	n, err := result.LastInsertId()
	if err != nil {
		fmt.Println("获取插入数据id错误")
		return
	}
	fmt.Printf("insert success, the affected row is:%d\n", n)
	c.JSON(http.StatusOK, gin.H{
		"message": "insert success",
	})
}

// HandlerGet 查询数据
func HandlerGet(c *gin.Context) {
	//查询全部数据
	sqlxStr := "select * from todo where id >?"
	//定义一个切片存储数据
	var todos []Todo
	err := db.Select(&todos, sqlxStr, 0)
	if err != nil {
		fmt.Printf("查询全部数据失败：%v\n", err)
		return
	}
	c.JSON(http.StatusOK, todos)
	fmt.Println("查询成功")
}

// HandlerUpdate 修改数据
func HandlerUpdate(c *gin.Context) {
	sqlxStr := "update todo set status = ? where id = ?"
	var todo Todo
	err := c.ShouldBind(&todo)
	if err != nil {
		fmt.Printf("绑定数据错误：%v\n", err)
		return
	}
	fmt.Println(todo)
	status := todo.Status
	id := todo.ID
	result, err := db.Exec(sqlxStr, status, id)
	if err != nil {
		fmt.Printf("修改数据失败：%v\n", err)
		return
	}
	n, err := result.RowsAffected()
	fmt.Printf("更新了第%d条数据", n)
	//将修改后的数据返回
	sqlStr2 := "select * from todo where id=?"
	var todoed Todo
	_ = db.Get(&todoed, sqlStr2, todo.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
		"result":  todoed,
	})
}

// HandlerDelete 删除数据
func HandlerDelete(c *gin.Context) {
	//获取要删除的id
	var todo Todo
	err := c.ShouldBind(&todo)
	if err != nil {
		fmt.Printf("绑定数据出错：%v\n", err)
		return
	}
	id := todo.ID
	sqlxStr := "delete from todo where id=?"
	result, _ := db.Exec(sqlxStr, id)
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("获取被删除id错误:%v\n", err)
		return
	}
	c.JSON(http.StatusOK, "第"+string(n)+"条数据被删除了")
}
