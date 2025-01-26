package error_three

import (
	"fmt"
	"strings"
)

type dummy struct{}
type wrapper interface {
	Unwrap() error
}
type listWrapper interface {
	Unwrap() []error
}

func PrintErrorThree(err error) string {
	visited := map[error]dummy{}
	errStack := &errorStack{}
	errStack.Push(err, 0)
	res := ""
	dfsError(errStack, visited, func(node *stackNode) {
		tabsCount := node.tabs + 1
		tabs := strings.Repeat("--", tabsCount)

		res += fmt.Sprintf("%s "+
			"%d: %s\n", tabs, tabsCount, node.err)
	})
	return res
}

type stackNode struct {
	err      error
	previous *stackNode
	tabs     int
}

type errorStack struct {
	head *stackNode
	size int
}

func (es *errorStack) Push(err error, tabs int) {
	es.head = &stackNode{err, es.head, tabs}
	es.size++
}

func (es *errorStack) Pop() *stackNode {
	if es.size == 0 {
		return nil
	}
	res := es.head
	es.head = es.head.previous
	es.size--
	return res
}

func dfsError(stack *errorStack, visited map[error]dummy, printer func(node *stackNode)) {
	nextNode := stack.Pop()
	for nextNode != nil {
		printer(nextNode)
		unwrapErr(stack, nextNode, visited)
		nextNode = stack.Pop()
	}
}

func unwrapErr(stack *errorStack, curNode *stackNode, visited map[error]dummy) {
	wrap, ok := curNode.err.(wrapper)
	if ok {
		child := wrap.Unwrap()
		describeNode(stack, child, visited, curNode.tabs+1)
	}
	listWrap, ok := curNode.err.(listWrapper)
	if ok {
		for _, child := range listWrap.Unwrap() {
			describeNode(stack, child, visited, curNode.tabs+1)
		}
	}
	return
}

func describeNode(stack *errorStack, childErr error, visited map[error]dummy, tabs int) {
	if childErr == nil {
		return
	}
	_, isVisited := visited[childErr]
	if isVisited {
		return
	}
	visited[childErr] = dummy{}
	stack.Push(childErr, tabs)
}
