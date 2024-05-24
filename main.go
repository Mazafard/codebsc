package main

import (
	"fmt"
	"sync"
)

type Producer struct {
	Mutex          sync.Mutex
	ItemsToProduce int
	ItemsProduced  int
	output         chan int
}

// @ todo: remove everything not useful, necessary lines codes ( channel)
func (p *Producer) produce() {
	for i := 0; i < p.ItemsToProduce; i++ {
		fmt.Printf("Produced item: %d\n", i)
		p.Mutex.Lock()
		p.ItemsProduced++
		p.Mutex.Unlock()
		p.output <- i
	}
}

type Consumer struct {
	input chan int
	quit  chan int
}

func (consumer *Consumer) consume() {
	defer close(consumer.input)
	for {
		select {
		case input := <-consumer.input:
			fmt.Printf("Consumed item: %d\n", input)
		case <-consumer.quit:
			consumer.quit <- 2
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	producer := Producer{output: make(chan int)}
	producer.ItemsToProduce = 3
	consumer := Consumer{input: producer.output, quit: make(chan int)}
	go func() {
		consumer.consume()
	}()
	go func() {
		for isProducing := true; isProducing; {
			producer.Mutex.Lock()
			isProducing = producer.ItemsProduced != producer.ItemsToProduce
			producer.Mutex.Unlock()
			if !isProducing {
				consumer.quit <- 1
				<-consumer.quit
				fmt.Println("Finished")
				wg.Done()
			}
		}

	}()
	producer.produce()
	wg.Wait()
}
