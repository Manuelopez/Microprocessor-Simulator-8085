package main

import (
	"fmt"
	"micp-sim/microprocessor"
	"micp-sim/util"
)

func main() {
  start()
}

func start() {
  m := microprocessor.New(2)
	m.Test()
}

func testUtill(){
  //a := util.DecimalToBinary(5)
  b := util.DecimalToBinary(10)
  //fmt.Println(util.BinaryToDecimal(a[:]))
  fmt.Println(util.BinaryToDecimal(b[:]))


}


