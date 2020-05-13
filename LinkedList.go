package models

import = (
    "errors"
    "fmt"
    "sync"
}

type LinkedList interface {
	Iter() <-chan *int
	Get(int) (*int, error)
        Append(int) error
	Remove(int) (*int, error)
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

// Iter() <-chan int

func (root *SingleLinkedList) Iter() <-chan *int {

    ch = make(chan *int)

    go sLLIterHelper(root, ch)
    return ch

}

func (root *DoubleLinkedList) Iter() <-chan *int {

    chLeft = make(chan *int)
    chRight = make(chan *int)

    var after int
    if root.prev != nil {
        after = root.Value
        go dLLLeftIterHelper(root.prev, root, chLeft)
    }

    go dLLRightIterHelper(root, chRight, chLeft)
    return chRight

}

// ToArray() []*int

func (root *SingleLinkedList) ToArray() {

    var ret []*int
    for val := range root.Iter() {
        ret = append(ret, val)
    }

    return ret

}

func (root *DoubleLinkedList) ToArray() {

    var ret []*int
    for val := range root.Iter() {
        ret = append(ret, val)
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

// Remove(int) *int

func (root *SingleLinkedList) Remove(int) (*int, error) {

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

func (root *SingleLinkedList) Remove(int) (*int, error) {

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

    fmt.Printf("(%d)", root.Value)

    first := root
    curr := root.next

    for ;curr != nil || curr != first; curr = curr.next {
        fmt.Printf("->(%d)", curr.Value)
    }

    if curr == first {
        fmt.Printf("->(LOOP)")
    } else {
        fmt.Printf("->(END)")
    }

}

func (root *DoubleLinkedList) Print() {

    leftNode := root
    rightNode := root

    leftWing := LIFO{}
    rightWing := FIFO{}

    end := "END"

    for ;; {

        if leftNode.prev == rightNode {
            end = "LOOP"
            break
        } else if leftNode.prev != nil {
            leftNode := leftNode.prev
            leftWing.Push(leftNode.Value)
        }

        if rightNode.next == leftNode {
            loop = true
            break
        } else if rightNode.next != nil {
            rightNode := rightNode.next
            rightWing.Push(rightNode.Value
        }

        if leftNode.prev == nil && rightNode.next == nil {
            break
        }

    }

    fmt.Print("(%s)", end)

    for ; !leftWing.IsEmpty(); {
        fmt.Printf("<-(%d)", leftWing.Pop())
    }

    fmt.Printf("<-(%d)->", root.Value)

    for ; !rightWing.IsEmpty(); {
        fmt.Printf("(%d)->", rightWing.Pop())
    }

    fmt.Print("(%s)", end)

}

// Len() int

func (root *SingleLinkedList) Len() int {

    curr := root
    counter := int(root.Value =! nil)

    for ; curr.next != nil || curr.next != root; curr = curr.next {
        counter++
    }

    if curr.next == root {
        return -1
    else {
        return counter
    }

}

func (root *DoubleLinkedList) Len() int {

    leftNode := root
    rightNode := root
    counter := int(root.Value =! nil)

    for ;; {

        if leftNode.prev == rightNode {
            break
        } else if leftNode.prev != nil {
            counter++
        }

        if rightNode.next == leftNode {
            break
        } else if rightNode.next != nil {
            counter++
        }

        if leftNode.prev == nil && rightNode.next == nil {
            break
        }

    }

    if leftNode.prev == rightNode || rightNode.next == leftNode {
        return -1
    } else {
        return counter
    }

}

// helper functions

func sLLIterHelper(node *SingleLinkedList, ch chan *int) {

    // if list is looped, traversal will continue yielding forever
    if node.Value != nil {
        ch<- node.Value
        sLLIterHelper(node.next, ch)
    } else {
        close(ch<-)
    }

}

func dLLLeftIterHelper(node, first *DoubleLinkedList, ch chan *int) {

    // if node is looped, chanel is closed and infinite generator is passed 
    // onto dLLRightIterHelper
    if node.Value != nil || node != first {
        if node.prev != nil {
            dLLLeftIterHelper(node.prev, ch)
        }
        ch<- node.Value
    } else {
        close(ch)
    }

}

func dLLRightIterHelper(node *DoubleLinkedList, ch, chLeft chan *int) {

    for val := range chLeft {
        ch<- val
    }

    if node.Value != nil {
        ch<- node.Value
        if node.next != nil {
            dLLRightIterHelper(node.next, ch)
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
