package entity

type User struct {
	ID         uint   `json:"id" binding:"gte=0"`
	Name       string `json:"name" binding:"required"`
	Age        int    `json:"age" binding:"gte=0,lte=120"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
