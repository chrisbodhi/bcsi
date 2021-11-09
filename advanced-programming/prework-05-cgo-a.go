package main

// #include "square.c"
import "C"

import "fmt"

func cgoA() {
	n := C.int(5)
	sq := C.square(n)
	fmt.Println(fmt.Sprintf("%d squared is %d, which is type %T",
		n, sq, sq))
}
