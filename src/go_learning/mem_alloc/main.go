package main

import (
	"fmt"
	"unsafe"
)

type Employee struct {
	/*
		Id     int
		Name   string
		Age    int16
		Gender string
		Active bool

		For above order, the memory alignment will be as follows in a 64bit architecture:
		CPU cycle  	| Memory
		1st cycle	| Id (8 bytes)
		2nd cycle	| Name (8 bytes)
		3rd cycle	| Name (8 bytes)
		4th cycle	| Age (2 bytes)
		5th cycle	| Gender (8 bytes)
		6th cycle	| Gender (8 bytes)
		7th cycle	| Active (1 byte)
	*/
	Id     int
	Name   string
	Gender string
	Age    int16
	Active bool
	/*
		For above order, the memory alignment will be as follows in a 64bit architecture:
		CPU cycle  	| Memory
		1st cycle	| Id (8 bytes)
		2nd cycle	| Name (8 bytes)
		3rd cycle	| Name (8 bytes)
		4th cycle	| Gender (8 bytes)
		5th cycle	| Gender (8 bytes)
		6th cycle	| Age (2 bytes) + Active (1 byte)
	*/
}

/*
Type | Size
Int	| 8 bytes
string | 16 bytes
int16 | 2 bytes
bool | 1 byte
*/
func main() {
	var e Employee
	fmt.Printf("Size of %T struct: %d bytes", e, unsafe.Sizeof(e))
}
