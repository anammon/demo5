package controler

import (
	"backend/model"
	"hash"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserControler struct {
}

func (u UserControler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	if user.Account == "" || user.Password == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "账户、密码和邮箱为必填字段",
		})
		return
	}
	var existingUser model.User
	if err := model.DB.Where("account = ?", user.Account).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "账户名已被注册",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := model.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "邮箱已被注册",
		})
		return
	}

	// 检查用户名是否已存在
	if err := model.DB.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "用户名已被使用",
		})
		return
	}
	hashpwd,err3:=bcrypt.GenerateFromPassword(
		[]byte(user.Password),bcrypt.DefaultCost
	)
	if err3!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"failed to hashpwd"
		})
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册失败: " + err.Error(),
		})
		return
	}
   
	// 返回成功响应（不返回密码）
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":       user.ID,
			"account":  user.Account,
			"email":    user.Email,
			"name":     user.Name,
			"birthday": user.Birthday.Format("2006-01-02"),
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
	c.JSON(200, gin.H{
		"message": "Registration successful",
		"user":    user,
	})
}
func (user UserControler) Login(c *gin.Context) {

}
