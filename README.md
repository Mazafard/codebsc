## Problem:
Code that utilizes many goroutines can run the risk of bombarding some remote resource with hundreds or even thousands of goroutines. 

A bombarded system can quickly become overloaded to the point where you experience extreme latency spikes or worse service outages.

## Revised code

In the revised version, I replaced the mutexes with semaphores to handle concurrency control more effectively. 

Semaphores help in limiting the number of concurrent goroutines, which can prevent resource contention and improve the overall efficiency of the program. 

By using semaphores, we also eliminate the need for explicit locking and unlocking, simplifying the code.

## Specific Changes:
1. Removed Mutexes: The explicit use of sync.Mutex was removed from the producer struct. This reduces the complexity of managing locks manually.

```go
type Producer struct {
    ItemsToProduce int
    output         chan int
    sema           chan struct{}
}
```
2. Added Semaphores: Introduced a sema channel to act as a semaphore, controlling the number of goroutines that can execute concurrently.
```go
type Producer struct {
    ItemsToProduce int
    output         chan int
    sema           chan struct{}
}
```
```go
producer := Producer{output: make(chan int), sema: sema}
```
```go
consumer := Consumer{input: producer.output, sema: sema}
```

3. Simplified Synchronization: Used sync.WaitGroup to handle the synchronization of goroutines, ensuring that the main function waits for the completion of all goroutines before exiting.
```go
func (p *Producer) produce(wg *sync.WaitGroup) {
    defer wg.Done()
    ...
```
```go
func (consumer *Consumer) consume(wg *sync.WaitGroup) {
    defer wg.Done()
    ...
```
4. Removed Manual State Management: By using semaphores and channels, we avoid manual state management, making the code cleaner and easier to maintain.

## Further Improvement

For more granular control of concurrency, a [weighted semaphore](https://github.com/golang/sync/tree/master/semaphore) can be used, which allows managing multiple units of a resource simultaneously. 

This can improve the performance and resource utilization of the application even further.

