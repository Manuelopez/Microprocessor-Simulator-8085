package opcode

const (
	DEC_AL = 0x6C
	DEC_AH = 0x6D
	DEC_B  = 0x6E
	DEC_C  = 0x6F
	DEC_D  = 0x70
	DEC_E  = 0x71
	DEC_H  = 0x72
	DEC_L  = 0x73
)

var DEC = map[string]byte{
	"AL": DEC_AL,
	"AH": DEC_AH,
	"B":  DEC_B,
	"C":  DEC_C,
	"D":  DEC_D,
	"E":  DEC_E,
	"H":  DEC_H,
	"L":  DEC_L,
}
