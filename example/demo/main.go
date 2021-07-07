package main

import (
	"github.com/cocotyty/oneshot"
	"log"
	"time"
)

func main() {
	s := &oneshot.Shot{}
	go func() {
		time.Sleep(time.Second)
		s.WakeUp()
	}()
	log.Println("wait")
	s.Wait()
	log.Println("finished")
}