package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IsometricCube struct {
	Size       Size                            // Размер Кубика Рубика XYZ
	Colors     map[Side][][]rune               // Карта для хранения цветов каждой стороны
	SideParams map[Side]IsometricSideParameter // Параметры боковой стороны кубика
}

// Структура, хранящая параметры для построения элементов на стороне кубика
type IsometricSideParameter struct {
	Base   Point  // Базовая точка
	Multi  Point  // X отвечает за горизонтальный шаг, Y — за вертикальный (0 по оси Y)
	Offset Point  // Смещение: 0 по X и шаг вниз по Y
	Drawn  string // Атрибут d тега path в SVG, хранящий в себе построение фигуры
}

// ParseIsometricParams парсит параметры для изометрической SVG картинки кубика
func ParseIsometricParams(pDimensions, pColors string) (IsometricCube, error) {

	// Извлечение размеров из строки pDimensions
	dimensions := strings.Split(pDimensions, "x")
	if len(dimensions) != 3 {
		return IsometricCube{}, fmt.Errorf("invalid dimensions: expected 3 dimensions")
	}

	dX, err1 := strconv.Atoi(dimensions[0])
	dY, err2 := strconv.Atoi(dimensions[1])
	dZ, err3 := strconv.Atoi(dimensions[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return IsometricCube{}, fmt.Errorf("invalid dimension values, expected integer values")
	}
	if dX < 1 || dY < 1 || dZ < 1 || dX > 64 || dY > 64 || dZ > 64 {
		return IsometricCube{}, fmt.Errorf("dimension values must be between 1 and 64")
	}

	Colors := strings.Split(strings.ToUpper(pColors), "-")

	// Инициализация структуры IsometricCube
	cube := IsometricCube{
		Size:   Size{X: dX, Y: dY, Z: dZ},
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
	cube.Colors[Up] = stringToRuneGrid(getColorOrEmpty(1, "X"), dZ, dX)
	cube.Colors[Right] = stringToRuneGrid(getColorOrEmpty(2, "X"), dZ, dY)

	// Цвет фона (base) будет последним в массиве Colors
	cube.Colors[Base] = stringToRuneGrid(getColorOrEmpty(3, "K"), 1, 1)

	return cube, nil
}

// GenerateIsometricCube генерирует изометрическую SVG картинку кубика
func GenerateIsometricCube(cube IsometricCube) string {
	var builder strings.Builder

	// Получаем размерность куба (в float64)
	dX := float64(cube.Size.X)
	dY := float64(cube.Size.Y)
	dZ := float64(cube.Size.Z)

	// // // // // ПРОИЗВОДИМ РАСЧЁТЫ

	// Считаем размер рамки (viewBox)
	viewBoxSize := Point{
		X: 2.85 + 42.43*(dZ+dX),
		Y: -1.38 + 24.5*(dZ+dX) + 49*dY,
	}

	// Считаем координаты точек, по которым рисуется основа (base)
	LX := Point{X: 13.55 - 42.43*dX, Y: 7.83 - 24.5*dX}
	LY := Point{Y: 15.67 - 49*dY}
	LZ := Point{X: 13.58 - 42.43*dZ, Y: -7.83 + 24.5*dZ}
	M := Point{X: 2.85 + 42.43*(dX+dZ), Y: -8.52 + 49*dY + 24.5*dX}

	// Считаем положение элементов на сторонах (side) кубика с размерами XxYxZ
	cube.SideParams = map[Side]IsometricSideParameter{
		Front: {
			Base:   Point{X: 41.2, Y: 31.48 + 24.5*dZ},
			Multi:  Point{X: 0, Y: 24.5},
			Offset: Point{X: 42.43, Y: 49},
			Drawn:  "v29.69c0,3.67-2.25,5.37-5,3.78l-27.23-15.72c-2.75-1.59-5-5.9-5-9.56v-29.69c0-3.67 2.25-5.37 5-3.78l27.23 15.72c2.75 1.6 5.01 5.9 5.01 9.57z",
		},
		Up: {
			Base:   Point{X: 48.83, Y: 17.92 + 24.5*dZ},
			Multi:  Point{X: 42.47, Y: -24.5},
			Offset: Point{X: 42.43, Y: 24.5},
			Drawn:  "l27.23-15.72c2.75-1.59 2.4-4.39-.78-6.23l-25.7-14.84c-3.18-1.84-8-2-10.79-.45l-27.23 15.72c-2.75 1.59-2.4 4.39.78 6.23l25.71 14.84c3.17 1.84 8.02 2.04 10.78.45z",
		},
		Right: {
			Base:   Point{X: 4.09 + 42.43*dX, Y: 7.98 + 24.5*(dZ+dX)},
			Multi:  Point{X: 0, Y: -24.5},
			Offset: Point{X: 42.43, Y: 49},
			Drawn:  "v29.69c0 3.66 2.25 5.37 5 3.78l27.23-15.72c2.76-1.59 5-5.9 5-9.56v-29.73c0-3.67-2.25-5.37-5-3.78l-27.23 15.72c-2.77 1.6-5 5.89-5 9.6z",
		},
	}

	// // // // // СТРОИМ SVG

	// Создаём рамку (viewBox)
	builder.WriteString(fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %.2f %.2f\">",
		viewBoxSize.X, viewBoxSize.Y))

	// Создаём основу (base)
	colorBase := cube.Colors[Base][0][0]
	builder.WriteString(fmt.Sprintf("\r\n\t<path id=\"base\" d=\"M%.2f %.2fv%.2fa15 15 0 00-7.49-13l%.2f %.2fa14.94 14.94 0 00-15 0l%.2f %.2fa15 15 0 00-7.49 13v%.2fa15 15 0 007.49 13l%.2f %.2fa15 15 0 0015 0l%.2f %.2fa15 15 0 007.49-13z\" style=\"fill: %s\"/>",
		M.X, M.Y, LY.Y, LX.X, LX.Y, LZ.X, LZ.Y, -LY.Y, -LX.X, -LX.Y, -LZ.X, -LZ.Y, colorMapRGBA[colorBase]))

	// Создаём стороны (side)
	GenerateIsometricSide(&builder, cube, Front)
	GenerateIsometricSide(&builder, cube, Up)
	GenerateIsometricSide(&builder, cube, Right)

	// Закрываем рамку (viewBox)
	builder.WriteString("\r\n</svg>")

	// Возвращаем сгенерированную SVG
	return builder.String()
}

func GenerateIsometricSide(builder *strings.Builder, cube IsometricCube, side Side) {
	sideParam := cube.SideParams[side]

	// Начало группы
	builder.WriteString(fmt.Sprintf("\r\n\t<g id=\"%s\">", side.String()))
	for x := 0; x < len(cube.Colors[side]); x++ {
		for y := 0; y < len(cube.Colors[side][x]); y++ {
			color := cube.Colors[side][x][y]
			startX := sideParam.Base.X + float64(x)*sideParam.Multi.X + float64(y)*sideParam.Offset.X
			startY := sideParam.Base.Y + float64(y)*sideParam.Multi.Y + float64(x)*sideParam.Offset.Y

			path := fmt.Sprintf("\r\n\t\t<path id=\"%c-%dx%d\" d=\"M%.2f %.2f %s\" style=\"fill: %s\"/>",
				side.String()[0], x+1, y+1, startX, startY, sideParam.Drawn, colorMapRGBA[color])
			builder.WriteString(path)
		}
	}
	// Закрытие группы
	builder.WriteString("\r\n\t</g>")
}
