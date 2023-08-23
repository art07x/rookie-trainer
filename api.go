package main

import (
	"example2/rookie/model"
	"regexp"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	//connect to database
	dsn := "host=localhost user=postgres password=1234 dbname=rookie port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.User{})

	r := gin.Default()

	//create user by form data
	r.POST("/users/form", func(c *gin.Context) {
		var user model.User
		user.Name = c.PostForm("name")
		ageStr := c.PostForm("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.JSON(500, gin.H{"error": "Invalid age value"})
			return
		}
		user.Age = age
		if user.Age < 1 || user.Age > 100 {
			c.JSON(500, gin.H{"error": "out of valid age range"})
			return
		}

		user.Year_of_birth = 2023 - user.Age
		user.Note = c.PostForm("note")

		user.Email = c.PostForm("email")

		//if email is not unique return error
		if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil {
			c.JSON(401, gin.H{"error": "Email already exists"})
			return
		}

		//check email by regular expression if not match return error
		if !regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(user.Email) {
			c.JSON(401, gin.H{"error": "Invalid email by regular expression"})
			return
		}

		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(500, gin.H{"error": "File not found!"})
			return
		}
		if file.Header["Content-Type"][0] != "image/png" && file.Header["Content-Type"][0] != "image/jpeg" {
			c.JSON(401, gin.H{"error": "File type not supported!"})
			return
		}
		if err := c.SaveUploadedFile(file, "avatar/"+file.Filename); err != nil {
			c.JSON(500, gin.H{"error": "File not saved!"})
			return
		}
		if file.Header["Content-Type"][0] == "image/png" {
			user.Avatar_type = "png"
			//delete string last 4 characters
			user.Avatar_name = file.Filename[:len(file.Filename)-4]
		}
		if file.Header["Content-Type"][0] == "image/jpeg" {
			user.Avatar_type = "jpg"
			user.Avatar_name = file.Filename[:len(file.Filename)-4]
		}

		db.Create(&user)
		c.JSON(201, user)

	})

	//get by id
	r.GET("/users/:id", func(c *gin.Context) {
		var user model.User
		if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found!"})
			return
		}
		c.JSON(200, user)

		//send response

	})

	//update user by id
	r.PUT("/users/:id", func(c *gin.Context) {
		var user model.User
		if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "User not found!"})
			return
		}
		user.Name = c.PostForm("name")
		ageStr := c.PostForm("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			c.JSON(500, gin.H{"error": "Invalid age value"})
			return
		}
		user.Age = age
		if user.Age < 1 || user.Age > 100 {
			c.JSON(500, gin.H{"error": "out of valid age range"})
			return
		}

		user.Year_of_birth = 2023 - user.Age
		user.Note = c.PostForm("note")
		if user.Note == "clean" {
			user.Note = ""
		}

		//if email  == "" not update email
		if c.PostForm("email") == "" {

		} else {
			user.Email = c.PostForm("email")
			if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil {
				c.JSON(500, gin.H{"error": "Email already exists"})
				return
			}
			if !regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(user.Email) {
				c.JSON(500, gin.H{"error": "Invalid email by regular expression"})
				return
			}

		}

		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(500, gin.H{"error": "File not found!"})
			return
		}
		if file.Header["Content-Type"][0] != "image/png" && file.Header["Content-Type"][0] != "image/jpeg" {
			c.JSON(500, gin.H{"error": "File type not sup	ported!"})
			return
		}
		if err := c.SaveUploadedFile(file, "avatar/"+file.Filename); err != nil {
			c.JSON(500, gin.H{"error": "File not saved!"})
			return
		}
		if file.Header["Content-Type"][0] == "image/png" {
			user.Avatar_type = "png"
			//delete string last 4 characters
			user.Avatar_name = file.Filename[:len(file.Filename)-4]
		}
		if file.Header["Content-Type"][0] == "image/jpeg" {
			user.Avatar_type = "jpg"
			user.Avatar_name = file.Filename[:len(file.Filename)-4]
		}

		db.Save(&user)
		c.JSON(200, user)
	})

	//delete user by id
	r.DELETE("/users/:id", func(c *gin.Context) {
		var user model.User
		if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "User not found!"})
			return
		}
		db.Delete(&user)
		c.JSON(200, gin.H{"success": "true"})
	})

	//get all by limit and offset
	r.GET("/userse", func(c *gin.Context) {
		var users []model.User
		limitStr := c.Query("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return
		}
		offsetStr := c.Query("offset")

		offset, err := strconv.Atoi(offsetStr)
		if offset <= 0 {
			c.JSON(400, gin.H{"error": "offset Invalid offset value"})
			return
		}
		if offsetStr == "1" {
			offsetStr = "0"
		}
		if offset > 0 {
			offset = offset * limit
			offset = offset - limit
		}

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return
		}
		//count all users
		var count int64
		db.Model(&model.User{}).Count(&count)

		db.Limit(limit).Offset(offset).Order("created_at asc").Find(&users)

		c.JSON(200, gin.H{"data": users, "count": count})
	})

	r.Run(":9000")

}
