package main

import (
	"fmt"
)

func QuickSort(a []int) {
	sort(a, 0, len(a)-1)
}

func sort(a []int, left int, right int) {
	if right - left < 2 {
		return
	}
	p := middle(a, left, right)
	i := left + 1
	j := right - 1
	for {
		for a[i] < p {
			i++
		}
		for a[j] > p {
			j++
		}
		if i < j {
			a[i], a[j] = a[j], a[i]
		} else {
			break
		}
	}
	a[i], a[j-1] = a[j-1], a[i]
	sort(a, left, i - 1)
	sort(a, i + 1, right)
}

func middle(a []int, left int, right int) int {
	center := (right + left)/2
	if a[left] > a[center] {
		a[left], a[center] = a[center], a[left]
	}
	if a[left] > a[right] {
		a[left], a[right] = a[right], a[left]
	}
	if a[center] > a[right] {
		a[center], a[right] = a[right], a[center]
	}
	a[right-1], a[center] = a[center], a[right-1]
	return a[right-1]
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	QuickSort(beforeSortSet)
	fmt.Println("After Sort:", beforeSortSet)
}

