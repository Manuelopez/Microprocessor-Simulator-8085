package main

import (
	"fmt"
	"micp-sim/microprocessor"
	"micp-sim/util"
)

func main() {
  start()
// testUtill()
}

func start() {
  m := microprocessor.New(5)
	m.Start()
}

func testUtill(){
  //a := util.DecimalToBinary(5)
  b, c := util.DecimalToBinary16(102)
  //fmt.Println(util.BinaryToDecimal(a[:]))
  d := []byte{}
  d = append(d, b[:]...)
  d= append(d, c[:]...)
  fmt.Println(util.BinaryToDecimal(d[:]))


}

 
