package main

import "fmt"

// Index returns the index of x in s, or -1 if not found.
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

// any is an alias for interface{}, meaning any type.
type List[T any] struct {
	elements []T
}

type Number interface {
	int | int64 | float64
}

func Sum[T Number](a, b T) T {
	return a + b
}

// Implement LinkedList using generics
type Node[T any] struct {
	value T
	next  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
}

func (ll *LinkedList[T]) Append(value T) {
	newNode := &Node[T]{value: value}
	if ll.tail != nil {
		ll.tail.next = newNode
	} else {
		ll.head = newNode
	}
	ll.tail = newNode
}

func (ll *LinkedList[T]) Print() {
	current := ll.head
	for current != nil {
		fmt.Print(current.value, " ")
		current = current.next
	}
	fmt.Println()
}

func main() {
	// Index works on a slice of ints
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	// Index also works on a slice of strings
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))

	intList := List[int]{elements: []int{1, 2, 3, 4, 5}}
	fmt.Println(intList)

	stringList := List[string]{elements: []string{"a", "b", "c"}}
	fmt.Println(stringList)

	fmt.Println(Sum(10, 20))
	fmt.Println(Sum(10.5, 20.3))

	// Demonstrate LinkedList with generics
	intLinkedList := LinkedList[int]{}
	intLinkedList.Append(1)
	intLinkedList.Append(2)
	intLinkedList.Append(3)
	intLinkedList.Print()

	stringLinkedList := LinkedList[string]{}
	stringLinkedList.Append("a")
	stringLinkedList.Append("b")
	stringLinkedList.Append("c")
	stringLinkedList.Print()
}
