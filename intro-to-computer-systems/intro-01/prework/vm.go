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
	registers := [3]byte{8, 0, 0} // PC, R1 and R2
	for i, v := range memory[0:30] {
		if i == 0 {
			fmt.Println("===Data===")
		}
		if i == 8 {
			fmt.Println("===Instruction===")
		}
		fmt.Println(i, ": ", v)
	}
	// Keep looping, like a physical computer's clock
	for {
		op := memory[registers[0]] // fetch the opcode
		mem1 := memory[registers[0]+1]
		mem2 := memory[registers[0]+2]

		// decode and execute
		switch op {
		// load		r1 addr		# Load value at given address into given register
		case Load:
			fmt.Println("Load", registers)
			registers[mem1] = memory[mem2]
			// fmt.Println("After load", registers)
		// store	r2 addr		# Store the value in register at the given memory address
		case Store:
			fmt.Println("Store", registers, memory[registers[0]:registers[0]+3])
			memory[mem2] = registers[mem1]
			// fmt.Println("After store", memory[0:8])
		// add		r1 r2		# Set r1 = r1 + r2
		case Add:
			fmt.Println("Add", registers)
			registers[mem1] += registers[mem2]
		// sub		r1 r2		# Set r1 = r1 - r2
		case Sub:
			fmt.Println("Sub", registers)
			registers[mem1] -= registers[mem2]
		// jump		r1			# Set pc to pc + r1
		case Jump:
			fmt.Println("Jump", registers)
			registers[0] = mem1
			continue
		// beqz		r1 r2		# If memory at r2 equals 0, increment the pc by the amount at r1; otherwise, continue
		case Beqz:
			fmt.Println("Beqz", registers, memory[registers[0]+2])
			if registers[1] == 0 {
				registers[0] += mem2
			}
			fmt.Println("End of Beqz", registers)
			continue
		// addi		r1 r2		# Add values passed as args, not adding what's in store
		case Addi:
			registers[mem1] += mem2
		// subi		r1 r2		# Subtract
		case Subi:
			fmt.Println("Subi", registers, memory[registers[0]:registers[0]+3])
			registers[mem1] -= mem2
		// no params
		case Halt:
			fmt.Println("~~~Halt~~~")
			return
		}
		// Incrememnt program counter
		registers[0] += byte(len(registers))
	}
}
