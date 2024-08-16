package ebpflog

import (
	"fmt"
	"strconv"
)

func (l *logContext) processLineEbpf(in string, lite bool) string {
	lineParts := l.lnRegEx.FindStringSubmatch(in)
	lnInt, err := strconv.Atoi(lineParts[1])
	if err != nil {
		panic(err)
	}

	lnStr := strconv.Itoa(lnInt)
	line := lineParts[3]
	if line == "safe" {
		return "<span class=\"bpf\"><b>" + in + "</b>(Visited state ✔)</span>"
	}
	if line == "exit" {
		return "<b>exit</b>"
	}
	if len(line) < 5 {
		panic("ebpf code:" + line)
	}
	add := ""
	if !lite {
		add = fmt.Sprintf("id=\"pc%d\"", lnInt)
	}
	line = l.handleEbpfCodeLine(lnInt, line)
	return fmt.Sprintf("<span %s class=\"bpf\"><b>%s</b>: %s</span>", add, lnStr, line)
}
