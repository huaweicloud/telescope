package collectors

import (
	"testing"
)

func TestGetMountPrefix(t *testing.T) {

	mountPointTable := []struct {
		input    string
		expected string
	}{
		{"C:", "C_"},
		{"/", "SlAsH_"},
		{"/var/log/ces", "SlAsHvarSlAsHlogSlAsHces_"},
	}

	for _, value := range mountPointTable {
		if getMountPrefix(value.input) != value.expected {
			t.Errorf("GetMountPrefix Error, input is %s, correct expected is %s\n", value.input, getMountPrefix(value.input))
		} else {
			t.Logf("Success, input is %s, expected is %s\n", value.input, value.expected)
		}
	}

}
