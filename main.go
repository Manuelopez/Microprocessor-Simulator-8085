package main

import (
	"fmt"
	"micp-sim/microprocessor"
	"micp-sim/util"
)

func main() {
  start()
  //testUtill()
}

func start() {
  m := microprocessor.New(5)
	m.Start()
}

func testUtill(){
  //a := util.DecimalToBinary(5)
  b := util.DecimalToBinary(500)
  //fmt.Println(util.BinaryToDecimal(a[:]))
  fmt.Println(util.BinaryToDecimal(b[:]))


}


