package ebpflog

import (
	"io"
	"regexp"
)

type LogContext interface {
	Process() error
}

type logContext struct {
	in              io.Reader
	out             io.Writer
	lnRegEx         *regexp.Regexp
	regRegEx        *regexp.Regexp
	goRegEx         *regexp.Regexp
	inBacktrace     bool
	stateCode       bool
	stateAnnotation bool
	insideBlock     bool
}

func NewLogContext(in io.Reader, out io.Writer) LogContext {
	l := &logContext{
		in:  in,
		out: out,
	}
	l.lnRegEx = regexp.MustCompile(`^(\d+):`)
	l.goRegEx = regexp.MustCompile(`goto pc(\D\d+)`)
	l.regRegEx = regexp.MustCompile(`(?m)[Rr][0-9]{1,2}`)
	return l
}
