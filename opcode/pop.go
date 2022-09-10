package opcode

const (
	POP_AL = 0xE7
	POP_AH = 0xE8
	POP_B  = 0xE9
	POP_C  = 0xEA
	POP_D  = 0xEB
	POP_E  = 0xEC
	POP_L  = 0xED
	POP_H  = 0xEE
)

var POP = map[string]byte{
	"AL": POP_AL,
	"AH": POP_AH,
	"B":  POP_B,
	"C":  POP_C,
	"D":  POP_D,
	"E":  POP_E,
	"L":  POP_L,
	"H":  POP_H,
}
