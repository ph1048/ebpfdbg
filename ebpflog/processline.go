package ebpflog

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"
)

func (l *logContext) processFromMessage(in string, w io.Writer) {
	in = strings.ReplaceAll(in, " (speculative execution)", "")
	from := 0
	to := 0
	cnt, _ := fmt.Sscanf(in, "from %d to %d: ", &from, &to)
	if cnt != 2 {
		panic(in)
	}
	parts := strings.SplitN(in, ": ", 2)
	rest := parts[1]
	fmt.Fprintf(w, "<span class=\"bbtrans\"><a class=\"pclnk\" href=\"#pc%d\">%d</a> ⇨ <b>%d</b></span><br>", from, from, to)
	if rest == "safe" {
		fmt.Fprintf(w, "Visited state ✔")
		return
	}
	l.processRegisterDump(rest, w, false)
}

func (l *logContext) processParentMessage(in string, w io.Writer) {
	didnt := false
	if strings.HasPrefix(in, "parent didn&#39;t have ") {
		didnt = true
		in = strings.TrimPrefix(in, "parent didn&#39;t have ")
	} else if strings.HasPrefix(in, "parent already had ") {
		in = strings.TrimPrefix(in, "parent already had ")
	} else {
		panic(in)
	}
	in = strings.TrimSuffix(in, " marks")
	in = l.processMasks(in)
	did := "already had"
	if didnt {
		did = "didn't have"
	}
	fmt.Fprintf(w, "%s<i>Parent %s: %s</i>", getIcon("ℹ"), did, in)
}

func (l *logContext) handleBacktraceStart(in string) string {
	li := 0
	fi := 0
	fmt.Sscanf(in, "last_idx %d first_idx %d", &li, &fi)
	return fmt.Sprintf("<a class=\"button\" onclick=\"toggleshow(this)\">%sBacktrace</a> <b>%d</b> ⇦ <a class=\"pclnk\" href=\"#pc%d\">%d</a>", getIcon("📃"), li, fi, fi)
}

func (l *logContext) processMasks(in string) string {
	regs := uint32(0)
	stack := uint64(0)
	fmt.Sscanf(in, "regs=%x stack=%x", &regs, &stack)
	out := ""
	if regs != 0 {
		first := true
		for i := 0; i < 32; i++ {
			if regs&(1<<i) != 0 {
				if first {
					first = false
				} else {
					out += ", "
				}
				out += l.regname(i)
			}
		}
	}
	if stack != 0 {
		first := true
		for i := 0; i < 64; i++ {
			if stack&(1<<i) != 0 {
				if first {
					first = false
				} else {
					out += ", "
				}
				out += "fp-" + strconv.Itoa((i+1)*8)
			}
		}
	}
	return out
}

func (l *logContext) processLineMasks(in string) string {
	parts := strings.Split(in, " before ")
	return "<tr><td>" + l.processMasks(parts[0]) + "</td><td>" + l.processLineEbpf(parts[1], true) + "</td></tr>"
}

func (l *logContext) handleEbpfCodeLine(pc int, in string) string {
	in = l.regstring(in)
	goexpr := l.goRegEx.FindString(in)
	if goexpr != "" {
		goexpr = strings.TrimPrefix(goexpr, "goto pc")
		off, err := strconv.Atoi(goexpr[1:])
		if err != nil {
			panic(err)
		}
		if goexpr[0] == '-' {
			off = 0 - off
		}
		tgt := pc + 1 + off
		return fmt.Sprintf("%s (target: <b><a class=\"pclnk\" href=\"#pc%d\">%d</a></b>)", in, tgt, tgt)
	}
	return in
}

func (l *logContext) regname(idx int) string {
	return fmt.Sprintf("<span class=\"reg_r%d\">R%d</span>", idx, idx)
}

func (l *logContext) regstr(s string) string {
	return fmt.Sprintf("<span class=\"reg_%s\">%s</span>", strings.ToLower(s), s)
}

func (l *logContext) regstring(in string) string {
	return l.regRegEx.ReplaceAllStringFunc(in, func(s string) string {
		return l.regstr(s)
	})
}

func (l *logContext) processLineSource(in string) string {
	return "<span class=\"src\">" + in + "</span>"
}

func (l *logContext) processLineErr(in string) string {
	return "<br><span class=\"err\">" + in + "</span>"
}

func (l *logContext) processLineRes(in string) string {
	return "<br><span class=\"result\">" + in + "</span>"
}

func (l *logContext) processLine(in string, skip *bool, w io.Writer) {
	in = strings.TrimSpace(in)
	prevStateCode := l.stateCode
	l.stateCode = false
	isBacktraceLine := false
	isStartingBlock := false
	isEndingBlock := false
	emptyline := false

	tmpBuf := []byte{}
	buf := bytes.NewBuffer(tmpBuf)

	if strings.HasPrefix(in, "regs=") {
		// handle backtrace line
		*skip = true
		fmt.Fprintln(buf, l.processLineMasks(in))
		isBacktraceLine = true
	} else if l.inBacktrace {
		// not backtrace line -- finish table
		l.inBacktrace = false
		fmt.Fprintln(w, "</table></div><!-- table for backtrace --><br>")
	}

	if isBacktraceLine {
	} else if strings.HasPrefix(in, "last_idx") {
		// this line is backtrace start
		bt := l.handleBacktraceStart(in)
		l.inBacktrace = true
		if !prevStateCode {
			fmt.Fprintf(buf, "<br>\n")
		}
		fmt.Fprintf(buf, "<div class=\"backtrace\">\n")
		fmt.Fprintln(buf, bt)
		fmt.Fprintln(buf, "<table class=\"bttbl regdumpcontent\">")
	} else if in == "" {
		*skip = true
		fmt.Fprintln(buf, "<div class=\"ending\"></div>")
		isEndingBlock = true
	} else if in == ";" {
		*skip = true
		emptyline = true
	} else if strings.HasPrefix(in, "; ") {
		fmt.Fprintln(buf, l.processLineSource(in))
		l.stateCode = true
	} else if l.lnRegEx.MatchString(in) {
		fmt.Fprintln(buf, l.processLineEbpf(in, false))
		l.stateCode = true
	} else if strings.HasPrefix(in, "R0") {
		l.processRegisterDump(in, buf, false)
	} else if strings.HasPrefix(in, "parent") {
		l.processParentMessage(in, buf)
	} else if strings.HasPrefix(in, "from ") {
		fmt.Fprintln(buf, "<div class=\"starting\"></div>")
		l.processFromMessage(in, buf)
		isStartingBlock = true
	} else if strings.HasPrefix(in, "processed ") {
		fmt.Fprintln(buf, l.processLineRes(in))
	} else if l.newRegDmpRegEx.MatchString(in) {
		l.processRegisterDump(in, buf, true)
	} else {
		fmt.Fprintln(buf, l.processLineErr(in))
	}

	if l.insideBlock && isEndingBlock {
		l.insideBlock = false
	}
	if !l.insideBlock && isStartingBlock {
		l.insideBlock = true
	}

	// handle annotations
	if l.insideBlock && !isStartingBlock && !l.stateAnnotation && !l.stateCode && !emptyline {
		l.stateAnnotation = true
		fmt.Fprintln(w, "<div class=\"annotation\">")
	} else if l.stateAnnotation && (l.stateCode || isEndingBlock) {
		l.stateAnnotation = false
		fmt.Fprintln(w, "</div><!-- annotation -->")
	}

	fmt.Fprintln(w, buf)
}

func (l *logContext) Process() error {

	l.out.Write([]byte(header))

	l.insideBlock = true

	scanner := bufio.NewScanner(l.in)
	for scanner.Scan() {
		line := html.EscapeString(scanner.Text())
		skip := false
		l.processLine(line, &skip, l.out)
		if skip {
			continue
		}
		fmt.Fprintln(l.out, "<br>")
	}

	l.out.Write([]byte(footer))

	return nil
}
