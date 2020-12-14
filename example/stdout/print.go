package stdout

import (
	"fmt"
	"time"
)

func PrintFunc(a string) {
	str := fmt.Sprintf("-------------  %s   -----------", a)
	fmt.Println(str)
	time.Sleep(100 * time.Millisecond)
}
