package main

import "fmt"

// RGBAtoHex преобразует RGBA значения в шестнадцатеричный формат
func RGBAtoHex(r, g, b, a int) string {
	if a == 0 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else {
		return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
	}
}

var colorMapRGBA = map[rune]string{
	'R': RGBAtoHex(213, 0, 0, 0), // Красный
	// 'R': RGBAtoHex(234, 6, 0, 0), // Красный
	'G': RGBAtoHex(0, 153, 0, 0), // Зеленый
	// 'G': RGBAtoHex(7, 164, 46, 0),    // Зеленый
	'B': RGBAtoHex(52, 52, 212, 0), // Синий
	// 'B': RGBAtoHex(24, 12, 138, 0),   // Синий
	'Y': RGBAtoHex(255, 255, 0, 0), // Желтый
	// 'Y': RGBAtoHex(255, 241, 68, 0),  // Желтый
	'W': RGBAtoHex(223, 223, 223, 0), // Белый
	'O': RGBAtoHex(239, 108, 0, 0),   // Оранжевый
	// 'O': RGBAtoHex(255, 164, 13, 0),  // Оранжевый
	'X': RGBAtoHex(86, 86, 86, 0), // Серый
	'K': RGBAtoHex(0, 0, 0, 0),    // Чёрный
	'T': "transparent",            // Прозрачный
}

// var colorMapRGBAPastel = map[rune]string{
// 	'R': RGBAtoHex(228, 113, 122, 0), // Красный
// 	'G': RGBAtoHex(62, 180, 137, 0),  // Зеленый
// 	'B': RGBAtoHex(175, 218, 252, 0), // Синий
// 	'Y': RGBAtoHex(252, 232, 131, 0), // Желтый
// 	'W': RGBAtoHex(253, 244, 227, 0), // Белый
// 	'O': RGBAtoHex(239, 169, 74, 0),  // Оранжевый
// 	'X': RGBAtoHex(128, 128, 128, 0), // Серый
// 	'K': RGBAtoHex(48, 48, 48, 0),    // Чёрный
// 	'T': "transparent",               // Прозрачный
// }

func stringToRuneGrid(input string, x, y int) [][]rune {
	// Создаём двумерный массив
	grid := make([][]rune, y)

	// Если строка состоит из одного символа, заполняем весь массив этим символом и возвращаем его
	if len(input) == 1 {
		fillRune := rune(input[0])
		for i := 0; i < y; i++ {
			grid[i] = make([]rune, x)
			for j := 0; j < x; j++ {
				grid[i][j] = fillRune
			}
		}
		return grid
	}

	// Преобразуем строку в срез рун
	runes := []rune(input)

	// Количество необходимых символов
	neededRunes := x * y

	// Если длина строки меньше, чем необходимо, добавляем 'X'
	if len(runes) < neededRunes {
		diff := neededRunes - len(runes)
		runes = append(runes, make([]rune, diff)...)
		for i := len(runes) - diff; i < len(runes); i++ {
			runes[i] = 'X'
		}
	} else if len(runes) > neededRunes {
		runes = runes[:neededRunes]
	}

	// Заполняем двумерный массив
	for i := 0; i < y; i++ {
		grid[i] = runes[i*x : (i+1)*x]
	}

	return grid
}
