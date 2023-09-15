package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func main() {
	command := exec.Command("go", "run", "code/user-code/main.go")
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	stinPipe, err := command.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stinPipe, "13 2\n")
	err = command.Run()
	if err != nil {
		log.Fatalln(err, stderr.String())
	}
	log.Println(out.String())
}
