package register

type FFRS struct {
	Value byte
	Load  byte
}

func (f *FFRS) SetLoad() {
	f.Load = 1
}

func (f *FFRS) UnsetLoad() {
	f.Load = 0
}

func (f *FFRS) LoadValue(value byte) {
	if f.Load != 1 {
		return
	}
	f.Value = value
}

func (f FFRS) GetValue() byte {
	return f.Value
}
