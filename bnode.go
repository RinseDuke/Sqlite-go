package main

import (
	"bytes"
	"encoding/binary"
)

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

// B+tree page size, the actual size of the node is less than this value, because we need to store some metadata
const BTREE_PAGE_SIZE = 4069
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

// B+tree Node
type Node struct {
	key      [][]byte //two-dimensional array, each row is a key
	vals     [][]byte //two-dimensional array,each row is a value
	children []*Node
}

// BNODE(a decode bnode tree)
type BNode []byte //can be drupmed to the disk

func init() {
	node1max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max >= BTREE_PAGE_SIZE {
		panic("node size exceeds page size")
	}
}

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
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + idx*8
	return binary.LittleEndian.Uint64(node[pos:])
	//return the value of the pointer at the position of idx
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx > node.nkeys() {
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
	//return the value of the offset at the position of idx key and value array
}
func (node BNode) setOffset(idx uint16, val uint16) {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	// binary.LittleEndian.PutUint64(node[pos:],val)
	// offset is 2 bytes,so we need to write the value in 2 bytes(uint16)
	binary.LittleEndian.PutUint16(node[pos:], val)
}

// get the position of the key and value in the node
func (node BNode) KvPos(idx uint16) uint16 {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

// get the key of the key
func (node BNode) getkey(idx uint16) []byte {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:]) //the first 2 bytes is the length of the key
	return node[pos+4:][:klen]
}

// get the value of the value
func (node BNode) getVal(idx uint16) []byte {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}

// Encode the node into a byte array
func Encode(node *Node) []byte { return []byte{} }

// Decode the byte array into a node
func Decode(page []byte) (*Node, error) {
	return &Node{}, nil
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	//directly return the position of the total nkeys
	return node.KvPos(node.nkeys()) // uses the offset value of the last key
}

// idx is the position of the item (a key, a value or a pointer).
// idx 是条目（键、值或指针）的位置。
// ptr is the nth child pointer, which is unused for leaf nodes.
// ptr 是第 n 个子节点指针，对于叶节点来说未使用。
// key and val is the KV pair. Use an empty value for internal nodes.
// key 和 val 是键值对。对于内部节点，请使用空值。
func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	//Ptrs
	new.setPtr(idx, ptr)
	//KVkeys
	pos := new.KvPos(idx) //use offset to get the position of the key and value
	//4 bytes KV size
	binary.LittleEndian.PutUint16(new[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
	//KV Data
	copy(new[pos+4:], key)
	//uint16(len(key)) the goal is to protect the value type is uint16 not int
	copy(new[pos+4+uint16(len(key)):], val)
	//update the offset of the next key and value
	//new.getOffset(idx)+4+uint16(len(key)+len(val)) is key ang value size
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

// nodeAppendRange() is just a loop for copying keys, values, and pointers:
// copy multiple keys, values, and pointers into the position
func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	//this for range just count time of the loop
	for i := uint16(0); i < n; i++ {
		//dstNew is the position of the new node,srcOld is the position of the old node
		//i is the offset of the key and value in the node
		dst, src := dstNew+i, srcOld+i
		nodeAppendKV(new, dst, old.getPtr(src), old.getkey(src), old.getVal(src))
	}
}

// leaf insert
func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	//copy before inserted item
	nodeAppendRange(new, old, 0, 0, idx)
	//inserted item
	nodeAppendKV(new, idx, 0, key, val)
	//copy after inserted item
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// leaf update
func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys())
	//copy before updated item
	nodeAppendRange(new, old, 0, 0, idx)
	//updated item
	nodeAppendKV(new, idx, 0, key, val)
	//copy after updated item
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-(idx+1))
}

// find key
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	for i := uint16(0); i < nkeys; i++ {
		cmp := bytes.Compare(node.getkey(i), key)
		if cmp == 0 {
			return i
		}
		if cmp > 0 {
			return i - 1
		}
	}
	return nkeys - 1
}

// insert or update after a key lookup
func nodeInsertOrUpdate(new BNode, old BNode, key []byte, val []byte) {
	idx := nodeLookupLE(old, key)
	if bytes.Equal(old.getkey(idx), key) {
		leafUpdate(new, old, idx, key, val)
	} else {
		leafInsert(new, old, idx+1, key, val) //insert after the position of the node
	}
}

// I dont sure this part is right for bnode.go,it is about the split of the node
// this is original content about it
// For an in-memory B+tree, an oversized node can be split into 2 nodes
// each with half of the keys. For a disk-based B+tree
// half of the keys may not fit into a page due to uneven key sizes.
// However, we can use the half position as an initial guess
// then move it left or right if the half is too large.
// 对于内存中的 B+树，过大的节点可以拆分为两个节点，每个节点包含一半的键。
// 而对于基于磁盘的 B+树，由于键大小不均，一半的键可能无法放入一个页面。
// 不过，我们可以将中间位置作为初始猜测，如果该位置过大，则向左或向右移动。
func nodeSplit2(left BNode, right BNode, old BNode) {

}
