package models

type LinkedList interface {
	Traverse() <-chan int
	ToArray() []*int
	Get(int) *int
	Delete(int) *int
	Split(int, int) LinkedList
}

// Struct definitions

type SingleLinkedList struct {
	Value *int
	next *SingleLinkedList
}

type DoubleLinkedList struct {
	Value *int
	next *DoubleLinkedList
	prev *DoubleLinkedList
}

// Traverse() <-chan int
// ToArray() []*int
// Get(int) *int
// Delete(int) *int
// Split(int, int) LinkedList
