package stdout

import "fmt"

func PrintFunc(a string) {
	str := fmt.Sprintf("-------------  %s   -----------", a)
	fmt.Println(str)
}
