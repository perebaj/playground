package main

import "slices"

func main() {

}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isSameTree(p *TreeNode, q *TreeNode) bool {
	pSlice := inorderSlice(p)
	qSlice := inorderSlice(q)

	if !slices.Equal(qSlice, pSlice) {
		return false
	}
	return true
}

func inorderSlice(n *TreeNode) []int {
	var resp []int
	if n != nil {
		resp = append(resp, inorderSlice(n.Left)...)
		resp = append(resp, n.Val)
		resp = append(resp, inorderSlice(n.Right)...)
	}
	return resp
}
