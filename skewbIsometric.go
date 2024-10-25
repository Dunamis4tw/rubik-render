package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IsometricSkewb struct {
	Colors     map[Side][]rune             // Карта для хранения цветов каждой стороны
	SideParams map[Side]IsometricSkewbSide // Параметры боковой стороны скьюба
}

// Структура, хранящая параметры для построения элементов на стороне скьюба
type IsometricSkewbSide struct {
	Drawn [5]string // Атрибут d тега path в SVG, хранящий в себе построение фигуры
}

// ParseIsometricSkewbParams парсит параметры для изометрической SVG картинки скьюба
func ParseIsometricSkewbParams(pDimensions, pColors string) (IsometricSkewb, error) {

	// Извлечение размеров из строки pDimensions
	dimensions := strings.Split(pDimensions, "x")
	if len(dimensions) != 1 {
		return IsometricSkewb{}, fmt.Errorf("invalid dimensions: expected 1 dimensions")
	}

	dX, err1 := strconv.Atoi(dimensions[0])
	if err1 != nil {
		return IsometricSkewb{}, fmt.Errorf("invalid dimension values, expected integer values")
	}
	if dX < 1 || dX > 1 {
		return IsometricSkewb{}, fmt.Errorf("dimension values must be between 1 and 1")
	}

	// Инициализация структуры IsometricSkewb
	skewb := IsometricSkewb{
		Colors: make(map[Side][]rune),
	}

	Colors := strings.Split(strings.ToUpper(pColors), "-")

	// Функция для безопасного извлечения цвета или возвращения пустой строки
	getColorOrEmpty := func(index int, defaultColor string) string {
		if index < len(Colors) && len(Colors[index]) > 0 {
			return Colors[index]
		}
		return defaultColor
	}

	// Парсинг цветов для каждой стороны
	skewb.Colors[Front] = stringToRuneGrid(getColorOrEmpty(0, "X"), 5, 1)[0]
	skewb.Colors[Up] = stringToRuneGrid(getColorOrEmpty(1, "X"), 5, 1)[0]
	skewb.Colors[Right] = stringToRuneGrid(getColorOrEmpty(2, "X"), 5, 1)[0]

	// Цвет фона (base) будет последним в массиве Colors
	skewb.Colors[Base] = stringToRuneGrid(getColorOrEmpty(3, "K"), 1, 1)[0]

	return skewb, nil
}

// GenerateIsometricSkewb генерирует изометрическую SVG картинку скьюба
func GenerateIsometricSkewb(skewb IsometricSkewb) string {
	var builder strings.Builder

	// // // // // ПРОИЗВОДИМ РАСЧЁТЫ

	// Определяем элементы на сторонах (side) скьюба с размерами
	skewb.SideParams = map[Side]IsometricSkewbSide{
		Front: {
			Drawn: [5]string{
				"m18.5 84.3-8.5 4.9c-3.5 2-6.1-.5-6.1-4.6v-25.6c0-3.7 3-4.9 5.7-3.4l23.6 13.6c2.6 1.6 2.2 5.2-1.3 7.3z",
				"m83.6 105c0-3.7-2.2-8-5-9.6l-22.5-13c-4-2.3-6 1-3.9 4.7l25.4 43.9c2.1 3.7 6 3.7 6-.1v-25.9z",
				"m8.9 109.8c-1.8-3.1-4.9-3-4.9.7v27.1c0 3.7 2.2 8 5 9.6l23.1 13.4c2.8 1.6 6-.1 4.2-3.3z",
				"m83.6 154v29.7c0 3.7-2.3 5.4-5 3.8l-24-13.9c-2.8-1.6-1.7-4.6.6-5.9l23.2-13.4c2.1-1.2 5.2-.4 5.2 3.2z",
				"m80 141.5-33.5-58.1c-2.7-4.7-8.6-4.3-11.2-2.8l-23.5 13.6c-2.8 1.6-4.5 7-2.6 10.3l33.7 58.3c1.2 2.1 5 3.3 6.8 2.2l27.8-16c3.9-2.3 3.3-6.1 2.5-7.4z"},
		},
		Up: {
			Drawn: [5]string{
				"m41.1 62.3-0-26c0-4.3-3.4-5.9-6-4.4l-23.5 13.6c-2.8 1.6-2.4 4.4.8 6.2l22.9 13.2c3.2 1.8 5.8 1.1 5.8-2.6z",
				"m115.8 23.7c2.6 0 4-2.1.8-3.9l-24.6-14.2c-3.2-1.8-8-2-10.8-.5l-25.3 14.6c-2.6 1.5-1.1 3.9 1.3 3.9z",
				"m91.3 91.5 24.4-14.3c2.8-1.6 1.4-4.7-1.9-4.7l-55.3 0c-2.8 0-3.9 3.4-1.7 4.7l23.9 13.8c3.2 1.8 8 2 10.6.5z",
				"m135.9 65.6 25.1-14.4c2.8-1.6 2.4-4.4-.8-6.2l-23.6-13.7c-2.2-1.2-5.2-1.2-5.2 2.1l0 29.4c0 3.2 3.1 3.7 4.5 2.8z",
				"m126.2 62.7-.1-30.4c0-3.5-3-5.6-6.9-5.6l-66.1 0c-3.6 0-6.8 2.5-6.8 5.6l0 30.4c0 3.2 2.9 5.6 6.8 5.6l66.2 0c3.6 0 6.8-2.5 6.8-5.6z"},
		},
		Right: {
			Drawn: [5]string{
				"m89 106-0 25c-0 3.7 4.3 3.7 6.4-.1l25.3-43.9c1.8-3.2.2-6.2-2.5-4.6l-24.2 14c-2.8 1.6-5 5.9-5 9.6z",
				"m139.7 76.3 21.7 12.4c4.2 2.4 7.3.7 7.3-3v-25.7c0-3.7-2.3-5.4-5-3.8l-24.1 13.9c-2.8 1.6-3.1 4.4.1 6.2z",
				"m89 157.5v27.2c0 3.7 2.3 5.4 5 3.8l25.8-14.9c2.4-1.4 1.6-4.5-.8-5.9l-23.1-13.3c-2.1-1.2-6.9-.4-6.9 3.2z",
				"m136.9 156.9c-1.8 3.2 2.2 5.5 5 3.9l21.8-12.6c2.8-1.6 5-5.9 5-9.6v-28c0-3.7-2.7-4.1-4.7-.6z",
				"m93 141.4 33.5-58.1c2.7-4.7 8.6-4.3 11.2-2.8l23.5 13.6c2.8 1.6 4.5 7 2.6 10.3l-33.7 58.3c-1.2 2.1-5 3.3-6.8 2.2l-27.8-16c-3.9-2.3-3.3-6.1-2.5-7.4z"},
		},
		Base: {
			Drawn: [5]string{"M172.57 138.48v-82.33a15 15 0 00-7.49-13l-71.31 -41.17a14.94 14.94 0 00-15 0l-71.28 41.17a15 15 0 00-7.49 13v82.33a15 15 0 007.49 13l71.31 41.17a15 15 0 0015 0l71.28 -41.17a15 15 0 007.49-13z"},
		},
	}

	// // // // // СТРОИМ SVG

	// Создаём рамку (viewBox)
	builder.WriteString("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 172.57 194.62\">")

	// Создаём основу (base)
	colorBase := skewb.Colors[Base][0]
	builder.WriteString(fmt.Sprintf("\r\n\t<path id=\"base\" d=\"%s\" style=\"fill: %s\"/>",
		skewb.SideParams[Base].Drawn[0], colorMapRGBA[colorBase]))

	// Создаём стороны (side)
	GenerateIsometricSkewbSide(&builder, skewb, Front)
	GenerateIsometricSkewbSide(&builder, skewb, Up)
	GenerateIsometricSkewbSide(&builder, skewb, Right)

	// Закрываем рамку (viewBox)
	builder.WriteString("\r\n</svg>")

	// Возвращаем сгенерированную SVG
	return builder.String()
}

func GenerateIsometricSkewbSide(builder *strings.Builder, skewb IsometricSkewb, side Side) {
	sideParam := skewb.SideParams[side]

	// Начало группы
	builder.WriteString(fmt.Sprintf("\r\n\t<g id=\"%s\">", side.String()))
	for x := 0; x < len(skewb.Colors[side]); x++ {
		color := skewb.Colors[side][x]
		path := fmt.Sprintf("\r\n\t\t<path id=\"%c-%d\" d=\"%s\" style=\"fill: %s\"/>",
			side.String()[0], x+1, sideParam.Drawn[x], colorMapRGBA[color])
		builder.WriteString(path)
	}
	// Закрытие группы
	builder.WriteString("\r\n\t</g>")
}
