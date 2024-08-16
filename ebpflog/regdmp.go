package ebpflog

import (
	"fmt"
	"io"
	"strings"
)

func getStateMap(in string) map[string]string {
	state := struct {
		isKey bool
		depth int
		key   string
		val   string
		kv    map[string]string
	}{
		kv: make(map[string]string),
	}
	state.isKey = true
	for _, ch := range in {
		switch ch {
		case ' ':
			if state.depth == 0 && !state.isKey {
				state.isKey = true
				state.kv[state.key] = state.val
				state.key = ""
				state.val = ""
				continue
			}
		case '=':
			if state.depth == 0 && state.isKey {
				state.isKey = false
				state.key = state.val
				state.val = ""
				continue
			}
		case '(':
			state.depth++
		case ')':
			state.depth--
		}

		state.val += string(ch)
	}
	return state.kv
}

type registerInfo struct {
	init     bool
	liveness struct {
		read  bool
		write bool
		done  bool
	}
	value string
}

const (
	StackInvalid = iota //'?',
	StackSpill          //'r',
	StackMisc           //'m',
	StackZero           //'0',
	StackDynptr         //'d',
	StackIter           //'i',
	StackMapValue
	StackCtx
	StackImm
)

type stateDump struct {
	regs  [11]registerInfo
	stack [512 / 8]registerInfo
}

func newRegisterDump(in string) stateDump {
	sm := getStateMap(in)
	sd := stateDump{}
	for k, v := range sm {
		if strings.HasPrefix(k, "R") {
			idx := 0
			flg := ""
			cnt, _ := fmt.Sscanf(k, "R%d_%s", &idx, &flg)
			if cnt == 0 {
				panic(k)
			}
			sd.regs[idx].init = true
			sd.regs[idx].value = v
			if strings.ContainsRune(flg, 'r') {
				sd.regs[idx].liveness.read = true
			}
			if strings.ContainsRune(flg, 'w') {
				sd.regs[idx].liveness.write = true
			}
			if strings.ContainsRune(flg, 'D') {
				sd.regs[idx].liveness.done = true
			}
		} else if strings.HasPrefix(k, "fp") {
			off := 0
			flg := ""
			cnt, _ := fmt.Sscanf(k, "fp-%d_%s", &off, &flg)
			if cnt == 0 {
				panic(k)
			}
			off -= 8
			idx := off / 8
			sd.stack[idx].init = true
			sd.stack[idx].value = v
			if strings.ContainsRune(flg, 'r') {
				sd.stack[idx].liveness.read = true
			}
			if strings.ContainsRune(flg, 'w') {
				sd.stack[idx].liveness.write = true
			}
			if strings.ContainsRune(flg, 'D') {
				sd.stack[idx].liveness.done = true
			}
		} else {
			panic(k)
		}
	}
	return sd
}

var gTblId = 0

func (l *logContext) processRegisterDump(in string, w io.Writer, newVariant bool) {
	if newVariant {
		parts := l.newRegDmpRegEx.FindStringSubmatch(in)
		in = parts[2]
	}
	sd := newRegisterDump(in)
	gTblId++
	fmt.Fprintf(w, "<div class=\"regdump\">\n")
	fmt.Fprintf(w, "<a class=\"button\" onclick=\"toggleshow(this)\">%s</a>\n", getIcon("👁"))
	fmt.Fprintf(w, "<div class=\"regdumpcontent\">\n")

	fmt.Fprintf(w, "<div class=\"registers\">Registers\n")

	fmt.Fprintf(w, "<table cellspacing=0 class=\"registerst\">\n")

	fmt.Fprintf(w, "<tr><th>Reg</th><th>Liveness</th><th>Value</th></tr>\n")

	for idx, r := range sd.regs {
		liv := ""
		if r.liveness.read {
			liv += "+read "
		}
		if r.liveness.write {
			liv += "+write "
		}
		if r.liveness.done {
			liv += "+done "
		}
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", l.regname(idx), liv, r.value)
	}

	fmt.Fprintf(w, "</table><!-- registerst -->\n") // registerst

	fmt.Fprintf(w, "</div><!-- registers -->\n") // registers

	// stack

	fmt.Fprintf(w, "<div class=\"registers\">Stack dump\n") // stack

	fmt.Fprintf(w, "<table cellspacing=0 class=\"registerst\">\n")

	fmt.Fprintf(w, "<tr><th>Offset</th><th>Liveness</th><th>Value</th></tr>\n")

	for idx, r := range sd.stack {
		if !r.init {
			continue
		}
		liv := ""
		if r.liveness.read {
			liv += "+read "
		}
		if r.liveness.write {
			liv += "+write "
		}
		if r.liveness.done {
			liv += "+done "
		}
		stk := fmt.Sprintf("fp-%d", 8*(idx+1))
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", stk, liv, r.value)
	}

	fmt.Fprintf(w, "</table><!-- registerst -->\n") // registerst

	fmt.Fprintf(w, "</div><!-- registerst -->\n") // registers

	fmt.Fprintf(w, "</div><!-- regdumpcontent -->\n") // regdumpcontent
	fmt.Fprintf(w, "</div><!-- regdump -->\n")        // regdump
}
