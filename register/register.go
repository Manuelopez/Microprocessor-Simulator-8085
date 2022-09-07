package register


type Register struct{
  Load byte
  registers [8]FFRS
}

func New() Register{
  return Register{}
}

func (r *Register) SetLoad(){
  r.Load = 1
}

func (r *Register) UnsetLoad(){
  r.Load = 0
}

func (r *Register) LoadValue(value [8]byte){
  if(r.Load != 1){
    return
  }
  
  for i, bit := range value{
    r.registers[i].SetLoad()
    r.registers[i].LoadValue(bit)
  }
  r.UnsetLoad()
}

func (r Register) GetValue() [8]byte{
  a := [8]byte{}
  for i, v := range r.registers{
    a[i] = v.Value
  }

  return a 
} 


//--------------------------------------------



