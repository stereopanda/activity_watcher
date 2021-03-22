package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)
var logfile = "/var/log/activity_watcher.log"

func main() {
	logfile, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	stride := 256
	filename := []string{
		"/dev/input/event5",
		"/dev/input/mice",
	}
	input_device := ""

	sliceLength := len(filename)
	var wg sync.WaitGroup
	wg.Add(sliceLength)

	for i := 0; i < sliceLength; i++ {
		go func(i int) {
			defer wg.Done()
			val := filename[i]
			f, err := os.Open(val)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			defer f.Close()

			r := bufio.NewReader(f)
			buf := make([]byte, 0, stride)
			for {
				n, err := io.ReadFull(r, buf[:cap(buf)])
				buf = buf[:n]
				if err != nil {
					if err == io.EOF {
						break
					}
					if err != io.ErrUnexpectedEOF {
						fmt.Fprintln(os.Stderr, err)
						break
					}
				}

				if val == "/dev/input/event5" {
					input_device = "Keyboard"
				} else {input_device = "Mouse"}
				log.Println(input_device,"activity detected")

				// process buf
			}
		}(i)
	}
	wg.Wait()
}


