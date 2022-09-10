package opcode

const (
	AND_AL_AH  = 0x24
	AND_AL_B   = 0x25
	AND_AL_C   = 0x26
	AND_AL_D   = 0x27
	AND_AL_E   = 0x28
	AND_AL_L   = 0x29
	AND_AL_H   = 0x2A
	AND_AL_M   = 0x2C
	AND_AL_VAL = 0x2E
	AND_S      = 0xF5
)

var AND = map[string]byte{
	"AL":  AND_AL_AH,
	"B":   AND_AL_B,
	"C":   AND_AL_C,
	"D":   AND_AL_D,
	"E":   AND_AL_E,
	"L":   AND_AL_L,
	"H":   AND_AL_H,
	"M":   AND_AL_M,
	"VAL": AND_AL_VAL,
	"S":   AND_S,
}
