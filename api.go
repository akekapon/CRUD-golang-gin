package main

import (
	"net/http"

	"allinone.mod/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root@tcp(127.0.0.1:3306)/go_orm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// สร้างตัวเซิร์ฟเวอร์ Gin
	r := gin.Default()

	// กำหนดเส้นทาง (route) และการตอบกลับ
	r.GET("/users", func(c *gin.Context) {

		var users []model.User // สร้างตัวแปร users เป็น slice ของ model.User
		db.Find(&users)        // ค้นหาข้อมูลทั้งหมดในตาราง users
		c.JSON(200, users)
	})
	r.GET("/users/:id", func(c *gin.Context) { // รับค่า id จาก url
		id := c.Param("id")
		var user model.User
		db.First(&user, id)
		c.JSON(200, user)
	})
	r.POST("/users", func(c *gin.Context) { // รับค่าจาก client และแปลงเป็น JSON
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil { // รับค่าจาก client และแปลงเป็น JSON
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // ถ้ามี error ให้แสดงค่า error และจบการทำงาน
			return
		}

		result := db.Create(&user)                              // สร้างข้อมูลใหม่
		c.JSON(200, gin.H{"RowsAffected": result.RowsAffected}) //มีเอ๊ฟฟเฟคมั้ย        // แสดงข้อมูลที่สร้าง
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user model.User
		db.Delete(&user, id)
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "delete success"})
	})
	r.PUT("/users/:id", func(c *gin.Context) { // รับค่าจาก client และแปลงเป็น JSON
		id := c.Param("id")
		var user model.User
		db.First(&user, id)
		user.Fname = c.PostForm("fname")
		user.Lname = c.PostForm("lname")
		user.Username = c.PostForm("username")
		user.Avatar = c.PostForm("avatar")
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&user)
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "update success"})
	})

	r.Use(cors.Default()) // ใช้งาน cors
	// เริ่มต้นเซิร์ฟเวอร์
	r.Run(":8080")
}
