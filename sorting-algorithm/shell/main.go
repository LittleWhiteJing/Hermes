package main

import (
	"fmt"
)

func ShellSort(a []int) ([] int, int){
	count := 0
	len := len(a)

	for i := len/2; i > 0; i = i/2 {
		for j := i; j < len; j++ {
			tmp := a[j]
			k := j
			for ; k >= i && tmp < a[k-1]; k = k - i {
				a[k] = a[k-1]
				count++
			}
			a[k] = tmp
		}
	}
	return a, count
}

func main() {
	beforeSortSet := []int {10, 22, 33, 21, 56, 32, 81, 73, 69, 83}
	fmt.Println("Before Sort:", beforeSortSet)
	afterSortSet, count := ShellSort(beforeSortSet)
	fmt.Println("After Sort:", afterSortSet, "Swap Count:", count)
}
