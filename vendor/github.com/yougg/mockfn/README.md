# Go mock function

## Using mockfn

Call `mockfn.Replace(<target function>, <replacement function>)` to replace a function.   
For example:

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yougg/mockfn"
)

func main() {
	mockfn.Replace(fmt.Println, func(a ...interface{}) (n int, err error) {
		s := make([]interface{}, len(a))
		for i, v := range a {
			s[i] = strings.Replace(fmt.Sprint(v), "ha", "*bleep*", -1)
		}
		return fmt.Fprintln(os.Stdout, s...)
	})
	fmt.Println("ha ha?")
}
```

call `mockfn.Revert(<target function>)` to revert the method again.

replace an instance method with `mockfn.ReplaceInstanceMethod(<type>, <name>, <replacement>)`.   
You get the type by using `reflect.TypeOf`, and your replacement function simply takes the instance as the first argument.   
To disable all network connections, you can do as follows for example:

```go
package main

import (
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/yougg/mockfn"
)

func main() {
	var d *net.Dialer // Has to be a pointer to because `Dial` has a pointer receiver
	mockfn.ReplaceInstanceMethod(reflect.TypeOf(d), "Dial", func(_ *net.Dialer, _, _ string) (net.Conn, error) {
		return nil, fmt.Errorf("no dialing allowed")
	})
	_, err := http.Get("http://google.com")
	fmt.Println(err) // Get http://google.com: no dialing allowed
}
```

Note that mocking the method for just one instance is currently not possible, `ReplaceInstanceMethod` will mock it for all instances. 
Don't bother trying `mockfn.Replace(instance.Method, replacement)`, it won't work. 
`mockfn.RevertInstanceMethod(<type>, <name>)` will undo `ReplaceInstanceMethod`.

If you want to remove all currently applied mock functions simply call `mockfn.RevertAll`. This could be useful in a test teardown function.

If you want to call the original function from within the replacement you need to use a `mockfn.FuncGuard`.   
A funcguard allows you to easily remove and restore the function so you can call the original function.   
For example:

```go
package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/yougg/mockfn"
)

func main() {
	var guard *mockfn.FuncGuard
	guard = mockfn.ReplaceInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
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
```

## Notes

1. mockfn sometimes fails to mock a function if inlining is enabled. Try running your tests with inlining disabled, for example: `go test -gcflags=-l`. The same command line argument can also be used for build.
2. mockfn won't work on some security-oriented operating system that don't allow memory pages to be both write and execute at the same time. With the current approach there's not really a reliable fix for this.
3. mockfn is not thread safe, Or any kind of safe.
4. mockfn is not recommend use it outside of a testing environment.