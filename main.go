package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
)

type Size struct {
	X, Y, Z int
}

type Point struct {
	X, Y float64
}

// Определяем структуру для параметров командной строки
type Options struct {
	IP       string `short:"i" long:"ip" description:"IP address to listen on" default:"0.0.0.0"`
	Port     string `short:"p" long:"port" description:"Port to start the server" default:"80"`
	CertFile string `short:"c" long:"cert" description:"Path to SSL certificate"`
	KeyFile  string `short:"k" long:"key" description:"Path to SSL key"`
	Help     bool   `short:"h" long:"help" description:"Display help information"`
}

func main() {
	var opts Options
	// Парсим параметры командной строки
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// Обработка флага help
	if opts.Help {
		PrintHelp()
		os.Exit(0)
	}

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/cube/:view/:dimensions/:colors", CubeHandler)
		v1.GET("/skewb/:view/:dimensions/:colors", SkewbHandler)
	}

	// Формирование адреса для прослушивания
	address := fmt.Sprintf("%s:%s", opts.IP, opts.Port)

	// Сообщение в консоль
	log.Printf("Listening on %s:%s", opts.IP, opts.Port)

	// Проверяем, указаны ли сертификаты для SSL
	if opts.CertFile != "" && opts.KeyFile != "" {
		log.Printf("Open https://%s:%s in your browser", opts.IP, opts.Port)
		// Запуск HTTPS сервера с сертификатами
		err := router.RunTLS(address, opts.CertFile, opts.KeyFile)
		if err != nil {
			log.Fatalf("Error starting HTTPS server: %v", err)
		}
	} else {
		log.Printf("Open http://%s:%s in your browser", opts.IP, opts.Port)
		// Запуск HTTP сервера
		err := router.Run(address)
		if err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}
}

// PrintHelp выводит информацию о параметрах программы
func PrintHelp() {
	fmt.Println("Program usage:")
	fmt.Println("  -i, --ip       IP address to listen on (default 0.0.0.0)")
	fmt.Println("  -p, --port     Port to start the server (default 80)")
	fmt.Println("  -c, --cert     Path to SSL certificate")
	fmt.Println("  -k, --key      Path to SSL key")
	fmt.Println("  -h, --help     Display this help")
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
	case "unfolded":
		// Парсим параметры
		unfoldedCube, err := ParseUnfoldedParams(pDimensions, pColors)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Генерация SVG
		svg := GenerateUnfoldedCube(unfoldedCube)

		// Установка заголовков и вывод SVG
		c.Header("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown view parameter"})
	}
}

// SkewbHandler обрабатывает запросы для генерации SVG Скьюба
func SkewbHandler(c *gin.Context) {
	// Получение параметров из URL
	pDimensions := c.Param("dimensions")
	pView := c.Param("view")
	pColors := c.Param("colors")

	switch pView {
	case "isometric":
		// Парсим параметры
		isometricCube, err := ParseIsometricSkewbParams(pDimensions, pColors)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Генерация SVG
		svg := GenerateIsometricSkewb(isometricCube)

		// Установка заголовков и вывод SVG
		c.Header("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown view parameter"})
	}
}
