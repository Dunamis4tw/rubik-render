package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseUnfoldedParams парсит параметры для развёртки SVG картинки кубика
func ParseUnfoldedParams(pDimensions, pColors string) (FlatCube, error) {

	// Извлечение размеров из строки pDimensions
	dimensions := strings.Split(pDimensions, "x")
	if len(dimensions) != 3 {
		return FlatCube{}, fmt.Errorf("invalid dimensions: expected 3 dimensions")
	}

	dX, err1 := strconv.Atoi(dimensions[0])
	dY, err2 := strconv.Atoi(dimensions[1])
	dZ, err3 := strconv.Atoi(dimensions[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return FlatCube{}, fmt.Errorf("invalid dimension values, expected integer values")
	}
	if dX < 1 || dY < 1 || dZ < 1 || dX > 64 || dY > 64 || dZ > 64 {
		return FlatCube{}, fmt.Errorf("dimension values must be between 1 and 64")
	}
	Colors := strings.Split(strings.ToUpper(pColors), "-")

	// Инициализация структуры FlatCube с использованием карты для хранения цветов
	cube := FlatCube{
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
	cube.Colors[Left] = stringToRuneGrid(getColorOrEmpty(1, "X"), dZ, dY)
	cube.Colors[Up] = stringToRuneGrid(getColorOrEmpty(2, "X"), dX, dZ)
	cube.Colors[Right] = stringToRuneGrid(getColorOrEmpty(3, "X"), dZ, dY)
	cube.Colors[Down] = stringToRuneGrid(getColorOrEmpty(4, "X"), dX, dZ)
	cube.Colors[Back] = stringToRuneGrid(getColorOrEmpty(5, "X"), dX, dY)

	// Цвет фона (base) будет последним в массиве Colors
	cube.Colors[Base] = stringToRuneGrid(getColorOrEmpty(6, "K"), 1, 1)

	return cube, nil
}

// GenerateUnfoldedCube генерирует развёртку SVG картинку кубика
func GenerateUnfoldedCube(cube FlatCube) string {
	var builder strings.Builder

	// Радиус скругления
	const round = 7.42

	// Получаем размерность куба (в float64)
	dX := float64(cube.Size.X)
	dY := float64(cube.Size.Y)
	dZ := float64(cube.Size.Z)

	// Расчитывем длины сторон
	lX := 8 + dX*49
	lY := 8 + dY*49
	lZ := 8 + dZ*49

	// Расчитывем длины сторон за вычетом скруглений
	// Вычитание 7 нужно, чтобы убрать большой отступ между сторонами
	lrX := lX - 2*round - 7
	lrY := lY - 2*round - 7
	lrZ := lZ - 2*round - 7

	// // // // // ПРОИЗВОДИМ РАСЧЁТЫ

	cube.SideParams = map[Side]FlatSideParameter{
		Front: {
			Base: Point{X: lZ - 7, Y: lZ - 7},
		},
		Left: {
			Base: Point{X: 0, Y: lZ - 7},
		},
		Up: {
			Base: Point{X: lZ - 7, Y: 0},
		},
		Right: {
			Base: Point{X: lZ + lX - 7*2, Y: lZ - 7},
		},
		Down: {
			Base: Point{X: lZ - 7, Y: lZ + lY - 7*2},
		},
		Back: {
			Base: Point{X: lZ*2 + lX - 7*3, Y: lZ - 7},
		},
	}

	// // // // // СТРОИМ SVG

	// Основной SVG-код для кубика Рубика
	mianSVG := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %d %d\">",
		int(lZ*2+lX*2-7*3), int(lZ*2+lY-7*2))
	builder.WriteString(mianSVG)

	// Генерируем фон
	colorBase := cube.Colors[Base][0][0]
	baseRect := fmt.Sprintf("\r\n\t<path id=\"base\" d=\"M7.42 %.2fh%.2fs7.42 0 7.42 -7.42v%.2fs0 -7.42 7.42 -7.42h%.2fs7.42 0 7.42 7.42"+
		"v%.2fs0 7.42 7.42 7.42h%.2fs7.42 0 7.42 7.42v%.2fs 0 7.42 -7.42 7.42h%.2fs-7.42 0 -7.42 7.42v%.2fs0 7.42 -7.42 7.42h%.2f"+
		"s-7.42 0 -7.42 -7.42v%.2fs0 -7.42 -7.42 -7.42h%.2fs-7.42 0 -7.42 -7.42v%.2fs0 -7.42 7.42 -7.42\" style=\"fill: %s\"/>",
		lZ-7, lrZ, -lrZ, lrX+7, lrZ, lrZ+lrX+2*round, lrY+7, -lrZ-lrX-2*round, lrZ, -lrX-7, -lrZ, -lrZ, -lrY-7, colorMapRGBA[colorBase])
	builder.WriteString(baseRect)

	// Генерация остальных сторон
	GenerateUnfoldedSide(&builder, cube, Front)
	GenerateUnfoldedSide(&builder, cube, Left)
	GenerateUnfoldedSide(&builder, cube, Up)
	GenerateUnfoldedSide(&builder, cube, Right)
	GenerateUnfoldedSide(&builder, cube, Down)
	GenerateUnfoldedSide(&builder, cube, Back)

	// Закрытие SVG
	builder.WriteString("\r\n</svg>")

	// Возвращаем финальную строку SVG
	return builder.String()
}

func GenerateUnfoldedSide(builder *strings.Builder, cube FlatCube, side Side) {
	startBase := cube.SideParams[side].Base

	builder.WriteString("\r\n\t<g id=\"" + side.String() + "\">")
	for y := 0; y < len(cube.Colors[side]); y++ {
		for x := 0; x < len(cube.Colors[side][y]); x++ {
			color := cube.Colors[side][y][x]

			startX := int(startBase.X) + 7 + x*49
			startY := int(startBase.Y) + 7 + y*49

			path := fmt.Sprintf("\r\n\t\t<rect id=\"%c-%dx%d\" x=\"%d\" y=\"%d\" width=\"43\" height=\"43\" rx=\"6.21\" style=\"fill: %s\"/>",
				side.String()[0], x+1, y+1, startX, startY, colorMapRGBA[color])
			builder.WriteString(path)
		}
	}

	// Закрытие группы
	builder.WriteString("\r\n\t</g>")
}
