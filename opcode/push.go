package opcode

const (
	PUSH_VAL = 0xDE
	PUSH_AL  = 0xDF
	PUSH_AH  = 0xE0
	PUSH_B   = 0xE1
	PUSH_C   = 0xE2
	PUSH_D   = 0xE3
	PUSH_E   = 0xE4
	PUSH_L   = 0xE5
	PUSH_H   = 0xE6
)

var PUSH = map[string]byte{
	"VAL": PUSH_VAL,
	"AL":  PUSH_AL,
	"AH":  PUSH_AH,
	"B":   PUSH_B,
	"C":   PUSH_C,
	"D":   PUSH_D,
	"E":   PUSH_E,
	"L":   PUSH_L,
	"H":   PUSH_H,
}
