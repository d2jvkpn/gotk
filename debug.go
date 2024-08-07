package gotk

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
)

var (
	_StackRE *regexp.Regexp = regexp.MustCompile("\n.*\n\t.*")
)

func _fn2() {
	defer func() {
		var intf any

		if intf = recover(); intf == nil {
			return
		}
		fmt.Println("!!!", intf)
		fmt.Println(">>>", Stack(""))
	}()

	_fn1()
}

func _fn1() {
	var mySlice []int
	j := mySlice[0]

	fmt.Printf("Hello, playground %d", j)
}

func Stack(prefix string) (slice []string) {
	bts := bytes.TrimSpace(debug.Stack())
	// fmt.Printf(">>>\n%s\n<<<\n", bts)
	out := _StackRE.FindAllStringSubmatch(string(bts), -1)
	const skip = 2
	if len(out) < skip {
		return make([]string, 0)
	}
	slice = make([]string, 0, len(out)-skip)

	for i := skip; i < len(out); i++ {
		t := strings.Split(strings.TrimSpace(out[i][0]), "\n\t")
		if len(t) <= 1 {
			continue
		}
		if prefix != "" && !strings.HasPrefix(t[0], prefix) {
			continue
		}

		fp := strings.Fields(t[1])[0]
		fn := filepath.Base(strings.Split(t[0], "(")[0])
		// gotk/debug_test.go:8(gotk.TestPanic)
		slice = append(slice, fmt.Sprintf("%s(%s)", fp, fn))
	}

	return
}
