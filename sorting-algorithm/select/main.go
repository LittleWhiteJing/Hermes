package main

import (
	"fmt"
)

//选择排序
func SelectSort (a []int) []int {
	length := len(a)
	minIndex := 0
	for i := 0; i < length; i++ {
		minIndex = i
		for j := i; j < length; j++ {
			if a[minIndex] > a[j] {
				minIndex = j
			}
		}
		if minIndex != i {
			a[minIndex], a[i] = a[i], a[minIndex]
		}
	}
	return a
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	afterSortSet := SelectSort(beforeSortSet)
	fmt.Println("After Sort:", afterSortSet)
}
