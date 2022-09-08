package opcode

const (
	ADD_VAL = 0x35
	ADD_AH  = 0x37
	ADD_B   = 0x38
	ADD_C   = 0x39
	ADD_D   = 0x3A
	ADD_E   = 0x3B
	ADD_L   = 0x3C
	ADD_H   = 0x3D
)

// MAP TO CONVERT ADD TO Number
var ADD = map[string]byte{
	"VAL": ADD_VAL,
	"AH":  ADD_AH,
	"B":   ADD_B,
	"C":   ADD_C,
	"D":   ADD_D,
	"E":   ADD_E,
	"L":   ADD_L,
	"H":   ADD_H,
}
