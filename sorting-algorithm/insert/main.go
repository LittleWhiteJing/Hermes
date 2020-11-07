package main

import (
	"fmt"
)

//插入排序
func InsertSort (a []int) ([]int, int) {
	count := 0
	length := len(a)
	for i := 1; i < length; i++ {
		tmp := a[i]
		j := i
		for ; (j > 0) && (tmp < a[j-1]); j-- {
			count++
			a[j] = a[j-1]
		}
		a[j] = tmp
	}
	return a, count
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	afterSortSet, count := InsertSort(beforeSortSet)
	fmt.Println("After Sort:", afterSortSet, "Swap Count:", count)
}