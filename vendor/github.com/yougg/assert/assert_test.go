package assert

import (
	"testing"
)

func TestNew(t *testing.T) {
	a := New(nil)
	if nil != a {
		t.Error("new assert error")
	}
	a = New(t)
	if nil == a {
		t.Error("new assert error")
	}
}

func TestNew2(t *testing.T) {
	a := New(t, true)
	a.Nil(&struct{}{}) // assert fail, but will continue testing
	a.Log("continue testing")
	a.NotNil(t)
}

func TestAssertNil(t *testing.T) {
	New(t).Nil(nil, "check nil input")
}

func TestAssertNotNil(t *testing.T) {
	New(t).NotNil(t, "check non nil input")
}

func TestAssertTrue(t *testing.T) {
	a := New(t)
	a.True(1)
	a.True("t")
	a.True("T")
	a.True("true")
	a.True("TRUE")
	a.True("True")
	a.True(true)
}

func TestAssertFalse(t *testing.T) {
	a := New(t)
	a.False(0)
	a.False("f")
	a.False("F")
	a.False("false")
	a.False("FALSE")
	a.False("False")
	a.False(false)
}

func TestAssertBool(t *testing.T) {
	a := New(t)
	a.Bool(true, 1)
	a.Bool(true, "t")
	a.Bool(true, "T")
	a.Bool(true, "true")
	a.Bool(true, "TRUE")
	a.Bool(true, "True")
	a.Bool(true, true)
	a.Bool(false, 0)
	a.Bool(false, "f")
	a.Bool(false, "F")
	a.Bool(false, "false")
	a.Bool(false, "FALSE")
	a.Bool(false, "False")
	a.Bool(false, false)
}

func TestAssertEqual(t *testing.T) {
	a := New(t)
	a.Equal(nil, nil)
	a.Equal(0, 0)
	a.Equal(1.23456, 1.23456)
	a.Equal(`abc`, "abc")
	a.Equal(true, true)
	a.Equal(struct{}{}, struct{}{})
}

func TestAssertNotEqual(t *testing.T) {
	a := New(t)
	a.NotEqual(nil, t)
	a.NotEqual(0, 0.0)
	a.NotEqual(1.23456, 1.234567)
	a.NotEqual(`abcd`, "abc")
	a.NotEqual(true, false)
	a.NotEqual(struct{}{}, struct{ *Assert }{})
}

func TestAssertContains(t *testing.T) {
	a := New(t)
	l := []int{1, 2, 3}
	a.Contains(l, 1)
	a.Contains(l, 2)
	a.Contains(l, 3)
	l0 := []string{"a", "b", "c"}
	a.Contains(l0, "a")
	a.Contains(l0, "b")
	a.Contains(l0, "c")
}

func TestAssertNotContains(t *testing.T) {
	a := New(t)
	l := []int{1, 2, 3}
	a.NotContains(l, 0)
	a.NotContains(l, 1.1)
	a.NotContains(l, "2")
	l0 := []string{"a", "b", "c"}
	a.NotContains(l0, "d")
	a.NotContains(l0, 'b')
	a.NotContains(l0, 6+7i)
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Error("IsNil return unexpected check result with nil input")
	}
	if IsNil(t) {
		t.Error("IsNil return unexpected check result with non nil input")
	}
}

func TestIsBoolMatch(t *testing.T) {
	if !IsBoolMatch(true, 1) {
		t.Error("IsBoolMatch return unexpected check result with input true, 1")
	}
	if !IsBoolMatch(true, "t") {
		t.Error("IsBoolMatch return unexpected check result with input true, \"t\"")
	}
	if !IsBoolMatch(true, "T") {
		t.Error("IsBoolMatch return unexpected check result with input true, \"T\"")
	}
	if !IsBoolMatch(true, "true") {
		t.Error("IsBoolMatch return unexpected check result with input true, \"true\"")
	}
	if !IsBoolMatch(true, "TRUE") {
		t.Error("IsBoolMatch return unexpected check result with input true, \"TRUE\"")
	}
	if !IsBoolMatch(true, "True") {
		t.Error("IsBoolMatch return unexpected check result with input true, \"True\"")
	}
	if !IsBoolMatch(true, true) {
		t.Error("IsBoolMatch return unexpected check result with input true, true")
	}
	if !IsBoolMatch(false, 0) {
		t.Error("IsBoolMatch return unexpected check result with input false, 0")
	}
	if !IsBoolMatch(false, "f") {
		t.Error("IsBoolMatch return unexpected check result with input false, \"f\"")
	}
	if !IsBoolMatch(false, "F") {
		t.Error("IsBoolMatch return unexpected check result with input false, \"F\"")
	}
	if !IsBoolMatch(false, "false") {
		t.Error("IsBoolMatch return unexpected check result with input false, \"false\"")
	}
	if !IsBoolMatch(false, "FALSE") {
		t.Error("IsBoolMatch return unexpected check result with input false, \"FALSE\"")
	}
	if !IsBoolMatch(false, "False") {
		t.Error("IsBoolMatch return unexpected check result with input false, \"False\"")
	}
	if !IsBoolMatch(false, false) {
		t.Error("IsBoolMatch return unexpected check result with input false, false")
	}
}

func TestIsEqual(t *testing.T) {
	if !IsEqual(nil, nil) {
		t.Error("IsEqual return unexpected check result with input nil")
	}
	if IsEqual(t, nil) {
		t.Error("IsEqual return unexpected check result with input t and nil")
	}
	if !IsEqual(t, t) {
		t.Error("IsEqual return unexpected check result with input t")
	}
	if !IsEqual(0, 0) {
		t.Error("IsEqual return unexpected check result with input 0")
	}
	if !IsEqual(1.2345, 1.2345) {
		t.Error("IsEqual return unexpected check result with input 1.2345")
	}
	if !IsEqual("xyz", `xyz`) {
		t.Error("IsEqual return unexpected check result with input \"xyz\"")
	}
}

func TestIsFunction(t *testing.T) {
	if IsFunction(nil) {
		t.Error("IsFunction return unexpected check result with input nil")
	}
	if !IsFunction(func() {}) {
		t.Error("IsFunction return unexpected check result with input func() {}")
	}
	if !IsFunction(t.Log) {
		t.Error("IsFunction return unexpected check result with input t.Log")
	}
}

func TestIncludeElement(t *testing.T) {
	l := []int{1, 2, 3}
	if ok, found := IncludeElement(l, 1); !ok || !found {
		t.Error("IncludeElement return unexpected check result with input l, 1")
	}
	if ok, found := IncludeElement(l, 0); !ok || found {
		t.Error("IncludeElement return unexpected check result with input l, 0")
	}
	l0 := []string{"a", "b", "c"}
	if ok, found := IncludeElement(l0, "a"); !ok || !found {
		t.Error("IncludeElement return unexpected check result with input l0, 'a'")
	}
	if ok, found := IncludeElement(l0, "d"); !ok || found {
		t.Error("IncludeElement return unexpected check result with input l0, \"d\"")
	}
}
