package main

import (
	"C"
	"fmt"
	"time"
)

//export say
func say(text *C.char) {
	fmt.Println(C.GoString(text))
}

type Employees struct {
	EmployeeID int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	HireDate   time.Time
	ManagerID  int
	JobTitle   string
}

func main() {}
