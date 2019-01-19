package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	file, err := os.Open("/dev/tty")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	fd := int(file.Fd())
	buf := make([]byte, 1)
	data := make(chan byte)

	go func() {
		for {
			_, err := syscall.Read(fd, buf)
			if err != nil {
				close(data)
				return
			}
			data <- buf[0]
		}
	}()

	for {
		b := <-data
		fmt.Println(b)
	}
}
