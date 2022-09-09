package opcode

const (
  NOT_AL = 0x15
  NOT_AH = 0x16
  NOT_B = 0x17
  NOT_C = 0x18
  NOT_D = 0x19
  NOT_E = 0x1A
)

var NOT = map[string]byte{
  "AL": NOT_AL,
  "AH": NOT_AH,
  "B": NOT_B,
  "C": NOT_C,
  "D": NOT_D,
  "E": NOT_E,
}
