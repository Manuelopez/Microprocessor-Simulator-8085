package opcode

const (
	INC_AL = 0x64
	INC_AH = 0x65
	INC_B  = 0x66
	INC_C  = 0x67
	INC_D  = 0x68
	INC_E  = 0x69
	INC_H  = 0x6A
	INC_L  = 0x6B
)

var INC = map[string]byte{
	"AL": INC_AL,
	"AH": INC_AH,
	"B":  INC_B,
	"C":  INC_C,
	"D":  INC_D,
	"E":  INC_E,
	"H":  INC_H,
	"L":  INC_L,
}
