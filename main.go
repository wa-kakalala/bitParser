package main

import (
	"bitParser/utils"
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("--------  start ---------")
	fp, err := os.Open("./data/info.dat")
	if err != nil {
		fmt.Println("open file error", err)
	}
	defer fp.Close()

	bufr := bufio.NewReader(fp)
	bs := utils.CreateBitStructure(1, 16)

	bitData, _ := bs.BitFifo2NBitsNum(6, 1, bufr, true)
	fmt.Println("result", bitData)
	bitData, _ = bs.BitFifo2NBitsNum(5, 1, bufr, true)
	fmt.Println("result", bitData)
}
