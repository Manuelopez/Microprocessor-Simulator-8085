package opcode

const (
	CLR_CA = 0xF0
	CLR_AL = 0xF7
	CLR_AH = 0xF8
	CLR_B  = 0xF9
	CLR_C  = 0xFA
	CLR_D  = 0xFB
	CLR_E  = 0xFC
	CLR_H  = 0xFD
	CLR_L  = 0xFE
)

var CLR = map[string]byte{
	"CA": CLR_CA,
	"AL": CLR_AL,
	"AH": CLR_AH,
	"B":  CLR_B,
	"C":  CLR_C,
	"D":  CLR_D,
	"E":  CLR_E,
	"H":  CLR_H,
	"L":  CLR_L,
}
