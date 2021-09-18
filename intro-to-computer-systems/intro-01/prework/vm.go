package vm

import "fmt"

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {
	fmt.Println("memory", memory)
	fmt.Println("len memory", len(memory))
	// Program command (Load, Store, Add, etc), first arg, second arg
	registers := [3]byte{8, 0, 0} // PC, R1 and R2
	pc := 8
	r1 := 9
	r2 := 10
	// Keep looping, like a physical computer's clock
	for {
		fmt.Println("reg before", registers)
		registers[0] = memory[pc]
		registers[1] = memory[r1]
		registers[2] = memory[r2]
		fmt.Println("reg after", registers)

		op := registers[0] // fetch the opcode

		// decode and execute
		switch op {
		// load    r1  addr    # Load value at given address into given register
		case Load:
			memory[registers[1]] = registers[2]
			// 	store   r2  addr    # Store the value in register at the given memory address
		case Store:
			fmt.Println("in store", registers)
			fmt.Println("mem at store", memory[0:21])
			memory[registers[2]] = memory[registers[1]]
			fmt.Println("Store", memory[registers[2]])
		// add     r1  r2      # Set r1 = r1 + r2
		case Add:
			fmt.Println("Add")
			fmt.Println("Going to sum", memory[r1], memory[r2])
			fmt.Println("What about", registers[1], registers[2])
			registers[1] += registers[2]
			fmt.Println("r1 is", registers[1])
		// sub     r1  r2      # Set r1 = r1 - r2
		case Sub:
			fmt.Println("Sub")
		// no params
		case Halt:
			fmt.Println("Halt")
			// break doesn't stop the loop
			// return doesn't stop the loop
			return
		}

		pc3 := pc + 3
		r13 := r1 + 3
		r23 := r2 + 3
		memory[pc] = memory[pc3]
		memory[r1] = memory[r13]
		memory[r2] = memory[r23]
		pc = pc3
		r1 = r13
		r2 = r23
	}
}
