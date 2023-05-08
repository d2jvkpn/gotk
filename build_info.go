package gotk

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	// "time"
)

func BuildInfo() (info map[string]any) {
	buildInfo, _ := debug.ReadBuildInfo()

	info = make(map[string]any, 8)
	info["goVersion"] = buildInfo.GoVersion

	parseFlags := func(str string) {
		for _, v := range strings.Fields(str) {
			k, v, _ := strings.Cut(v, "=")
			if strings.HasPrefix(k, "main.") && v != "" {
				info[k[5:]] = v
			}
		}
	}

	for _, v := range buildInfo.Settings {
		if v.Key == "-ldflags" || v.Key == "--ldflags" {
			parseFlags(v.Value)
		}
	}

	info["os"] = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	return info
}

func BuildInfoText(info map[string]any) string {
	strs := make([]string, 0, len(info))
	for k, v := range info {
		strs = append(strs, fmt.Sprintf("%s: %v", strings.Title(k), v))
	}

	sort.Strings(strs)
	return strings.Join(strs, "\n")
}
