package utils

import (
	"strings"
	"time"
)

type p struct {
	find   string
	format string
	reg    string
}

/*
	Formats:
	M    - month (1)
	MM   - month (01)
	MMM  - month (Jan)
	MMMM - month (January)
	D    - day (2)
	DD   - day (02)
	DDD  - day (Mon)
	DDDD - day (Monday)
  T    - Time  (T, 2006-01-02T15:04:05)
	YY   - year (06)
	YYYY - year (2006)
  hh   - hours (15)
	mm   - minutes (04)
	ss   - seconds (05)
	AM/PM hours: 'h' followed by optional 'mm' and 'ss' followed by 'pm', e.g.
  hpm        - hours (03PM)
  h:mmpm     - hours:minutes (03:04PM)
  h:mm:sspm  - hours:minutes:seconds (03:04:05PM)
  Time zones: a time format followed by 'ZZZZ', 'ZZZ' or 'ZZ', e.g.
  hh:mm:ss ZZZZ (16:05:06 +0100)
  hh:mm:ss ZZZ  (16:05:06 CET)
	hh:mm:ss ZZ   (16:05:06 +01:00)
*/

var Placeholder = []p{
	p{"hh", "15", "\\d{2}"},
	p{"h", "03", "\\d{2}"},
	p{"mm", "04", "\\d{2}"},
	p{"ss", "05", "\\d{2}"},
	p{"SSS", "999", "\\d{3}"},
	p{"MMMM", "January", "\\w{3,9}"},
	p{"MMM", "Jan", "\\w{3}"},
	p{"MM", "01", "\\d{2}"},
	p{"M", "1", "\\d{1}"},
	p{"T", "T", "\\w{1}"},
	p{"pm", "PM", "\\w{2}"},
	p{"ZZZZ", "-0700", "(-|\\+)\\d{4}"},
	p{"ZZZ", "MST", "\\w{3}"},
	p{"ZZ", "Z07:00", "\\w{1}"},
	p{"YYYY", "2006", "\\d{4}"},
	p{"YY", "06", "\\d{2}"},
	p{"DDDD", "Monday", "\\w{6,9}"},
	p{"DDD", "Mon", "\\w{3}"},
	p{"DD", "02", "\\d{2}"},
	p{"D", "2", "\\d{1}"},
}

var (
	DefaultTimeFormat     = "hh:mm:ss"
	DefaultDateFormat     = "YYYY-MM-DD"
	DefaultDateTimeFormat = "YYYY-MM-DDThh:mm:ss"
)

func replace(in string) (out string) {
	out = in
	for _, ph := range Placeholder {
		out = strings.Replace(out, ph.find, ph.format, -1)
	}
	return
}

func GetReg(in string) (reg string) {
	reg = in
	for _, ph := range Placeholder {
		reg = strings.Replace(reg, ph.find, ph.reg, -1)
	}
	return
}

// Format formats a date based on Microsoft Excel (TM) conventions
func Format(format string, date time.Time) string {
	return date.Format(replace(format))
}

// Parse parses a value to a date based on Microsoft Excel (TM) formats
func Parse(format string, value string) (time.Time, error) {
	return time.ParseInLocation(replace(format), value, time.Local)
}
