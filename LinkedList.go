package models

import = (
    "errors"
    "fmt"
    "sync"
}

type LinkedList interface {
	Traverse() <-chan *int
	ToArray() []*int
	Get(int) (*int, error)
        Append(int) error
	Delete(int) (*int, error)
	Split(int, int) (LinkedList, error)
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

func (root *SingleLinkedList) Traverse() <-chan *int {

    ch = make(chan *int)

    go sLLTraverseHelper(root, ch)
    return ch

}

func (root *DoubleLinkedList) Traverse() <-chan *int {

    chLeft = make(chan *int)
    chRight = make(chan *int)

    var after int
    if root.prev != nil {
        after = root.Value
        go dLLLeftTraverseHelper(root.prev, root, chLeft)
    }

    go dLLRightTraverseHelper(root, chRight, chLeft)
    return chRight

}

// ToArray() []*int

func (root *SingleLinkedList) ToArray() {

    var ret []*int
    for val := range root.Traverse() {
        append(ret, val)
    }

    return ret

}

func (root *DoubleLinkedList) ToArray() {

    var ret []*int
    for val := range root.Traverse() {
        append(ret, val)
    }

    return ret

}

// Get(int) *int

func (root *SingleLinkedList) Get(i int) (*int, error) {

    node, err := root.getNode(i)
    return node.Value, err

}

func (root *DoubleLinkedList) Get(i int) (*int, error) {

    node, err := root.getNode(i)
    return node.Value, err

}

// Delete(int) *int

func (root *SingleLinkedList) Delete(int) (*int, error) {

    if node, err := root.getNode(i-1); err {
        return nil, err
    }

    if node.next == nil {
        return nil, errors.New("index out of range")
    } else {
        ret := node.next.Value
        node.next := node.next.next
        return ret, nil
    }

}

func (root *SingleLinkedList) Delete(int) (*int, error) {

    if node, err := root.getNode(i); err {
        return nil, err
    }

    ret := node.Value

    nextNode := node.next
    prevNode := node.prev

    nextNode.prev = prevNode
    prevNode.next = nextNode

    return ret, nil

}

// Split(int, int) LinkedList
// These allow for looped lists

func (root *SingleLinkedList) Split(i, j int) (LinkedList, error) {

    if j < i {
        return nil, errors.New("right index must be smaller than left")
    }

    if left, err := root.getNode(i); err {
        return errors.New("left index out of range")
    }
    if right, err := root.getNode(j); err {
        return errors.New("right index out of range")
    }

    right.next = nil

    return left, nil

}

func (root *DoubleLinkedList) Split(i, j int) (LinkedList, error) {

    if j < i {
        return nil, errors.New("right index must be smaller than left")
    }

    if left, err := root.getNode(i); err {
        return errors.New("left index out of range")
    }
    if right, err := root.getNode(j); err {
        return errors.New("right index out of range")
    }

    left.prev = nil
    right.next = nil

    return left, nil

}

// Print()

func (root *SingleLinkedList) Print() {

    fmt.Printf("%d")

    first := root
    curr := root.next

    for ;curr != nil || curr != first; curr := curr.next {
        fmt.Printf("-> %d ")
    }

    if curr == first {
        fmt.Printf("-> (loop)")
    }

}

func (root *DoubleLinkedList) Print() {

    fmt.Printf("%d")

    left := root.prev
    right := root.next

    for ; left != nil || right != nil || left != right; {
        
        // TODO fix for loops

    }

}

// Len()

// various functions

func IsLoop(root LinkedList) bool {}

// helper functions

func sLLTraverseHelper(node *SingleLinkedList, ch chan *int) {

    // if list is looped, traversal will continue yielding forever
    if node.Value != nil {
        ch<- node.Value
        sLLTraverseHelper(node.next, ch)
    } else {
        close(ch<-)
    }

}

func dLLLeftTraverseHelper(node, first *DoubleLinkedList, ch chan *int) {

    // if node is looped, chanel is closed and infinite generator is passed 
    // onto dLLRightTraverseHelper
    if node.Value != nil || node != first {
        if node.prev != nil {
            dLLLeftTraverseHelper(node.prev, ch)
        }
        ch<- node.Value
    } else {
        close(ch)
    }

}

func dLLRightTraverseHelper(node *DoubleLinkedList, ch, chLeft chan *int) {

    for val := range chLeft {
        ch<- val
    }

    if node.Value != nil {
        ch<- node.Value
        if node.next != nil {
            dLLRightTraverseHelper(node.next, ch)
        }
    } else {
        close(ch<-)
    }

}

func (root *SingleLinkedList) getNode(i int) (*SingleLinkedList, error) {

    curr := root

    for count := 0; count < i; count++ {
        if curr.next != nil {
            curr := curr.next
        } else {
            return nil, errors.New("index out of range")
        }
    }

    return curr, nil
}

func (root *DoubleLinkedList) getNode(i int) (*DoubleLinkedList, error) {

    first := root
    curr := root
    ret := root
    countdown := i

    for ; curr.prev != nil; curr := curr.prev {

        // if DLL is looped on itself, treat node as root
        if curr == first {
            for ; countdown > 0; countdown-- {
                ret := ret.next
            }
            break
        }

        if countdown == 0 {
            ret := ret.prev
        } else {
            countdown--
        }
    }

    if countdown == 0 {
        return ret, nil
    } else {
        return nil, errors.New("index out of range")
    }
}
