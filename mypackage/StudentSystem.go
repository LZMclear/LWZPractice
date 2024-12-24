package mypackage

import "fmt"

type Address struct {
	Province string
	City     string
}
type Student struct {
	Id   int
	Name string
	Age  int
	Address
}
type Class struct {
	Id       int
	Students []Student
}

func Init() *Class {
	//初始化班级
	c := &Class{
		Id:       1,
		Students: make([]Student, 0, 10),
	}
	for i := 0; i < 5; i++ {
		stu := Student{
			Id:   i,
			Name: fmt.Sprintf("stu%d", i),
			Age:  i + 10,
			Address: Address{
				"河南",
				"新乡",
			},
		}
		c.Students = append(c.Students, stu)
	}
	fmt.Println("班级初始化")
	return c
}

// 添加学生信息，传入的参数是学生
func (c *Class) AddStudent(s *Student) *Class {
	c.Students = append(c.Students, *s)
	fmt.Sprintf("添加学生后的学生信息为：", c)
	return c
}
