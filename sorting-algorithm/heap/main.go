package main

import (
	"fmt"
)

func HeapSort(a []int) {
	length := len(a)
	for i := length/2-1; i >= 0; i-- {
		HeapAjust(a, i, length)
	}
	for i := length - 1; i > 0; i-- {
		a[i], a[0] = a[0], a[i]
		HeapAjust(a, 0, i)
	}
}

func HeapAjust(a []int, start int, length int) {
	tmp := a[start]
	for i := 2*start + 1; i < length; i = i * 2 {
		if i+1 < length && a[i] < a[i+1] {
			i++
		}
		if tmp > a[i] {
			break
		}
		a[start] = a[i]
		start = i
	}
	a[start] = tmp
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	HeapSort(beforeSortSet)
	fmt.Println("After Sort:", beforeSortSet)
}
