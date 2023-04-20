package main

import (
	"fmt"
	"sort"
)

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	// is i<j return true but...
	// j=0 is always bigger than other number...
	// to remain the first
	if j == 0 {
		return false
	} else {
		return r.Interface.Less(j, i)
	}

}

func Custom(data sort.Interface) sort.Interface {
	return &reverse{data}
}

func main() {
	lst := []int{0, 4, 5, 2, 8, 1, 9, 3}
	sort.Sort(sort.IntSlice(lst))
	fmt.Println(lst)

	sort.Sort(Custom(sort.IntSlice(lst)))
	fmt.Println(lst)

}
