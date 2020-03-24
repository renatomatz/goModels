package models

import (
    "fmt"
    "sync"
    "math/rand"
)

type Tree interface {
    Print()
    Traverse(int) []int
    Equals(Tree) bool
    Identical(Tree) bool
}

type TreeGenParams interface {
    Generate(int) Tree
}

type TreeInfo interface {
    printTreeInfo() error
}

// Define structs and basic functionality

type NTree struct {
    Value int
    children Children
    size int
}
type Children = []Ntree

type nTreeGenParams struct {
    r *Rand
    Balanced bool
    MaxChild int
    PossibleValues []int
}

func NewNTreeGenParams(r *Rand, maxChild int, possibleValues ...int) TreeGenParams {
    return nTreeGenParams{
        r: r,
        MaxChild: maxChild,
        PossibleValues: possibleValues,
    }
}

func NewDefaultNTreeGenParams() TreeGenParams {
    return NewNTreeGenParams(rand.New(rand.NewSource(42)), 5, []int{1, 2, 3, 4, 5}...)
}

type nTreeInfo struct {
    maxChild int
    uniqueValues []int
    maxDepth int
}

type BSTree struct {
    Value int
    Right Tree
    Left Tree
    size int
}

type bSTreeGenParams struct {
    r *Rand
    ChildProb float64
    PossibleValues []int
}

func NewBSTreeGenParams(r *Rand, childProb int, possibleValues ...int) TreeGenParams {
    return bSTreeGenParams{
        r: r,
        ChildProb: childProb,
        PossibleValues: possibleValues,
    }
}

func NewDefaultBSTreeGenParams() TreeGenParams {
    return NewBSTreeGenParams(rand.New(rand.NewSource(42)), 0.5, []int{1, 2, 3, 4, 5}...)
}

type bsTreeInfo struct {
    avgNChild float64
    uniqueValues []int
    maxDepth int
}

// Printing functions

func (tree *NTree) Print() {
    fmt.Println("TO BE IMPLEMENTED")
}

func (tree *BSTree) Print() {
    fmt.Println("TO BE IMPLEMENTED")
}

func PrintTreeInfo(info treeInfo) {
    info.printTreeInfo()
}

func (info *nTreeInfo) printTreeInfo() (err error) {

}

func (info *bSTreeInfo) printTreeInfo() (err error) {

}

// Traversing functions

// Traverse Binary Search Tree in selected mode

const (
    preorder int = -1
    inorder int = 0
    postorder int = 1
)

func (tree *NTree) Traverse(mode int) []int {
    switch mode {
    case preorder:
    case inorder:
    case postorder:
    default:
    }
}

func (tree *BSTree) Traverse(mode int) []int {
    switch mode {
    case preorder:
    case inorder:
    case postorder:
    default:
    }
}

// Equality functions

func (tree0 *NTree) Equals(tree1 *NTree) bool {

    // this comarisson just checks to make sure the values contained 
    // in one tree are present in the other and do not take hierarchy into
    // consideration

}

func (tree0 *NTree) Identical(tree1 *NTree) bool {

    // This ordered traversal ensures identical structures and values
    return compareTrees(tree0, tree1, orderedT)

}

func (tree0 *BSTree) Equals(tree1 *BSTree) bool {

    // inorder traversal is blind to the tree's structure and only care about
    // the order of values
    return compareTrees(tree0, tree1, inorderT)

}

func (tree0 *BSTree) Identical(tree1 *BSTree) bool {

    // preorder traversal takes exact structure and value order into account
    // for its comparisson and is more efficient than postorder traversals 
    // to spot early deviations 
    return compareTrees(tree0, tree1, preorderT)

}

func compareTrees(tree0, tree1 Tree, traversalFunc func(Tree, chan<- int)) bool {

    if tree0.size != tree1.size {
        return false
    }

    chT0 := make(chan int, tree0.size)
    chT1 := make(chan int, tree1.size)

    go traversalFunc(tree0, chT0)
    go traversalFunc(tree1, chT1)

    for {
        select {
        case val0, ok := <-chT0:
            if !ok {
                chT0 = nil
            } else if chT1 == nil {
                return false
            } else if val1, ok := <-chT1; !ok || val0 != val1 {
                return false
            }
        case val1, ok := <-chT1:
            if !ok {
                chT1 = nil
            } else if chT0 == nil {
                return false
            } else if val0, ok := <-chT0; !ok || val1 != val0 {
                return false
            }
        }

        if chT0 == nil && chT1 == nil {
            break
        }
    }

    return true
}

// Generate functions

// These imply that trees are also parameters for generating bootstrapped 
// versions of themselves

func (params *NTreeGenParams) Generate(maxDepth int) Tree {

}

func (tree *NTree) GetParams(maxDepth int) NTreeGenParams {

}

func (params *BSTreeGenParams) Generate(maxDepth int) Tree {

}

func (tree *NTree) GetParams(maxDepth int) BSTreeGenParams {

}

// Helper functions

func getTreeInfo(tree *Tree) (treeInfo, error) {

    switch tree.(type) {

    }

}

func orderedT(tree *NTree, ch chan<- int) {

    defer func(ch chan<- int) {
        close(ch)
    }(ch)

    if tree.size == 1 {
        ch<- tree.Value
    } else if tree.size >= 2 {

        var chChildren []chan int

        for i, child := range tree.children {
            chChildren.append(make(chan int, child.size))
            go orderedT(child, chChildren[i])
        }

        ch<- tree.Value
        for _, chChannel := range chChildren {
            for val := range chChannel {
                ch<- val
            }
        }

    }

}

func unorderedT(tree *Ntree, wg *sync.WaitGroup, ch chan<- int) {

    defer wg.Done()

    switch size := tree.size; {
    case size >= 1:
        ch<- tree.Value
        fallthrough
    case size >= 2:
        for _, child := range tree.children {
            wg.Add(1)
            go unorderedT(child, wg, ch)
        }

    }

}

func inorderT(tree *BSTree, ch chan<- int) {

    defer func(ch chan int) {
        close(ch)
    }(ch)

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    for val := range chLeft {
        ch<- val
    }

    ch<- tree.Value

    for val := range chRight {
        ch<- val
    }
}

func preorderT(tree *BSTree, ch chan<- int) {

    defer func(ch chan int) {
        close(ch)
    }(ch)

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    ch<- tree.Value
    for val := range chLeft {
        ch<- val
    }
    for val := range chRight {
        ch<- val
    }
}

func postorderT(tree *BSTree, ch chan<- int) {

    defer func(ch chan int) {
        close(ch)
    }(ch)

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    for val := range chLeft {
        ch<- val
    }
    for val := range chRight {
        ch<- val
    }
    ch<- tree.Value

}

