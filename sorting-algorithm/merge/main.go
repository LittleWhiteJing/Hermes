package main

import (
	"fmt"
)

const INFINITY = 0xffff

func MergeSort(a []int) {
	merge(a, 0, len(a) - 1)
}

func merge(a []int, start int, end int) {
	if start < end {
		mid := (start + end)/2
		merge(a, start, mid)
		merge(a, mid + 1, end)

		arr1 := make([]int, mid - start + 2)
		copy(arr1, a[start:mid+1])
		arr1[mid - start + 1] = INFINITY

		arr2 := make([]int, end - mid + 1)
		copy(arr2, a[mid+1:end+1])
		arr2[end - mid] = INFINITY

		j, k := 0, 0
		for i := start; i <= end; i++ {
			if arr1[j] <= arr2[k] {
				a[i] = arr1[j]
				j++
			} else {
				a[i] = arr2[k]
				k++
			}
		}
	}
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	MergeSort(beforeSortSet)
	fmt.Println("After Sort:", beforeSortSet)
}
