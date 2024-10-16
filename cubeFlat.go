package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Side перечисление для сторон кубика
type Side int

const (
	Front Side = iota
	Left
	Right
	Top
	Down
	Back
	Base
)

func (s Side) String() string {
	return [...]string{"front", "left", "right", "top", "down", "back", "base"}[s]
}

type FlatCube struct {
	Size       Size                    // Размер Кубика Рубика XYZ
	Colors     map[Side][][]rune       // Карта для хранения цветов каждой стороны
	SideParams map[Side]FlatFaceConfig // Параметры боковой стороны кубика
}

type FlatFaceConfig struct {
	Base  Point // Базовая точка
	Size  Point // Размер объекта
	Multi Point // X отвечает за горизонтальный шаг, Y — за вертикальный (0 по оси Y)
}

// ParseFlatParams парсит параметры для плоской SVG картинки кубика
func ParseFlatParams(pDimensions, pColors string) (FlatCube, error) {

	// Извлечение размеров из строки pDimensions
	dimensions := strings.Split(pDimensions, "x")
	if len(dimensions) != 2 {
		return FlatCube{}, fmt.Errorf("invalid dimensions: expected 2 dimensions")
	}

	dX, err1 := strconv.Atoi(dimensions[0])
	dY, err2 := strconv.Atoi(dimensions[1])
	if err1 != nil || err2 != nil {
		return FlatCube{}, fmt.Errorf("invalid dimension values, expected integer values")
	}
	if dX < 1 || dY < 1 || dX > 64 || dY > 64 {
		return FlatCube{}, fmt.Errorf("dimension values must be between 1 and 64")
	}
	Colors := strings.Split(strings.ToUpper(pColors), "-")

	// Инициализация структуры FlatCube с использованием карты для хранения цветов
	cube := FlatCube{
		Size:   Size{X: dX, Y: dY},
		Colors: make(map[Side][][]rune),
	}

	// Функция для безопасного извлечения цвета или возвращения пустой строки
	getColorOrEmpty := func(index int, defaultColor string) string {
		if index < len(Colors) {
			return Colors[index]
		}
		return defaultColor
	}

	// Парсинг цветов для каждой стороны
	cube.Colors[Front] = stringToRuneGrid(getColorOrEmpty(0, "X"), dX, dY)
	cube.Colors[Left] = stringToRuneGrid(getColorOrEmpty(1, "T"), 1, dY)
	cube.Colors[Top] = stringToRuneGrid(getColorOrEmpty(2, "T"), dX, 1)
	cube.Colors[Right] = stringToRuneGrid(getColorOrEmpty(3, "T"), 1, dY)
	cube.Colors[Down] = stringToRuneGrid(getColorOrEmpty(4, "T"), dX, 1)

	// Цвет фона (base) будет последним в массиве Colors
	cube.Colors[Base] = stringToRuneGrid(getColorOrEmpty(5, "K"), 1, 1)

	return cube, nil
}

// GenerateFlatCube генерирует плоскую SVG картинку кубика
func GenerateFlatCube(cube FlatCube) string {
	var builder strings.Builder

	// // // // // ПРОИЗВОДИМ РАСЧЁТЫ

	cube.SideParams = map[Side]FlatFaceConfig{
		Left: {
			Base: Point{X: 0, Y: 10},
			Size: Point{X: 6, Y: 43},
		},
		Right: {
			Base: Point{X: 0, Y: 10},
			Size: Point{X: 6, Y: 43},
		},
		Down: {
			Base: Point{X: 10, Y: 0},
			Size: Point{X: 43, Y: 6},
		},
		Top: {
			Base: Point{X: 10, Y: 0},
			Size: Point{X: 43, Y: 6},
		},
	}

	// // // // // СТРОИМ SVG

	// Основной SVG-код для кубика Рубика
	mianSVG := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %d %d\">",
		14+cube.Size.X*49, 14+cube.Size.Y*49)
	builder.WriteString(mianSVG)

	colorBase := cube.Colors[Base][0][0]
	baseRect := fmt.Sprintf("<rect id=\"base\" width=\"%d\" height=\"%d\" rx=\"7.42\" style=\"fill: %s\" x=\"3\" y=\"3\"/>",
		8+cube.Size.X*49, 8+cube.Size.Y*49, colorMapRGBA[colorBase])
	builder.WriteString(baseRect)

	// Генерация фронтальной стороны
	builder.WriteString("\r\n\t<g id=\"front\">")
	for y := 0; y < len(cube.Colors[Front]); y++ {
		for x := 0; x < len(cube.Colors[Front][y]); x++ {
			color := cube.Colors[Front][y][x]

			startX := 10 + x*49
			startY := 10 + y*49

			path := fmt.Sprintf("\r\n\t\t<rect id=\"%s-%dx%d\" x=\"%d\" y=\"%d\" width=\"43\" height=\"43\" rx=\"6.21\" style=\"fill: %s\"/>",
				"front", x+1, y+1, startX, startY, colorMapRGBA[color])
			builder.WriteString(path)
		}
	}
	// Закрытие группы
	builder.WriteString("\r\n\t</g>")

	// Генерация остальных сторон
	builder.WriteString(GenerateFlatSide(cube, Left))
	builder.WriteString(GenerateFlatSide(cube, Top))
	builder.WriteString(GenerateFlatSide(cube, Right))
	builder.WriteString(GenerateFlatSide(cube, Down))

	// Закрытие SVG
	builder.WriteString("\r\n</svg>")

	// Возвращаем финальную строку SVG
	return builder.String()
}

// GenerateFlatSide генерирует сторону кубика
func GenerateFlatSide(cube FlatCube, side Side) string {
	var builder strings.Builder
	sideParam := cube.SideParams[side]

	// Начало группы
	builder.WriteString(fmt.Sprintf("\r\n\t<g id=\"%s\">", side.String()))
	for y := 0; y < len(cube.Colors[side]); y++ {
		for x := 0; x < len(cube.Colors[side][y]); x++ {
			colorRune := cube.Colors[side][y][x]

			startX := int(sideParam.Base.X) + x*49
			startY := int(sideParam.Base.Y) + y*49

			if side == Right {
				startX += 8 + cube.Size.X*49
			} else if side == Down {
				startY += 8 + cube.Size.Y*49
			}

			rect := fmt.Sprintf("\r\n\t\t<rect id=\"%s-%dx%d\" x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"2.32\" style=\"fill: %s\"/>",
				side.String(), x+1, y+1, startX, startY, int(sideParam.Size.X), int(sideParam.Size.Y), colorMapRGBA[colorRune])
			builder.WriteString(rect)
		}
	}
	// Закрытие группы
	builder.WriteString("\r\n\t</g>")

	return builder.String()
}
