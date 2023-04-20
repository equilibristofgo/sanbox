package main

import (
	"C"
	"fmt"
	"time"

	"gorm.io/gorm"
)
import "github.com/captain-corgi/golang-oracledb-example/pkg/oracle"

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

//export call_oracle_db
func call_oracle_db(text *C.char) {

	dsn := "oracle://OT:yourpassword@localhost/XE"
	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Read
	var employees []Employees
	db.Debug().
		Where("employee_id = ?", 100).
		Limit(5).
		Find(&employees)

	// Print
	for _, employee := range employees {
		fmt.Printf("%+v\n", employee)
	}

}

func main() {}
