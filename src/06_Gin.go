package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Student 结构体成员变量名称首字母要大写，否则ShouldBindJSON失效
type Student struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"gte=1"`
}

func getting(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(student.Name)
	fmt.Println(student.Age)
	c.JSON(200, student)
}

func main() {
	router := gin.Default()
	// 路由分组
	v1 := router.Group("/v1")
	{
		v1.POST("/someGet", getting)
	}
	v2 := router.Group("/v2")
	{
		v2.POST("/someGet", getting)
	}
	// 设置监听端口并启动
	err := router.Run(":9090")
	if err != nil {
		return
	}
}
