package opcode

const (
	DIV_AL_VAL = 0x50
	DIV_AL_M   = 0x51
	DIV_AL_AH  = 0x52
	DIV_AL_B   = 0x53
	DIV_AL_C   = 0x54
	DIV_AL_D   = 0x55
	DIV_AL_E   = 0x56
	DIV_AL_L   = 0x57
	DIV_AL_H   = 0x58
)

var DIV = map[string]byte{
	"VAL": DIV_AL_VAL,
	"M":   DIV_AL_M,
	"AH":  DIV_AL_AH,
	"B":   DIV_AL_B,
	"C":   DIV_AL_C,
	"D":   DIV_AL_D,
	"E":   DIV_AL_E,
	"L":   DIV_AL_L,
	"H":   DIV_AL_H,
}
