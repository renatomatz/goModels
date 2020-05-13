package models

import (
    "errors"
    "fmt"
)

type Queue interface {
    Push(int)
    Pop() int
    IsEmpty() bool
}

type LIFO struct {
    items []int
}

type FIFO struct {
    items []int
}

//Push(int) 
func (queue *LIFO) Push(item int) {
    queue.items = append(queue.items, item)
}

func (queue *FIFO) Push(item int) {
    queue.items = append(queue.items, item)
}

//Pop() int
func (queue *LIFO) Pop() int {
    var ret int
    ret, queue.items = queue.items[len(queue.items)-1], queue.items[:len(queue.items)-1]
    return ret
}

func (queue *FIFO) Pop() int {
    var ret int
    ret, queue.items = queue.items[0], queue.items[1:]
    return ret
}

//IsEmpty() bool
func (queue *LIFO) IsEmpty() bool {
    return len(queue.items) == 0
}

func (queue *FIFO) IsEmpty() bool {
    return len(queue.items) == 0
}

//Print()
func (queue *LIFO) Print() {
    fmt.printf("%v (OUT)", queue.items)
}

func (queue *FIFO) Print() {
    fmt.printf("(OUT) %v", queue,items)
}
