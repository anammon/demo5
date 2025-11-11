package controler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatrixRequest struct {
	Rows int         `json:"rows" binding:"required,min=1"`
	Cols int         `json:"cols" binding:"required,min=1"`
	Data [][]float64 `json:"data" binding:"required"`
}

type CombinedMatrixRequest struct {
	A MatrixRequest `json:"a" binding:"required"`
	B MatrixRequest `json:"b" binding:"required"`
}

func validateMatrixReq(m MatrixRequest) error {
	if m.Rows <= 0 || m.Cols <= 0 {
		return fmt.Errorf("矩阵维度无效: %dx%d", m.Rows, m.Cols)
	}
	if len(m.Data) != m.Rows {
		return fmt.Errorf("数据行数 %d 与指定的行数 %d 不匹配", len(m.Data), m.Rows)
	}
	for i := range m.Data {
		if len(m.Data[i]) != m.Cols {
			return fmt.Errorf("第 %d 行的列数 %d 与指定的列数 %d 不匹配", i, len(m.Data[i]), m.Cols)
		}
	}
	return nil
}

func addMatrices(a, b MatrixRequest) (MatrixRequest, error) {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		return MatrixRequest{}, fmt.Errorf("matrix dimensions mismatch: %dx%d vs %dx%d", a.Rows, a.Cols, b.Rows, b.Cols)
	}
	res := make([][]float64, a.Rows)
	for i := 0; i < a.Rows; i++ {
		res[i] = make([]float64, a.Cols)
		for j := 0; j < a.Cols; j++ {
			res[i][j] = a.Data[i][j] + b.Data[i][j]
		}
	}
	return MatrixRequest{Rows: a.Rows, Cols: a.Cols, Data: res}, nil
}

func subtractMatrices(a, b MatrixRequest) (MatrixRequest, error) {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		return MatrixRequest{}, fmt.Errorf("matrix维度不匹配: %dx%d vs %dx%d", a.Rows, a.Cols, b.Rows, b.Cols)
	}
	res := make([][]float64, a.Rows)
	for i := 0; i < a.Rows; i++ {
		res[i] = make([]float64, a.Cols)
		for j := 0; j < a.Cols; j++ {
			res[i][j] = a.Data[i][j] - b.Data[i][j]
		}
	}
	return MatrixRequest{Rows: a.Rows, Cols: a.Cols, Data: res}, nil
}

func multiplyMatrices(a, b MatrixRequest) (MatrixRequest, error) {
	if a.Cols != b.Rows {
		return MatrixRequest{}, fmt.Errorf("matrix不能相乘: %dx%d * %dx%d", a.Rows, a.Cols, b.Rows, b.Cols)
	}
	res := make([][]float64, a.Rows)
	for i := 0; i < a.Rows; i++ {
		res[i] = make([]float64, b.Cols)
		for j := 0; j < b.Cols; j++ {
			sum := 0.0
			for k := 0; k < a.Cols; k++ {
				sum += a.Data[i][k] * b.Data[k][j]
			}
			res[i][j] = sum
		}
	}
	return MatrixRequest{Rows: a.Rows, Cols: b.Cols, Data: res}, nil
}

func HandleMatrixAddition(c *gin.Context) {
	var req CombinedMatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.A); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.B); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := addMatrices(req.A, req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows": result.Rows, "cols": result.Cols, "data": result.Data})
}

func HandleMatrixSubtraction(c *gin.Context) {
	var req CombinedMatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.A); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.B); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := subtractMatrices(req.A, req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows": result.Rows, "cols": result.Cols, "data": result.Data})
}

func HandleMatrixMultiplication(c *gin.Context) {
	var req CombinedMatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.A); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMatrixReq(req.B); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := multiplyMatrices(req.A, req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows": result.Rows, "cols": result.Cols, "data": result.Data})
}
