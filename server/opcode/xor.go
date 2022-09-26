package opcode

const (
	XOR_AL_H   = 0x06
	XOR_AL_M   = 0x07
	XOR_AL_VAL = 0x08
	XOR_AL_AH  = 0x2F
	XOR_AL_B   = 0x30
	XOR_AL_C   = 0x31
	XOR_AL_D   = 0x32
	XOR_AL_E   = 0x33
	XOR_AL_L   = 0x34
	XOR_S      = 0xF6
)

var XOR = map[string]byte{
	"H":   XOR_AL_AH,
	"M":   XOR_AL_M,
	"VAL": XOR_AL_VAL,
	"AH":  XOR_AL_AH,
	"B":   XOR_AL_B,
	"C":   XOR_AL_C,
	"D":   XOR_AL_D,
	"E":   XOR_AL_E,
	"L":   XOR_AL_L,
	"S":   XOR_S,
}
