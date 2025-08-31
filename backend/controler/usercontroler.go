package controler

import (
	"backend/model"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserControler struct {
}
type RegisterDTO struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Birthday string `json:"birthday" binding:"required"` // 前端传 YYYY-MM-DD
}

func (u UserControler) Register(c *gin.Context) {
	var dto RegisterDTO
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	birthday, _ := time.Parse("2006-01-02", dto.Birthday)
	user := model.User{
		Account:   dto.Account,
		Password:  dto.Password,
		Email:     dto.Email,
		Name:      dto.Name,
		Birthday:  birthday, // 直接存 time.Time
		CreatedAt: time.Now(),
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
	hashpwd, err3 := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost, // 添加逗号
	)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hashpwd",
		})
		return // 添加return
	}
	user.Password = string(hashpwd) // 修正变量名
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
			"id":         user.ID,
			"account":    user.Account,
			"email":      user.Email,
			"name":       user.Name,
			"birthday":   user.Birthday.Format("2006-01-02"),
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
	// 删除重复的c.JSON调用
}
func (u UserControler) Login(c *gin.Context) { // 修正参数名
	var loginData struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	var user model.User
	if err := model.DB.Where("account = ?", loginData.Identifier).Or("email = ?", loginData.Identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "账户或邮箱输入错误",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "密码错误",
		})
		return
	}
	token, err := utils.GenerateJWT(user.Account) // 添加utils包名
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return // 添加return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"account":    user.Account,
			"email":      user.Email,
			"name":       user.Name,
			"birthday":   user.Birthday.Format("2006-01-02"),
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
