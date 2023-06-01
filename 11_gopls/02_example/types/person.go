package types

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) String() string {
	return fmt.Sprintf("%v (%v)", p.Name, p.Age)
}

type PersonToBool func(*Person) bool
type PersonList []*Person

func (pl PersonList) Filter(f PersonToBool) PersonList {
	var ret []*Person
	for _, p := range pl {
		if f(p) {
			ret = append(ret, p)
		}
	}
	return ret
}
