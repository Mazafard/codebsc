package main

import (
	"fmt"
	"sync"
	"time"
)

type Producer struct {
	// Mutex and ItemsProduced are removed, as Semaphore directly controls the number
	// Of active goroutines, eliminating the need for manual locking.

	// Mutex sync.Mutex
	// ItemsProduced int
	ItemsToProduce int
	output         chan int
	sema           chan struct{}
}

// @ todo: remove everything not useful, necessary lines codes ( channel)
func (p *Producer) produce(wg *sync.WaitGroup) {
	defer close(p.output)
	defer wg.Done()

	for i := 0; i < p.ItemsToProduce; i++ {
		p.sema <- struct{}{}
		fmt.Printf("Produced item: %d\n", i)
		p.output <- i
		<-p.sema
	}
}

type Consumer struct {
	// The quit channel is removed and a sema channel is added to allow the consumer to use Semaphore
	quit  chan int
	input chan int
	sema  chan struct{}
}

func (consumer *Consumer) consume(wg *sync.WaitGroup) {
	// The code for handling the quit channel and closing the input channel is removed.
	// The sync.WaitGroup is used to synchronize the completion of the consumer.
	defer wg.Done()
	//defer close(consumer.input)

	//for {
	//	select {
	//	case input := <-consumer.input:
	//		fmt.Printf("Consumed item: %d\n", input)
	//		time.Sleep(5 * time.Second) // Simulating work
	//	case <-consumer.quit:
	//		consumer.quit <- 2
	//		return
	//	}
	//}
	for item := range consumer.input {
		// the consumer also uses the sema channel to control the number of active goroutines
		consumer.sema <- struct{}{}
		fmt.Printf("Consumed item: %d\n", item)
		time.Sleep(2 * time.Second) // Simulating work
		<-consumer.sema
	}
}

func main() {
	var wg sync.WaitGroup

	sema := make(chan struct{}, 2) // Semaphore with capacity 2

	producer := Producer{output: make(chan int), sema: sema}
	producer.ItemsToProduce = 10
	consumer := Consumer{input: producer.output, sema: sema}

	wg.Add(1)
	go consumer.consume(&wg)

	wg.Add(1)
	go producer.produce(&wg)

	fmt.Println("Start")
	wg.Wait()

	fmt.Println("Finished")
}
