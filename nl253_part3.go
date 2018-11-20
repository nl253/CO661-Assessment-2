package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"math/rand"
	"time"
)

func logLec(name string, msg string) {
	fmt.Printf("[lecturer %s] %s\n", name, msg)
}

func logStudent(name string, msg string) {
	fmt.Printf("[student %s] %s\n", name, msg)
}

func sleepBetween(min int, max int) {
	time.Sleep((time.Duration(min) + time.Duration(rand.Intn(min+max))) * time.Second)
}

func student(wait chan chan string, name string) {
	logStudent(name, fmt.Sprintf("wants to meet"))
	var dropInSess = make(chan string, 1)
	dropInSess <- name
	select {
	case wait <- dropInSess:
		logStudent(name, "found a place in the waiting room")
	default:
		logStudent(name, "waiting from is full, bye!")
	}
}

func lecturer(wait chan chan string, name string, done chan bool) {

	// prevent endless waiting once all students have been seen
	waitTime := 0
	maxWaitTime := 100

	for {
		if waitTime == maxWaitTime {
			logLec(name, fmt.Sprintf("waited for [%d/%d] time segments, done", waitTime, maxWaitTime))
			done <- true // signal to main()
			return
		}
		select {
		case dropInSess := <-wait:
			waitTime = 0 // reset counter
			logLec(name, "a student in the queue, inviting in")
			student := <-dropInSess
			logLec(name, fmt.Sprintf("student name is %s", student))
			logLec(name, fmt.Sprintf("drop in session started with %s...", student))
			sleepBetween(2, 5)
			logLec(name, fmt.Sprintf("drop in session finished, student %s has left", student))
		default:
			logLec(name, fmt.Sprintf("no students, waiting room empty, going back to reading papers [waited %d/%d] ...", waitTime, maxWaitTime))
			waitTime++
			sleepBetween(0, 3)
		}
	}
}

func main() {
	noLects := 10
	noStuds := 1000
	fmt.Printf("creating %d students & %d lecturers\n", noStuds, noLects)
	wait := make(chan chan string)
	done := make(chan bool, noLects)
	for i := 0; i < noLects; i++ {
		go lecturer(wait, randomdata.SillyName(), done)
	}
	for i := 0; i < noStuds; i++ {
		sleepBetween(0, 1)
		go student(wait, randomdata.SillyName())
	}
	// wait for lecturers to finish
	for i := 0; i < noLects; i++ {
		<-done
	}
}
