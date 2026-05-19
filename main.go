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

// read and write the child pointer array
func (node BNode) getPtr(idx uint16) uint64 {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + idx*8
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + idx*8
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// read the 'offset' of array
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

// get the position of the key and value in the node
func (node BNode) KvPos(idx uint16) uint16 {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getkey(idx uint16) []byte {
	if idx >= node.nkeys() {
		panic("index out of range")
	}
	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
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
