package controler

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppControler struct {
}

func (AppControler) Create(c *gin.Context) {
	var app model.App
	if err := c.ShouldBind(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := model.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, app)
}
func (AppControler) GetApp(c *gin.Context) {
	var applist []model.App
	query := model.DB.Model(&model.App{}) //构建查询
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if description := c.Query("description"); description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}
	if err := query.Find(&applist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, applist)
}
