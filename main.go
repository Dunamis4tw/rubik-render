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
	IP       string `short:"i" long:"ip" description:"IP адрес для прослушивания" default:"0.0.0.0"`
	Port     string `short:"p" long:"port" description:"Порт для запуска сервера" default:"80"`
	CertFile string `short:"c" long:"cert" description:"Путь к SSL сертификату"`
	KeyFile  string `short:"k" long:"key" description:"Путь к SSL ключу"`
	Help     bool   `short:"h" long:"help" description:"Вывести информацию о параметрах"`
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
	}

	// Формирование адреса для прослушивания
	address := fmt.Sprintf("%s:%s", opts.IP, opts.Port)

	// Сообщение в консоль
	log.Printf("Сервер запущен на %s:%s", opts.IP, opts.Port)
	log.Printf("Откройте http://%s:%s в браузере", opts.IP, opts.Port)

	// Проверяем, указаны ли сертификаты для SSL
	if opts.CertFile != "" && opts.KeyFile != "" {
		// Запуск HTTPS сервера с сертификатами
		err := router.RunTLS(address, opts.CertFile, opts.KeyFile)
		if err != nil {
			log.Fatalf("Ошибка при запуске HTTPS сервера: %v", err)
		}
	} else {
		// Запуск HTTP сервера
		err := router.Run(address)
		if err != nil {
			log.Fatalf("Ошибка при запуске HTTP сервера: %v", err)
		}
	}
}

// PrintHelp выводит информацию о параметрах программы
func PrintHelp() {
	fmt.Println("Использование программы:")
	fmt.Println("  -i, --ip       IP адрес для прослушивания (по умолчанию 0.0.0.0)")
	fmt.Println("  -p, --port     Порт для запуска сервера (по умолчанию 80)")
	fmt.Println("  -c, --cert     Путь к SSL сертификату")
	fmt.Println("  -k, --key      Путь к SSL ключу")
	fmt.Println("  -h, --help     Вывести эту справку")
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
