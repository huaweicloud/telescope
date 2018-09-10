package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yougg/mockfn"
	"testing"
	"time"
)

//InitConfig
func Test_InitConfig1(t *testing.T) {
	//go InitConfig()
	Convey("Test_InitConfig", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})

		})
	})
}

//Code
func TestCode(t *testing.T) {
	//go InitConfig()
	Convey("Test_Code", t, func() {
		Convey("test case 1", func() {
			field := errInvalidField{}
			code := field.Code()
			So(code, ShouldBeBlank)
		})
	})
}

//Msg
func TestMsg(t *testing.T) {
	//go InitConfig()
	Convey("Test_Msg", t, func() {
		Convey("test case 1", func() {
			field := errInvalidField{}
			msg := field.Msg()
			So(msg, ShouldEqual, ": .")
		})
	})
}

//Error
func TestErrorErrInvalidField(t *testing.T) {
	//go InitConfig()
	Convey("Test_errInvalidField", t, func() {
		Convey("test case 1", func() {
			field := errInvalidField{}
			err := field.Error()
			So(err, ShouldEqual, ": ")
		})
	})
}

//Field
func TestField(t *testing.T) {
	//go InitConfig()
	Convey("Test_Field", t, func() {
		Convey("test case 1", func() {
			field := errInvalidField{object: "abc"}
			f := field.Field()
			So(f, ShouldEqual, "abc.")
		})
	})
}

//SetObject
func TestSetObject(t *testing.T) {
	//go InitConfig()
	Convey("Test_SetObject", t, func() {
		Convey("test case 1", func() {
			field := errInvalidField{object: "abc"}
			field.SetObject("b")
			So(field.object, ShouldEqual, "b")
		})
	})
}

//Add
func TestAdd(t *testing.T) {
	//go InitConfig()
	Convey("Test_Add", t, func() {
		Convey("test case 1", func() {
			/*	fields := &ErrInvalidFields{}
				field := ErrInvalidField{}
				fields.Add(field)*/
		})
	})
}

//Len
func TestLen(t *testing.T) {
	//go InitConfig()
	Convey("Test_Len", t, func() {
		Convey("test case 1", func() {
			fields := &ErrInvalidFields{}
			i := fields.Len()
			So(i, ShouldEqual, 0)
		})
	})
}

//Error
func TestError(t *testing.T) {
	//go InitConfig()
	Convey("Test_Error", t, func() {
		Convey("test case 1", func() {
			defer mockfn.RevertAll()
			mockfn.Replace(time.Sleep, func(d time.Duration) {
				return
			})
			fields := &ErrInvalidFields{
				Errs: []ErrInvalidField{},
			}
			s := fields.Error()
			So(s, ShouldBeBlank)
		})
	})
}

//NewErrFieldRequired
func TestNewErrFieldRequired(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewErrFieldRequired", t, func() {
		Convey("test case 1", func() {
			required := NewErrFieldRequired("123")
			So(required.field, ShouldEqual, "123")
		})
	})
}

//NewErrFieldMaxSize
func TestNewErrFieldMaxSize(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewErrFieldMaxSize", t, func() {
		Convey("test case 1", func() {
			size := NewErrFieldMaxSize("abc", 1)
			So(size, ShouldNotBeNil)
		})
	})
}

//NewErrFieldMaxLen
func TestNewErrFieldMaxLen(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewErrFieldMaxLen", t, func() {
		Convey("test case 1", func() {
			size := NewErrFieldMaxLen("abc", 1)
			So(size, ShouldNotBeNil)
		})
	})
}

//NewErrFieldValueExisted
func TestNewErrFieldValueExisted(t *testing.T) {
	//go InitConfig()
	Convey("Test_NewErrFieldValueExisted", t, func() {
		Convey("test case 1", func() {
			size := NewErrFieldValueExisted("abc")
			So(size, ShouldNotBeNil)
		})
	})
}
