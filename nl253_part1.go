package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"math/rand"
	"time"
)

func logLec(msg string) {
	fmt.Printf("[lecturer] %s\n", msg)
}

func logStudent(name string, msg string) {
	fmt.Printf("[student %s] %s\n", name, msg)
}

func sleepBetween(min int, max int) {
	time.Sleep((time.Duration(min) + time.Duration(rand.Intn(min+max))) * time.Second)
}

func lecturer(wait chan chan string, lect chan bool) {
	// prevent endless waiting once all students have been seen
	waitTime := 0
	maxWaitTime := 100
	for {
		if waitTime == maxWaitTime {
			logLec(fmt.Sprintf("waited for [%d/%d] time segments, done", waitTime, maxWaitTime))
			return
		}
		select {
		case dropInSess := <-wait:
			logLec("a student in the queue, inviting in")
			student := <-dropInSess
			logLec(fmt.Sprintf("student name is %s", student))
			logLec(fmt.Sprintf("drop in session started with %s...", student))
			sleepBetween(2, 5)
			logLec(fmt.Sprintf("drop in session finished, student %s has left", student))
		default:
			logLec("no students, going back to reading papers ...")
			sleepBetween(1, 3)
			lect <- true // signal free
		}
	}
}

func student(wait chan chan string, lect chan bool, name string) {
	logStudent(name, "arrived")
	<-lect // check if the lecturer is free
	logStudent(name, "checked and the lecturer is free!")
	dropIn := make(chan string, 1)
	dropIn <- name
	wait <- dropIn
	logStudent(name, "invited in")
}

func main() {
	queueSize := 3
	noStudents := 10
	fmt.Printf("creating %d students, queue size is %d\n", noStudents, queueSize)
	lect := make(chan bool)
	wait := make(chan chan string, queueSize)
	for i := 0; i < noStudents; i++ {
		go student(wait, lect, randomdata.SillyName())
	}
	// when running lecturer on the main thread you don't need sleep
	lecturer(wait, lect)
}
