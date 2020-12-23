package main

import (
	"fmt"
	"github.com/ivalue2333/pkg/src/stringx"
)

func main() {

	fmt.Println("Rand:", stringx.Rand())
	fmt.Println("RandId:", stringx.RandId())
	fmt.Println("Randn:", stringx.Randn(8))
}
