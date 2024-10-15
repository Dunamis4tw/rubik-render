package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Size struct {
	X, Y, Z int
}

type Point struct {
	X, Y float64
}

func main() {
	router := gin.Default()
	router.GET("/cube/:view/:dimensions/:colors", CubeHandler)
	router.Run(":8080")
}

// CubeHandler обрабатывает запросы для генерации SVG кубика Рубика
func CubeHandler(c *gin.Context) {
	// Получение параметров из URL
	pDimensions := c.Param("dimensions")
	pView := c.Param("view")
	pColors := c.Param("colors")

	switch pView {
	case "isometric":
		// Парсим параметры
		isometricCube, err := ParseIsometricParams(pDimensions, pColors)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Генерация SVG
		svg := GenerateIsometricCube(isometricCube)

		// Установка заголовков и вывод SVG
		c.Header("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
		return
	case "flat":
		// Парсим параметры
		flatCube, err := ParseFlatParams(pDimensions, pColors)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Генерация SVG
		svg := GenerateFlatCube(flatCube)

		// Установка заголовков и вывод SVG
		c.Header("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown view parameter"})
	}
}
