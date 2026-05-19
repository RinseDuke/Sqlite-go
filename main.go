package main

import (
	"encoding/binary"
	"fmt"
)

// B+tree page size, the actual size of the node is less than this value, because we need to store some metadata
const BTREE_PAGE_SIZE = 4069
const BTREE_MAX_KEYS_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

// B+tree Node
type Node struct {
	key      [][]byte //two-dimensional array, each row is a key
	vals     [][]byte //two-dimensional array,each row is a value
	children []*Node
}

// BNODE(a decode bnode tree)
type BNode []byte //can be drupmed to the disk

// getter
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// setter
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)

}

// Encode the node into a byte array
func Encode(node *Node) []byte { return []byte{} }

// Decode the byte array into a node
func Decode(page []byte) *Node {
	return &Node{}
}

func main() {
	fmt.Println("Hello,B+tree!")
}
