package main

import (
	"fmt"
	"sync"
)

type gimmieStuff interface {
	GetPrinter() func([]string)
}

func main() {

	s := []string{"hello", "dolly", "world"}

	s2 := s[1:]
	s3 := s[0:1]
	s4 := make([]string, len(s[1:]))
	copy(s4, s[1:])

	fmt.Println(s)
	fmt.Println(s2)
	fmt.Println(s3)

	s[0] = "goodbye"
	s[1] = "cruel"

	fmt.Println(s)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)

	f := GetPrinter()
	f(s)
	f(s2)
	f(s3)

	g := giver{}

	f2 := g.GetPrinter()
	f2(s)

	f3 := GetAPrinter(g)
	f3(s)

	m := make(map[int]string)

	m[1] = "one"
	m[2] = "two"
	m[3] = "three"

	for k, v := range m {
		fmt.Println(k, ":", v)
	}

	gg := newGoGiver()
	gg.Print(s)
	gg.Done()
}

func GetPrinter() func([]string) {
	return func(s []string) {
		fmt.Println(s)
	}
}

type giver struct{}

func (g giver) GetPrinter() func([]string) {
	return func(s []string) {
		fmt.Println("giver:", s)
	}
}

type printer func([]string)

func GetAPrinter(g gimmieStuff) printer {
	return g.GetPrinter()
}

type goGiver struct {
	s    chan []string
	done chan bool
	wg   *sync.WaitGroup
}

func newGoGiver() (g goGiver) {
	g = goGiver{s: make(chan []string),
		done: make(chan bool),
		wg:   &sync.WaitGroup{}}
	g.wg.Add(1)
	go g.Run()
	return g
}

func (g goGiver) Run() {
	defer g.wg.Done()
	for {
		select {
		case msg := <-g.s:
			fmt.Println("goGiver:", msg)
		case <-g.done:
			fmt.Println("closing")
			return

		}
	}
}

func (g goGiver) Done() {
	g.done <- true
	g.wg.Wait()
	fmt.Println("closed")
}
func (g goGiver) Print(s []string) {
	g.s <- s
}
