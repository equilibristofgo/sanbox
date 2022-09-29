package main

import (
	"fmt"
	"sort"
)

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func Reverse(data sort.Interface) sort.Interface {
	return &reverse{data}
}

func main() {
	lst := []int{4, 5, 2, 8, 1, 9, 3}
	sort.Sort(sort.IntSlice(lst))
	fmt.Println(lst)

	sort.Sort(sort.Reverse(sort.IntSlice(lst)))
	fmt.Println(lst)
}
