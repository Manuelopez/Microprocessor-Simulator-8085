package opcode

const (
	XOR_AL_H   = 0x06
	XOR_AL_M   = 0x07
	XOR_AL_VAL = 0x08
	XOR_AL_AH  = 0x2F
	XOR_AL_B   = 0x30
	XOR_AL_C   = 0x31
	XOE_AL_D   = 0x32
	XOR_AL_E   = 0x33
	XOR_AL_L   = 0x34
)

var XOR = map[string]byte{
	"H":   0x06,
	"M":   0x07,
	"VAL": 0x08,
	"AH":  0x2F,
	"B":   0x30,
	"C":   0x31,
	"D":   0x32,
	"E":   0x33,
	"L":   0x34,
}
