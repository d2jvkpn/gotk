package gotk

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	// "time"
)

func BuildInfo(prefixList ...string) (info map[string]any) {
	var (
		length    int
		prefix    string
		buildInfo *debug.BuildInfo
	)

	if prefix = "main."; len(prefixList) > 0 {
		prefix = prefixList[0]
	}
	length = len(prefix)

	buildInfo, _ = debug.ReadBuildInfo()

	info = make(map[string]any, 8)
	info["go_version"] = buildInfo.GoVersion
	info["os"] = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	parseFlags := func(str string) {
		for _, v := range strings.Fields(str) {
			k, v, _ := strings.Cut(v, "=")
			if strings.HasPrefix(k, prefix) && v != "" {
				info[k[length:]] = v
			}
		}
	}

	for _, v := range buildInfo.Settings {
		if v.Key == "-ldflags" || v.Key == "--ldflags" {
			parseFlags(v.Value)
		}
	}

	return info
}

func BuildInfoText(info map[string]any, pList ...string) string {
	var (
		strs []string
		p    string
	)

	if len(pList) > 0 {
		p = pList[0]
	}

	strs = make([]string, 0, len(info))
	for k, v := range info {
		// strs = append(strs, fmt.Sprintf("%s: %v", strings.Title(k), v))
		strs = append(strs, fmt.Sprintf("%s%s: %v", p, k, v))
	}

	sort.Strings(strs)
	return strings.Join(strs, "\n")
}
