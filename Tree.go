/*
BIG QUESTIONS:
- Implement interface metods using pointers or references?
 */

package models

import (
    "errors"
    "fmt"
    "math"
    "math/rand"
    "reflect"
    "sync"
    "time"
)

type Tree interface {
    Traverse(int) <-chan int
}

type TreeGenParams interface {
    Generate(int) Tree
}

type TreeInfo interface {
    printTreeInfo()
}

// Define structs and basic functionality

type NTree struct {
    Value *int
    children []*NTree
    size int
}

func (tree NTree) GetChildren() []*NTree {
    return tree.children
}

type NTreeGenParams struct {
    r *rand.Rand
    MaxChild int
    PossibleValues []int
}

func NewNTreeGenParams(r *rand.Rand, maxChild int, possibleValues ...int) TreeGenParams {
    return &NTreeGenParams{
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
    depth int
}

type BSTree struct {
    Value *int
    Right *BSTree
    Left *BSTree
    size int
}

type BSTreeGenParams struct {
    r *rand.Rand
    ChildProb float64
    PossibleValues []int
}

func NewBSTreeGenParams(r *rand.Rand, childProb float64, possibleValues ...int) TreeGenParams {
    return &BSTreeGenParams{
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
    depth int
}

// Printing functions

func (tree *NTree) Print() {
    prefix := "-"
    if tree.size == 1 {
        nTreePrintHelper(tree, &prefix, true)
    } else if tree.size > 1 {
        nTreePrintHelper(tree, &prefix, false)
    }
}

func nTreePrintHelper(node *NTree, prefix *string, last bool) {
	if node.Value != nil {
        fmt.Print(*prefix)
        if last {
            fmt.Print("`-")
            *prefix += "   "
        } else {
            fmt.Print("|-")
            *prefix += "|  "
        }
        fmt.Print(node.Value)
        for i, child := range node.children {
            last := i == len(node.children)
            nTreePrintHelper(child, prefix, last)
        }
    }
}

func (tree *BSTree) Print() {
    prefix := "-"
    if tree.size == 1 {
        bsTreePrintHelper(tree, &prefix, true)
    } else if tree.size > 1 {
        bsTreePrintHelper(tree, &prefix, false)
    }
}
func bsTreePrintHelper(node *BSTree, prefix *string, last bool) {
    if node.Value != nil {
        fmt.Print(*prefix)
        if last {
            fmt.Print("`-")
            *prefix += "   "
        } else {
            fmt.Print("|-")
            *prefix += "|  "
        }
        fmt.Print(node.Value)
        bsTreePrintHelper(node.Right, prefix, false)
        bsTreePrintHelper(node.Right, prefix, true)
    }
}

func (info *nTreeInfo) printTreeInfo() {
    fmt.Println("TO BE IMPLEMENTED")
}

func (info *bsTreeInfo) printTreeInfo() {
    fmt.Println("TO BE IMPLEMENTED")
}

// Traversing functions

// Traverse Binary Search Tree in selected mode

type Order = int
const (
    ordered int = 0
    unordered int = 1
)
const (
    preorder int = -1
    inorder int = 0
    postorder int = 1
)

func (tree *NTree) Traverse(mode Order) <-chan int {

    var tFunc func(*NTree, chan<- int)

    switch mode {
    case ordered:
        tFunc = orderedT
    case unordered:
        tFunc = unorderedT
    default:
        errors.New("Option not available")
    }

    ch := make(chan int, tree.size)
    go tFunc(tree, ch)

    return ch
}

func (tree *BSTree) Traverse(mode Order) <-chan int {

    var tFunc func(*BSTree, chan<- int)

    switch mode {
    case preorder:
        tFunc = preorderT
    case inorder:
        tFunc = inorderT
    case postorder:
        tFunc = postorderT
    default:
        errors.New("Order not available")
    }

    ch := make(chan int, tree.size)
    go tFunc(tree, ch)

    return ch
}

// Equality functions

func (tree0 *NTree) Equals(tree1 *NTree) bool {

    // this comarisson just checks to make sure the values contained 
    // in one tree are present in the other and do not take hierarchy into
    // consideration

    if tree0.size != tree1.size {
        return false
    }

    mapT0 := make(map[int] int)
    mapT1 := make(map[int] int)

    chT0 := make(chan int, tree0.size)
    chT1 := make(chan int, tree1.size)

    go unorderedT(tree0, chT0)
    go unorderedT(tree1, chT1)

    var wg sync.WaitGroup
    wg.Add(2)

    go addToMap(chT0, mapT0, &wg)
    go addToMap(chT1, mapT1, &wg)

    wg.Wait()

    return reflect.DeepEqual(mapT0, mapT1)
}

func addToMap(ch <-chan int, m map[int] int, wg *sync.WaitGroup) {

   defer wg.Done()

   for val := range ch {
        if _, ok := m[val]; ok {
            m[val]++
        } else {
            m[val] = 0
        }
   }

}

func (tree0 *NTree) Identical(tree1 *NTree) bool {

    // This ordered traversal ensures identical structures and values
    if tree0.size != tree1.size {
        return false
    }

    chT0 := make(chan int, tree0.size)
    chT1 := make(chan int, tree1.size)

    go orderedT(tree0, chT0)
    go orderedT(tree1, chT1)

    return compareTrees(chT0, chT1)

}

func (tree0 *BSTree) Equals(tree1 *BSTree) bool {

    // inorder traversal is blind to the tree's structure and only care about
    // the order of values
    if tree0.size != tree1.size {
        return false
    }

    chT0 := make(chan int, tree0.size)
    chT1 := make(chan int, tree1.size)

    go inorderT(tree0, chT0)
    go inorderT(tree1, chT1)

    return compareTrees(chT0, chT1)
}

func (tree0 *BSTree) Identical(tree1 *BSTree) bool {

    // preorder traversal takes exact structure and value order into account
    // for its comparisson and is more efficient than postorder traversals 
    // to spot early deviations 
    if tree0.size != tree1.size {
        return false
    }

    chT0 := make(chan int, tree0.size)
    chT1 := make(chan int, tree1.size)

    go preorderT(tree0, chT0)
    go preorderT(tree1, chT1)

    return compareTrees(chT0, chT1)
}

func compareTrees(chT0, chT1 <-chan int) bool {
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
        default:
            if chT0 == nil && chT1 == nil {
                return true
            }
        }
    }
}

// Info functions

type uniqueCounter struct {
    v map[int] int
    mux sync.Mutex
}

func (unique *uniqueCounter) Add(val int) {

    defer unique.mux.Unlock()

    unique.mux.Lock()
    if _, ok := unique.v[val]; !ok {
        unique.v[val] = 1
    }
}

func (unique *uniqueCounter) toArray() []int {
    var ret []int
    for key := range unique.v {
        ret = append(ret, key)
    }
    return ret
}

func (tree *NTree) getInfo() nTreeInfo {

    unique := &uniqueCounter{
        v: make(map[int] int),
    }
    ch := make(chan nTreeInfo)

    go nTreeInfoHelper(tree, unique, ch)

    info := <-ch
    close(ch)

    info.uniqueValues = unique.toArray()

    return info
}

func nTreeInfoHelper(tree *NTree, unique *uniqueCounter, ch chan nTreeInfo) {

    if tree.size == 0 {
        ch<- nTreeInfo{
            maxChild: 0,
            uniqueValues: nil,
            depth: 0,
        }
    } else if tree.size == 1 {
        unique.Add(*tree.Value)
        ch<- nTreeInfo{
            maxChild: 0,
            uniqueValues: nil,
            depth: 1,
        }
    } else {

        info := nTreeInfo{
            maxChild: len(tree.children),
            uniqueValues: nil,
            depth: 0,
        }

        chChild := make(chan nTreeInfo)

        for _, child := range tree.children {
            go nTreeInfoHelper(child, unique, chChild)
        }

        for i := 0; i < len(tree.children); i++ {
            chInfo := <-chChild
            if chInfo.maxChild > info.maxChild {
                info.maxChild = chInfo.maxChild
            }
            if chInfo.depth > info.depth {
                info.depth = chInfo.depth
            }
        }
        close(chChild)

        unique.Add(*tree.Value)
        info.depth++

        ch<- info

    }

}

func (tree *BSTree) getInfo() bsTreeInfo {

    unique := &uniqueCounter{
        v: make(map[int] int),
    }
    ch := make(chan int)

    go bsTreeInfoHelper(tree, unique, ch)

    info := bsTreeInfo{
        avgNChild: nil,
        uniqueValues: unique.toArray(),
        depth: <-ch,
    }

    // size of current tree divided by the maximum size possible for a root
    // with this depth
    info.avgNChild = float64(tree.size) / (math.Pow(2, float64(info.depth)) - 1)

    return info
}

func bsTreeInfoHelper(tree *BSTree, unique *uniqueCounter, ch chan int) {

    defer func(ch chan int) {
        close(ch)
    }(ch)

    if tree.Value == nil {
        ch<- 0
    } else {
        ret := 0

        chLeft := make(chan int)
        chRight := make(chan int)

        go bsTreeInfoHelper(tree.Left, unique, chLeft)
        go bsTreeInfoHelper(tree.Right, unique, chRight)

        if val := <-chLeft; ret < val {
                ret = val
        }
        if val := <-chRight; ret < val {
                ret = val
        }

        unique.Add(*tree.Value)
        ch<- ret + 1
    }
}

// Generate functions

// These imply that trees are also parameters for generating resampled 
// versions of themselves

func (params *NTreeGenParams) Generate(depth int) Tree {
    /* NTreeGenParams
    r *Rand
    MaxChild int
    PossibleValues []int
    */
    ch := make(chan *NTree)
    go nTreeGenHelper(params, depth, ch)
    return <-ch
}

func nTreeGenHelper(params *NTreeGenParams, depth int, ch chan *NTree) {

    tree := &NTree{
        Value:nil,
        children:nil,
        size:0,
    }

    defer func(tree *NTree, ch chan *NTree) {
        ch<- tree
    }(tree, ch)

    if depth >= 1 {
        tree.Value = &params.PossibleValues[params.r.Intn(len(params.PossibleValues))]
        var children []*NTree
        tree.children = children
        tree.size = 1
    }
    if depth > 1 {

        nChild := params.r.Intn(params.MaxChild)

        chChild := make(chan *NTree)
        defer close(chChild)

        for i := 0; i < nChild; i++ {
            go nTreeGenHelper(params, depth-1, chChild)
        }

        for i := 0; i < nChild; i++ {
            child := <-chChild
            tree.size += child.size
            tree.children = append(tree.children, child)
        }
    }
}

func (tree *NTree) GetParams() NTreeGenParams {

    info := tree.getInfo()
    return NTreeGenParams{
        r: rand.New(rand.NewSource(time.Now().UnixNano())),
        MaxChild: info.maxChild,
        PossibleValues: info.uniqueValues,
    }
}

func (params *BSTreeGenParams) Generate(depth int) Tree {
    /*
    r *Rand
    ChildProb float64
    PossibleValues []int
    */
    ch := make(chan *BSTree)
    go bsTreeGenHelper(params, depth, ch)
    return <-ch
}

func bsTreeGenHelper(params *BSTreeGenParams, depth int, ch chan *BSTree) {

        tree := &BSTree{
            Value:nil,
            Right:nil,
            Left:nil,
            size:0,
        }

        defer func(tree *BSTree, ch chan *BSTree) {
            ch<- tree
            close(ch)
        }(tree, ch)

        if depth >= 1 {
            tree.Value = &params.PossibleValues[params.r.Intn(len(params.PossibleValues))]
            tree.size = 1
        }
        if depth > 1 {

            chRight := make(chan *BSTree)
            chLeft := make(chan *BSTree)

            if params.r.Float64() < params.ChildProb {
                go bsTreeGenHelper(params, depth-1, chRight)
            } else {
                go bsTreeGenHelper(params, 0, chRight)
            }

            if params.r.Float64() < params.ChildProb {
                go bsTreeGenHelper(params, depth-1, chLeft)
            } else {
                go bsTreeGenHelper(params, 0, chLeft)
            }

            tree.Right = <-chRight
            tree.Left = <-chLeft

            tree.size += (tree.Right.size + tree.Left.size)
        }
}

func (tree *BSTree) GetParams() BSTreeGenParams {
    info := tree.getInfo()
    return BSTreeGenParams{
        r: rand.New(rand.NewSource(time.Now().UnixNano())),
        ChildProb: info.avgNChild,
        PossibleValues: info.uniqueValues,
    }
}

// Helper functions

func orderedT(tree *NTree, ch chan<- int) {

    if tree.size >= 1 {
        ch<- *tree.Value
    }
    if tree.size > 1 {

        var chChildren []chan int

        for i, child := range tree.children {
            chChildren = append(chChildren, make(chan int, child.size))
            go orderedT(child, chChildren[i])
        }

        for i, chChannel := range chChildren {
            for j := 0; j < tree.children[i].size; j++ {
                ch<- <-chChannel
            }
            close(chChannel)
        }

    }

}

func unorderedT(tree *NTree, ch chan<- int) {

    if tree.size >= 1 {
        ch<- *tree.Value
    }
    if tree.size > 1 {
        for _, child := range tree.children {
            go unorderedT(child, ch)
        }
    }
}

func inorderT(tree *BSTree, ch chan<- int) {

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    for i := 0; i < tree.Left.size; i++ {
        ch<- <-chLeft
    }
    close(chLeft)

    ch<- *tree.Value

    for i := 0; i < tree.Right.size; i++ {
        ch<- <-chRight
    }
    close(chRight)
}

func preorderT(tree *BSTree, ch chan<- int) {

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    ch<- *tree.Value
    for i := 0; i < tree.Left.size; i++ {
        ch<- <-chLeft
    }
    close(chLeft)
    for i := 0; i < tree.Right.size; i++ {
        ch<- <-chRight
    }
    close(chRight)
}

func postorderT(tree *BSTree, ch chan<- int) {

    if tree.Value == nil {
        return
    }

    chLeft := make(chan int, tree.Left.size)
    chRight := make(chan int, tree.Right.size)

    go inorderT(tree.Left, chLeft)
    go inorderT(tree.Right, chRight)

    for i := 0; i < tree.Left.size; i++ {
        ch<- <-chLeft
    }
    close(chLeft)
    for i := 0; i < tree.Right.size; i++ {
        ch<- <-chRight
    }
    close(chRight)
    ch<- *tree.Value

}
