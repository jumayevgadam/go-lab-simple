package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Visual learing  is the best way of learning
// About concurrency

func main() {
	ch := make(chan order, 3)

	wg := &sync.WaitGroup{} // More on WaitGroup another day
	wg.Add(2)

	for i := 0; i < 10; i++ {
		waitForOrders()
		o := order(i)
		log.Printf("Partier: I %v, I will pass it to the channel\n", o)
		ch <- o
	}

	go func() {
		defer wg.Done()
		worker("Candier", ch)
	}()

	go func() {
		defer wg.Done()
		worker("Stringer", ch)
	}()

	log.Println("No More Orders, closing the channel to signify workers to stop")
	close(ch)

	log.Println("Wait for workers to gracefully stop")
	wg.Wait()

	log.Println("All done")
}

func waitForOrders() {
	processingTime := time.Duration(rand.Intn(2)) * time.Second
	time.Sleep(processingTime)
}

func worker(name string, ch <-chan order) {
	for o := range ch {
		log.Printf("%s: I got %v, I will process it\n", name, o)
		processOrder(o)
		log.Printf("%s: I completed %v, I'm ready to take a new order\n", name, o)
	}

	log.Printf("%s: I'm done \n", name)
}

func processOrder(_ order) {
	processingTime := time.Duration(2+rand.Intn(2)) * time.Second
	time.Sleep(processingTime)
}

type order int

func (o order) String() string {
	return fmt.Sprintf("order-%02d", o)
}
