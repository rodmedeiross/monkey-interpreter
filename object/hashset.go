package object

import "hash/fnv"

var (
	true  = HashSet{ObjectType: BOOLEAN_OBJ, Value: uint64(1)}
	false = HashSet{ObjectType: BOOLEAN_OBJ, Value: uint64(0)}
)

type HashSet struct {
	ObjectType ObjectType
	Value      uint64
}

type Hashable interface {
	Hash() HashSet
}

func (b *Boolean) Hash() HashSet {
	if b.Value {
		return true
	} else {
		return false
	}
}

func (i *Integer) Hash() HashSet {
	return HashSet{ObjectType: INTEGER_OBJ, Value: uint64(i.Value)}
}

func (i *String) Hash() HashSet {
	hash := fnv.New64a()
	hash.Write([]byte(i.Value))

	return HashSet{ObjectType: INTEGER_OBJ, Value: hash.Sum64()}
}
