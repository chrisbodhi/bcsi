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
	// Program command (Load, Store, Add, etc), first arg, second arg
	pc := byte(8)
	registers := [3]byte{pc, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		registers[0] = memory[pc]
		registers[1] = memory[pc+1]
		registers[2] = memory[pc+2]

		op := registers[0] // fetch the opcode
		r1 := registers[1]
		r2 := registers[2]

		// decode and execute
		switch op {
		// load    r1  addr    # Load value at given address into given register
		case Load:
			memory[r1] = memory[r2]
		// 	store   r2  addr    # Store the value in register at the given memory address
		case Store:
			memory[r2] = memory[r1]
		// add     r1  r2      # Set r1 = r1 + r2
		case Add:
			memory[r1] += memory[r2]
		// sub     r1  r2      # Set r1 = r1 - r2
		case Sub:
			memory[r1] -= memory[r2]
		// jump    r1          # Set pc to pc + r1
		case Jump:
			pc = r1
			return
		// beqz    r1  r2      # If memory at r1 equals 0, increment the pc by the amount at r2; otherwise, continue
		case Beqz:
			if memory[r1] == 0 {
				pc += memory[r2]
				return
			}
		case Addi:
			fmt.Println("Addi")
		// no params
		case Halt:
			return
		}
		// Incrememnt program counter
		pc += byte(len(registers))
	}
}
