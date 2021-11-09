package lang

import "fmt"

type ObjectType uint8

const (
	ObjNil ObjectType = iota
	ObjStr
	ObjInt
	ObjBool
	ObjInv
)

type Object struct {
	Type  ObjectType
	StrV  string
	IntV  int
	BoolV bool
}

func (o *Object) Repr() string {
	switch o.Type {
	case ObjNil:
		return "[Nil]"
	case ObjStr:
		return "[Str `" + o.StrV + "`]"
	case ObjInt:
		return fmt.Sprintf("[Int %d]", o.IntV)
	case ObjBool:
		return fmt.Sprintf("[Bool %t]", o.BoolV)
	}
	return "?"
}

func (o *Object) String() string {
	switch o.Type {
	case ObjNil:
		return "Nil"
	case ObjStr:
		return o.StrV
	case ObjInt:
		return fmt.Sprint(o.IntV)
	case ObjBool:
		return fmt.Sprint(o.BoolV)
	}
	return "?"
}
