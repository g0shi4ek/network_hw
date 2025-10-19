package main

import (
	"fmt"
	"math/bits"
)

const n = 7
const k = 4
const informationVector = 11 // 1011b

// Функции для работы с бинарными числами
func IntToBits(num uint64, length int) []byte {
	bits := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		bits[i] = byte(num & 1)
		num >>= 1
	}
	return bits
}

func BitsToInt(bits []byte) uint64 {
	var result uint64
	for _, bit := range bits {
		result = (result << 1) | uint64(bit)
	}
	return result
}

// Кодирование Хэмминга
func encodeHamming(data uint64) uint64 {
	d := IntToBits(data, 4)

	// Вычисляем проверочные биты
	p1 := d[0] ^ d[1] ^ d[3] // 1^0^1^1 = 0
	p2 := d[0] ^ d[2] ^ d[3] // 1^1^1^1 = 1
	p4 := d[1] ^ d[2] ^ d[3] // 0^1^1^1 = 0

	// Формируем кодовое слово: P1,P2,D1,P4,D2,D3,D4
	encoded := []byte{p1, p2, d[0], p4, d[1], d[2], d[3]}

	return BitsToInt(encoded)
}

// Декодирование Хэмминга с исправлением ошибок
func decodeHamming(encoded uint64) uint64 {
	bits := IntToBits(encoded, 7)

	// Вычисляем синдром
	s1 := bits[0] ^ bits[2] ^ bits[4] ^ bits[6]
	s2 := bits[1] ^ bits[2] ^ bits[5] ^ bits[6]
	s3 := bits[3] ^ bits[4] ^ bits[5] ^ bits[6]

	syndrome := (uint64(s3) << 2) | (uint64(s2) << 1) | uint64(s1)

	// Исправляем ошибку если нужно
	if syndrome != 0 {
		errorPos := syndrome - 1
		if errorPos < uint64(len(bits)) {
			bits[errorPos] ^= 1
		}
	}

	// Извлекаем информационные биты (позиции 3,5,6,7 - индексы 2,4,5,6)
	infoBits := []byte{bits[2], bits[4], bits[5], bits[6]}

	return BitsToInt(infoBits)
}

// Генерация всех векторов ошибок
func getErrorsByClasses() [][]uint64 {
	errorClasses := make([][]uint64, n+1)

	for i := uint64(1); i < (1 << n); i++ {
		class := bits.OnesCount64(i)
		errorClasses[class] = append(errorClasses[class], i)
	}
	return errorClasses
}

func main() {
	fmt.Println("=== ВАРИАНТ 21 - КОД ХЭММИНГА [7,4] ===")
	fmt.Printf("Информационный вектор: 1011\n")
	fmt.Printf("Код: Хэмминга [7,4]\n")
	fmt.Printf("Способность: корректирующая (Ck)\n\n")

	// Кодируем информационный вектор
	encoded := encodeHamming(informationVector)
	fmt.Printf("Кодирование:\n")
	fmt.Printf("  Информационный вектор:  %04b\n", informationVector)
	fmt.Printf("  Закодированное слово:   %07b\n", encoded)

	// Пример с ошибкой
	errorVector := uint64(0b0000100) // Ошибка в 5-й позиции
	received := encoded ^ errorVector
	decoded := decodeHamming(received)

	fmt.Printf("\nПример исправления ошибки:\n")
	fmt.Printf("  Вектор ошибки:          %07b\n", errorVector)
	fmt.Printf("  Принятое слово:         %07b\n", received)
	fmt.Printf("  Декодированное слово:   %04b", decoded)
	if decoded == informationVector {
		fmt.Printf(" ✓ Ошибка исправлена\n")
	} else {
		fmt.Printf(" ✗ Ошибка не исправлена\n")
	}

	// Расчет корректирующей способности
	fmt.Printf("\nКорректирующая способность:\n")
	fmt.Printf("i  C(n,i)  Nk   Ck\n")
	fmt.Printf("───────────────────\n")

	errorClasses := getErrorsByClasses()

	for class := 1; class <= n; class++ {
		if len(errorClasses[class]) == 0 {
			continue
		}

		corrected := 0
		for _, errorVector := range errorClasses[class] {
			received := encoded ^ errorVector
			decoded := decodeHamming(received)
			if decoded == informationVector {
				corrected++
			}
		}

		total := len(errorClasses[class])
		Ck := float64(corrected) / float64(total)

		fmt.Printf("%d  %6d  %3d  %.3f", class, total, corrected, Ck)

		if Ck == 1.0 {
			fmt.Printf(" ✓\n")
		} else if Ck == 0.0 {
			fmt.Printf(" ✗\n")
		} else {
			fmt.Printf(" ~\n")
		}
	}
}