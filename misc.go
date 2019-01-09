package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func outputStatm() {
	buf, err := ioutil.ReadFile("/proc/self/statm")
	if err != nil {
		return
	}
	split := bytes.Split(buf, []byte{0x20})
	fmt.Println("Statm: size", string(split[0]), "resident", string(split[1]), "share", string(split[2]))
}
