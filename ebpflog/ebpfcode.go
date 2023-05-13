package ebpflog

import (
	"fmt"
	"strconv"
	"strings"
)

func (l *logContext) processLineEbpf(in string, lite bool) string {
	ln := l.lnRegEx.FindString(in)
	ln = ln[0 : len(ln)-1]
	lnInt, err := strconv.Atoi(ln)
	if err != nil {
		panic(err)
	}

	lnStr := strconv.Itoa(lnInt)
	line := strings.TrimPrefix(in, lnStr+": ")
	if line == "safe" {
		return "<span class=\"bpf\"><b>" + in + "</b>(Visited state ✔)</span>"
	}
	if len(line) < 5 {
		panic("ebpf code:" + line)
	}
	add := ""
	if !lite {
		add = fmt.Sprintf("id=\"pc%d\"", lnInt)
	}
	line = line[5:]
	line = l.handleEbpfCodeLine(lnInt, line)
	return fmt.Sprintf("<span %s class=\"bpf\"><b>%s</b>: %s</span>", add, lnStr, line)
}
