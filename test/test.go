package main

import (
	"fmt"
	"regexp"
)

func main() {
	r := regexp.MustCompile(`^([a-z]+)([0-9]{1})$`)
	fmt.Println(r.FindStringSubmatch("d5"))

}
