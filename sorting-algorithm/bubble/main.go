package main

import (
	"fmt"
)

//冒泡排序(稳定排序)
func BubbleSort(a []int) ([]int, int) {
	 count := 0
	 length := len(a)
	 for i := length - 1; i >= 0 ; i-- {
	 	flag := 0
	 	for j := 0; j < i; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
				count++
				flag = 1
			}
		}
		if flag == 0 {
			break
		}
	 }
	 return a, count
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	afterSortSet, count := BubbleSort(beforeSortSet)
	fmt.Println("After Sort:", afterSortSet, "Swap Count:", count)
}