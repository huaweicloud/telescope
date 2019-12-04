package mockfn

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type (
	// function is an applied function
	//  needed to undo a function
	function struct {
		originalBytes      []byte
		replacement        *reflect.Value
		aliasMockedPos     uintptr
		aliasOriginalBytes []byte
		addr               *uintptr
	}

	value struct {
		_   uintptr
		ptr unsafe.Pointer
	}

	FuncGuard struct {
		target      reflect.Value
		replacement reflect.Value
		alias       *reflect.Value // Use this interface to access the original target
	}
)

var (
	lock = sync.Mutex{}

	funcs = make(map[uintptr]function)
)

func (f *FuncGuard) Revert() {
	revertValue(f.target)
}

func (f *FuncGuard) Restore() {
	mockValue(f.target, f.replacement, f.alias)
}

func Replace(target, replacement interface{}) *FuncGuard {
	t := reflect.ValueOf(target)
	r := reflect.ValueOf(replacement)
	mockValue(t, r, nil)

	return &FuncGuard{t, r, nil}
}

// ReplaceEx replaces a function with another
//  alias: A wrapper, to access the original target when mocked
func ReplaceEx(target, alias, replacement interface{}) *FuncGuard {
	t := reflect.ValueOf(target)
	r := reflect.ValueOf(replacement)
	a := reflect.ValueOf(alias)
	mockValue(t, r, &a)

	return &FuncGuard{t, r, &a}
}

// ReplaceInstanceMethod replaces an instance method methodName for the type target with replacement
//  Replacement should expect the receiver (of type target) as the first argument
func ReplaceInstanceMethod(target reflect.Type, methodName string, replacement interface{}) *FuncGuard {
	m, ok := target.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("unknown method %s", methodName))
	}
	r := reflect.ValueOf(replacement)
	mockValue(m.Func, r, nil)

	return &FuncGuard{m.Func, r, nil}
}

// ReplaceInstanceMethod replaces an instance method methodName for the type target with replacement
//  Replacement should expect the receiver (of type target) as the first argument
func ReplaceInstanceMethodEx(target reflect.Type, methodName string, alias, replacement interface{}) *FuncGuard {
	m, ok := target.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("unknown method %s", methodName))
	}
	r := reflect.ValueOf(replacement)
	a := reflect.ValueOf(alias)
	mockValue(m.Func, r, &a)

	return &FuncGuard{m.Func, r, &a}
}

func mockValue(target reflect.Value, replacement reflect.Value, alias *reflect.Value) {
	lock.Lock()
	defer lock.Unlock()

	if target.Kind() != reflect.Func {
		panic("target has to be a Func")
	}

	if replacement.Kind() != reflect.Func {
		panic("replacement has to be a Func")
	}

	if target.Type() != replacement.Type() {
		panic(fmt.Sprintf("target and replacement have to have the same type %s != %s", target.Type(), replacement.Type()))
	}

	if alias != nil {
		if alias.Kind() != reflect.Func {
			panic("alias has to be a Func")
		}

		if target.Type() != alias.Type() {
			panic(fmt.Sprintf("target and alias have to have the same type %s != %s", target.Type(), alias.Type()))
		}
	}

	if fn, ok := funcs[target.Pointer()]; ok {
		revert(target.Pointer(), fn)
	}

	var addr *uintptr
	var aliasPtr uintptr
	var aliasBytes []byte
	if alias != nil {
		targetOffset, aliasOffset, aliasOriginal := replaceJBE(target.Pointer(), (*alias).Pointer())

		addr = new(uintptr)
		*addr = *(*uintptr)(pointer(target)) + targetOffset
		aliasPos := (*alias).Pointer() + aliasOffset
		originalBytes := replaceFunction(aliasPos, (uintptr)(unsafe.Pointer(addr)))
		aliasPtr = (*alias).Pointer()
		aliasBytes = make([]byte, len(aliasOriginal)+len(originalBytes))
		copy(aliasBytes, aliasOriginal)

		capacity := len(aliasBytes)
		len1 := len(aliasOriginal)
		for i := len1; i < capacity; i++ {
			aliasBytes[i] = originalBytes[i-len1]
		}
	}

	originalBytes := replaceFunction(target.Pointer(), (uintptr)(pointer(replacement)))
	funcs[target.Pointer()] = function{originalBytes, &replacement, aliasPtr, aliasBytes, addr}
}

func pointer(v reflect.Value) unsafe.Pointer {
	return (*value)(unsafe.Pointer(&v)).ptr
}

// Revert removes any mock funcs on target
// returns whether target was mocked in the first place
func Revert(target interface{}) bool {
	return revertValue(reflect.ValueOf(target))
}

// RevertInstanceMethod removes the function on methodName of the target
//  returns whether it was mocked in the first place
func RevertInstanceMethod(target reflect.Type, methodName string) bool {
	m, ok := target.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("unknown method %s", methodName))
	}
	return revertValue(m.Func)
}

// RevertAll removes all applied mock functions
func RevertAll() {
	lock.Lock()
	defer lock.Unlock()
	for target, p := range funcs {
		revert(target, p)
		delete(funcs, target)
	}
}

// Revert removes a function from the specified function
//  returns whether the function was mock in the first place
func revertValue(target reflect.Value) bool {
	lock.Lock()
	defer lock.Unlock()
	fn, ok := funcs[target.Pointer()]
	if !ok {
		return false
	}
	revert(target.Pointer(), fn)
	delete(funcs, target.Pointer())
	return true
}

func revert(target uintptr, p function) {
	copyToLocation(target, p.originalBytes)
	if p.addr != nil {
		copyToLocation(p.aliasMockedPos, p.aliasOriginalBytes)
	}
}
