package edit

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Print string

var eprint = n

type Emitted struct {
	Name string
	Dot  []Dot
}

type parser struct {
	cmd       []*Command
	last, tok item
	in        chan item
	out       chan func()
	err       error
	stop      chan error
	addr      Address

	recache map[string]*regexp.Regexp

	Emit    *Emitted
	Options *Options
}

func parse(i chan item, opts ...*Options) *parser {
	var o *Options
	if len(opts) != 0 {
		o = opts[0]
	}
	p := &parser{
		in:      i,
		stop:    make(chan error),
		Emit:    &Emitted{},
		Options: o,
		recache: make(map[string]*regexp.Regexp),
	}
	go p.run()
	return p
}

func (p *parser) compileRegexp(s string) (re *regexp.Regexp, err error) {
	re, ok := p.recache[s]
	if !ok {
		re, err = regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		p.recache[s] = re
	}
	return re, nil
}

func parseAddr(p *parser) (a Address) {
	a0 := parseSimpleAddr(p)
	p.Next()
	op, a1 := parseOp(p)
	if op == '\x00' {
		return a0
	}
	p.Next()
	return &Compound{a0: a0, a1: a1, op: op}
}

func parseOp(p *parser) (op byte, a Address) {
	//	fmt.Printf("parseOp:1 %s\n", p.tok)
	if p.tok.kind != kindOp {
		return
	}
	v := p.tok.value
	if v == "" {
		eprint("no value" + v)
		return
	}
	if !strings.ContainsAny(v, "+-;,") {
		//		eprint(fmt.Sprintf("bad op: %q", v))
	}
	p.Next()
	return v[0], parseSimpleAddr(p)
}

func tryRelative(p *parser) int {
	v := p.tok.value
	k := p.tok
	if k.kind == kindRel {
		defer p.Next()
		if v == "+" {
			return 1
		}
		return -1
	}
	return 0
}

// Put
func parseSimpleAddr(p *parser) (a Address) {
	//fmt.Printf("parseSimpleAddr:1 %s\n", p.tok)
	back := false
	rel := tryRelative(p)
	v := p.tok.value
	k := p.tok
	//fmt.Printf("%s\n", k)
	switch k.kind {
	case kindRegexpBack:
		back = true
		fallthrough
	case kindRegexp:
		re, err := regexp.Compile(v)
		if err != nil {
			p.fatal(err)
			return
		}
		if rel != -1 {
			rel = 1
		}
		return &Regexp{re, back, 1}
	case kindLineOffset, kindByteOffset:
		i := p.mustatoi(v)
		if rel < 0 {
			i = -i
		}
		if k.kind == kindLineOffset {
			return &Line{i, rel}
		}
		return &Byte{i, rel}
	case kindDot:
		return &Dot{}
	}
	p.err = fmt.Errorf("bad address: %q", v)
	return
}

func parseArg(p *parser) (arg string) {
	p.Next()
	if p.tok.kind != kindArg {
		p.fatal(fmt.Errorf("want arg, have %q", p.tok.value))
	}
	return p.tok.value
}

func (p *parser) Dot(f Editor) (q0, q1 int64) {
	return f.Dot()
}

type Sender interface {
	Send(e interface{})
	SendFirst(e interface{})
}

// Put
func parseCmd(p *parser) (c *Command) {
	v := p.tok.value
	c = &Command{}
	c.s = v
	switch v {
	case "h":
		parseArg(p)
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			p.Emit.Dot = append(p.Emit.Dot, Dot{q0, q1})
		}
		return
	case "=":
		if p.Options == nil || p.Options.Sender == nil {
			return
		}
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			str := fmt.Sprintf("%s:#%d,#%d", p.Options.Origin, q0+1, q1)
			p.Options.Sender.Send(Print(str))
		}
		return
	case "p":
		if p.Options == nil || p.Options.Sender == nil {
			return
		}
		c.fn = func(f Editor) {
			q0, q1 := p.Dot(f)
			str := fmt.Sprintf("%s", f.Bytes()[q0:q1])
			p.Options.Sender.Send(Print(str))
		}
		return
	case "a":
		c.fn = Append{Data: []byte(parseArg(p))}.Apply
		return
	case "i":
		c.fn = Insert{Data: []byte(parseArg(p))}.Apply
		return
	case "c":
		c.fn = Change{To: []byte(parseArg(p))}.Apply
		return
	case "d":
		c.fn = Delete{}.Apply
		return
	case "e":
	case "k":
	case "r":
		c.fn = ReadFile{Name: parseArg(p)}.Apply
		return
	case "s":

		matchn := int64(1)
		sre := parseArg(p)
		if p.tok.kind == kindCount {
			matchn = p.mustatoi(sre)
			sre = parseArg(p)
		}
		replamp := compileReplaceAmp(parseArg(p))

		// And at this point I realized that instead
		// of a one token look-ahead parser, I have
		// a one token look-behind parser. How unfortunate.
		//
		// TODO(as): check for 'g' here after fixing the parser
		// try parsing the last part of the construction anyway
		// and look for 'g'
		if g := parseArg(p); p.tok.kind == kindGlobal {
			if g != "g" {
				p.fatal("s: suffix not supported: " + g)
				return
			}
			matchn = -1
		}
		if sre == "" {
			eprint("s: no regexp to find")
			return
		}

		re, err := regexp.Compile(sre)
		if err != nil {
			p.fatal(err)
			return
		}
		c.fn = S{
			Regexp:     re,
			ReplaceAmp: replamp,
			Limit:      matchn,
		}.Apply
		return
	case "w":
		c.fn = WriteFile{Name: parseArg(p)}.Apply
		return
	case "m":
		a1 := parseSimpleAddr(p)
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			p := append([]byte{}, f.Bytes()[q0:q1]...)
			a1.Set(f)
			_, a1 := f.Dot()
			f.Delete(q0, q1)
			f.Insert(p, a1)
		}
		return
	case "t":
		a1 := parseSimpleAddr(p)
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			p := f.Bytes()[q0:q1]
			a1.Set(f)
			_, a1 := f.Dot()
			f.Insert(p, a1)
		}
		return
	case "g":
		re, err := regexp.Compile(parseArg(p))
		if err != nil {
			p.fatal(err)
			return
		}
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			if re.Match(f.Bytes()[q0:q1]) {
				if nextfn := c.nextFn(); nextfn != nil {
					nextfn(f)
				}
			}
		}
		return
	case "v":
		re, err := regexp.Compile(parseArg(p))
		if err != nil {
			p.fatal(err)
			return
		}
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			if !re.Match(f.Bytes()[q0:q1]) {
				if nextfn := c.nextFn(); nextfn != nil {
					nextfn(f)
				}
			}
		}
		return
	case "|":
		c.fn = Pipe{To: parseArg(p)}.Apply
		return
	case ">":
		filename := parseArg(p)
		c.fn = func(f Editor) {
			fd, err := os.Create(filename)
			if err != nil {
				eprint(err)
				return
			}
			defer fd.Close()
			q0, q1 := f.Dot()
			_, err = io.Copy(fd, bytes.NewReader(f.Bytes()[q0:q1]))
			if err != nil {
				eprint(err)
			}
		}
		return
	case "x":
		re, err := regexp.Compile(parseArg(p))
		if err != nil {
			p.fatal(err)
			return
		}
		buf := new(bytes.Reader)
		c.fn = func(f Editor) {

			sp, ep := f.Dot()
			buf.Reset(f.Bytes()[sp:ep])
			q0 := int64(0)
			for {
				loc := re.FindReaderIndex(buf)
				if loc == nil {
					break
				}
				q1 := q0 + int64(loc[1])
				q0 += int64(loc[0])
				//				log.Printf("match: %q location (%d,%d)", f.Bytes()[sp+q0:sp+q1], sp+q0, sp+q1)
				f.Select(sp+q0, sp+q1)
				if nextfn := c.nextFn(); nextfn != nil {
					nextfn(f)
				}
				q0 = q1
				buf.Seek(q0, 0)
				if sp+q0 == ep {
					break
				}
			}
			f.Select(ep, ep)
		}
		return
	case "y":
		re, err := regexp.Compile(parseArg(p))
		if err != nil {
			p.fatal(err)
			return
		}
		c.fn = func(f Editor) {
			q0, q1 := f.Dot()
			x0, x1 := int64(0), int64(0)
			y0, y1 := int64(0), q1
			buf := bytes.NewReader(f.Bytes()[q0:q1])
			for {
				loc := re.FindReaderIndex(buf)
				if loc == nil {
					buf.Seek(x1, 0)
					eprint("not found")
					break
				}
				y0 = x1
				x0, x1 = int64(loc[0])+x1, int64(loc[1])+x1
				y1 = x0
				f.Select(q0+y0, q0+y1)
				if nextfn := c.nextFn(); nextfn != nil {
					nextfn(f)
				}
				buf.Seek(x1, 0)
			}
			if x1 != q1 {
				f.Select(q0+x1, q1)
				if nextfn := c.nextFn(); nextfn != nil {
					nextfn(f)
				}
			}
		}
		return
	}
	return nil
}

type ReplaceAmp []func([]byte) string

func (r ReplaceAmp) Run(ed Editor, q1 int64, sel []byte) (n int) {
	for _, fn := range r {
		b := []byte(fn(sel))
		n += len(b)
		ed.Insert(b, q1)
	}
	return n
}
func (r ReplaceAmp) Gen(replace []byte) (b []byte) {
	for _, fn := range r {
		b = append(b, fn(replace)...)
	}
	return b
}
func compileReplaceAmp(in string) (s ReplaceAmp) {
	// we can use strings.Map but then invalid runes
	// are replaced. we'll do it the old fashioned
	// way for now

	// strings are immutable, but this is faster than
	// the string 'builders' on platforms tested
	for {
		i := strings.Index(in, `&`)
		if i == -1 {
			s = append(s, func([]byte) string { return in })
			return
		}
		if i > 0 && in[i] == '\\' {
			s = append(s, func(b []byte) string { return in[:i-2] })
			s = append(s, func(b []byte) string { return "&" })
			i += 2
		} else {
			s = append(s, func(b []byte) string { return in[:i-1] })
			s = append(s, func(b []byte) string { return string(b) })
			i++
		}
		if i == len(in) {
			break
		}
		in = in[i:]
	}
	if s == nil {
		s = append(s, func(b []byte) string { return "" })
	}
	return
}
func (p *parser) Next() *item {
	p.last = p.tok
	p.tok = <-p.in
	return &p.tok
}

func (p *parser) run() {
	tok := p.Next()
	if tok.kind == kindEof || p.err != nil {
		if tok.kind == kindEof {
			p.fatal(fmt.Errorf("run: unexpected eof"))
			return
		}
		p.fatal(fmt.Errorf("run: %s", p.err))
		return
	}
	p.addr = parseAddr(p)
	for {
		c := parseCmd(p)
		if c == nil {
			break
		}
		p.cmd = append(p.cmd, c)
		p.Next()
	}
	p.stop <- p.err
	close(p.stop)
}

func (p *parser) mustatoi(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		p.fatal(err)
	}
	return int64(i)
}
func (p *parser) fatal(why interface{}) {
	switch why := why.(type) {
	default:
		//fmt.Println(why)
		_ = why
	}
}

func n(i ...interface{}) (n int, err error) {
	return
}
