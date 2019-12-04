//go:generate go test -v -gcflags=-l
package mockfn

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/yougg/assert"
)

func yes() bool { return true }
func no() bool  { return false }

func originNow() time.Time {
	fmt.Println("nothing")
	return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func TestReplace(t *testing.T) {
	Replace(no, yes)
	defer RevertAll()
	runtime.GC()
	assert.New(t).True(no())
}

func TestReplace1(t *testing.T) {
	a := assert.New(t)
	a.False(no())
	Replace(no, yes)
	a.True(no())
	a.True(Revert(no))
	a.False(no())
	a.False(Revert(no))
}

func TestReplace2(t *testing.T) {
	Replace(fmt.Println, func(a ...interface{}) (n int, err error) {
		s := make([]interface{}, len(a))
		for i, v := range a {
			s[i] = strings.Replace(fmt.Sprint(v), "ha", "*bleep*", -1)
		}
		return fmt.Fprintln(os.Stdout, s...)
	})
	n, err := fmt.Println("ha ha?")
	a := assert.New(t)
	a.Nil(err)
	a.Equal(17, n)
}

func TestReplaceEx(t *testing.T) {
	before := time.Now()
	ReplaceEx(time.Now, originNow, func() time.Time {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	during := time.Now()
	a := assert.New(t)
	a.True(Revert(time.Now))
	after := time.Now()

	a.Equal(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), during)
	a.NotEqual(before, during)
	a.NotEqual(during, after)
}

func TestReplaceEx1(t *testing.T) {
	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	ReplaceEx(time.Now, originNow, func() time.Time {
		return now
	})
	assert.New(t).Equal(time.Now(), now)
	Revert(time.Now)
}

func TestGuard(t *testing.T) {
	var guard *FuncGuard
	guard = Replace(no, func() bool {
		guard.Revert()
		defer guard.Restore()
		return !no()
	})
	a := assert.New(t)
	for i := 0; i < 100; i++ {
		a.True(no())
	}
	Revert(no)
}

func TestRevertAll(t *testing.T) {
	a := assert.New(t)
	a.False(no())
	Replace(no, yes)
	a.True(no())
	RevertAll()
	a.False(no())
}

type s struct{}

func (s *s) yes() bool { return true }

func TestReplaceRevert(t *testing.T) {
	i := &s{}

	a := assert.New(t)
	a.False(no())
	Replace(no, i.yes)
	a.True(no())
	Revert(no)
	a.False(no())
}

type f struct{}

func (f *f) No() bool  { return false }
func (f *f) Yes() bool { return true }

func TestReplaceInstanceMethod(t *testing.T) {
	i := &f{}
	a := assert.New(t)
	a.False(i.No())
	ReplaceInstanceMethod(reflect.TypeOf(i), "No", func(_ *f) bool { return true })
	a.True(i.No())
	a.True(RevertInstanceMethod(reflect.TypeOf(i), "No"))
	a.False(i.No())

	a.True(i.Yes())
	a.False(i.No())
	ReplaceInstanceMethod(reflect.TypeOf(i), "No", func(*f) bool { return true })
	a.True(i.Yes())
	a.True(i.No())
	a.True(RevertInstanceMethod(reflect.TypeOf(i), "No"))
	a.True(i.Yes())
	a.False(i.No())
	ReplaceInstanceMethod(reflect.TypeOf(i), "Yes", func(*f) bool { return false })
	a.False(i.Yes())
	a.False(i.No())
	a.True(RevertInstanceMethod(reflect.TypeOf(i), "Yes"))
	a.True(i.Yes())
	a.False(i.No())
	ReplaceInstanceMethod(reflect.TypeOf(i), "Yes", func(*f) bool { return false })
	ReplaceInstanceMethod(reflect.TypeOf(i), "No", func(*f) bool { return true })
	a.False(i.Yes())
	a.True(i.No())
	a.True(RevertInstanceMethod(reflect.TypeOf(i), "Yes"))
	a.True(RevertInstanceMethod(reflect.TypeOf(i), "No"))
	a.True(i.Yes())
	a.False(i.No())
}

func TestReplaceInstanceMethod1(t *testing.T) {
	var d *net.Dialer
	ReplaceInstanceMethod(reflect.TypeOf(d), "Dial", func(_ *net.Dialer, _, _ string) (net.Conn, error) {
		return nil, fmt.Errorf("no dialing allowed")
	})
	_, err := http.Get("http://google.com")
	fmt.Println(err) // Get http://google.com: no dialing allowed
}

func TestReplaceInstanceMethod2(t *testing.T) {
	var guard *FuncGuard
	guard = ReplaceInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
		guard.Revert()
		defer guard.Restore()

		if !strings.HasPrefix(url, "https://") {
			return nil, fmt.Errorf("only https requests allowed")
		}

		return c.Get(url)
	})

	_, err := http.Get("http://google.com")
	fmt.Println(err) // only https requests allowed
	resp, err := http.Get("https://google.com")
	fmt.Println(resp.Status, err) // 200 OK <nil>
}
