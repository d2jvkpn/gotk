package gotk

import (
	"fmt"
	"strconv"
	"strings"
)

// parse ports from string: 8000,8001,8002-8009
func ParsePorts(str string) (ports []uint64, err error) {
	var (
		p1, p2 uint64
		strs   []string
	)

	strs = strings.Split(str, ",")
	ports = make([]uint64, 0, len(strs))

	for _, v := range strs {
		list := strings.Split(strings.TrimSpace(v), "-")
		if len(list) == 0 {
			continue
		}

		if p1, err = strconv.ParseUint(strings.TrimSpace(list[0]), 10, 64); err != nil {
			return nil, err
		}

		if len(list) == 1 {
			p2 = p1
		} else {
			if p2, err = strconv.ParseUint(strings.TrimSpace(list[1]), 10, 64); err != nil {
				return nil, err
			}
		}

		for p := p1; p <= p2; p++ {
			ports = append(ports, p)
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports")
	}

	return ports, nil
}
