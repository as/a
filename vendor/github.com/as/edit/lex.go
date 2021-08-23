package edit

// Put
import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	kind  Kind
	value string
}

func (i item) String() string {
	return fmt.Sprintf("%v %s", i.kind, i.value)
}

const (
	ralpha  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ="
	rcmd    = ralpha + "<>|"
	rdigit  = "0123456789"
	rop     = "+-;,"
	rmod    = "^$#/?"
	rescape = `#/?+-;,\abnrtx`
)

const maxBytes = ^uint64(0)

func max() string {
	return fmt.Sprintf("%v", maxBytes)
}

type Kind int

const (
	kindOp Kind = iota
	kindString
	kindSlash
	kindQuest
	kindRel
	kindComma
	kindDot
	kindEof
	kindColon
	kindSemi
	kindHash
	kindErr
	kindGlobal
	kindRegexp
	kindRegexpBack
	kindByteOffset
	kindLineOffset
	kindCount
	kindCmd
	kindArg
)
const (
	eof    = '\x00'
	slash  = '/'
	quest  = '?'
	comma  = ','
	plus   = "+"
	dot    = '.'
	colon  = ':'
	semi   = ';'
	hash   = '#'
	dollar = '$'
	caret  = '^'
)

type statefn func(*lexer) statefn

type lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	items  chan item
	lastop item
	first  bool
	esc    bool
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:   name,
		input:  input,
		items:  make(chan item),
		lastop: item{kindOp, "+"},
		first:  true,
	}
	go l.run() // run state machine
	return l, l.items
}

func (l *lexer) run() {
	for state := lexAny; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) acceptUntil(delim string) {
	lim := 8192
	i := 0
	for !strings.ContainsRune(delim, l.next()) {
		i++
		if i > lim {
			l.errorf("missing terminating char %q: %q\n", delim, l)
			l.ignore()
			l.emit(kindEof)
			return
		}
	}
	l.backup()
}

func (l *lexer) acceptEOF() {
	lim := 8192
	i := 0
	for l.next() != eof {
		i++
		if i > lim {
			l.emit(kindEof)
			return
		}
	}
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t Kind) {
	s, err := strconv.Unquote(`"` + l.String() + `"`)
	if err != nil {
		l.errorf(err.Error())
	}
	l.items <- item{t, s}
	l.start = l.pos
}

func (l *lexer) inject(it item) {
	l.items <- it
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) String() string {
	return string(l.input[l.start:l.pos])
}

func space(r rune) bool {
	return unicode.IsSpace(r)
}

func ignoreSpaces(l *lexer) {
	if l.accept(" 	") {
		l.acceptRun(" 	")
		l.ignore()
	}
}

func lexAny(l *lexer) statefn {
	ignoreSpaces(l)
	if l.accept(rdigit + rop + rmod) {
		l.backup()
		return lexAddr
	}
	l.emit(kindDot)
	return lexCmd
}

func lexAddr(l *lexer) statefn {
	ignoreSpaces(l)
	switch l.peek() {
	case eof:
		l.next()
		l.emit(kindEof)
		return nil
	case ',', ';':
		// LHS is empty so use #0
		if l.first {
			l.inject(item{kindByteOffset, "0"})
			l.first = false
		}
		return lexOp
	case '+', '-':
		if l.first {
			l.first = false
		}
		l.accept("+-")
		l.emit(kindRel)
		return lexAddr
	case slash, quest:
		return lexRegexp
	case dot:
		l.accept(".")
		l.emit(kindDot)
		return lexOp
	case hash:
		l.accept("#")
		l.ignore()
		ignoreSpaces(l)
		if !l.accept(rdigit) {
			return l.errorf("non-numeric offset")
		}
		l.acceptRun(rdigit)
		l.emit(kindByteOffset)
		return lexOp
	default:
		if l.accept("$") {
			l.ignore()
			l.inject(item{kindByteOffset, max()})
			return lexCmd
		}
		if l.accept("^") {
			l.ignore()
			l.inject(item{kindByteOffset, "0"})
			return lexOp
		}
		if l.accept(rdigit) {
			l.acceptRun(rdigit)
			l.emit(kindLineOffset)
			return lexOp
		}
	}
	return lexCmd
}

func lexArgsTuple(l *lexer) statefn {
	if !l.accept("s") {
		return l.errorf("want 's', have %q", l.String())
	}
	l.emit(kindCmd)

	if l.accept(rdigit) {
		// optional repetition count
		l.acceptRun(rdigit)
		l.emit(kindCount)
	}
	r := string(l.next())
	l.ignore()

	l.acceptUntil(r)
	l.emit(kindArg)
	if !l.accept(r) {
		return l.errorf("bad opening delimiter: %q", r)
	}
	l.ignore()
	l.acceptUntil(r)
	l.emit(kindArg)
	if !l.accept(r) {
		return l.errorf("bad closing delimiter: %q", r)
	}
	l.ignore()
	ignoreSpaces(l)
	if l.accept("g") {
		l.emit(kindGlobal)
	}
	return lexCmd
}

func lexCmd(l *lexer) statefn {
	ignoreSpaces(l)
	if l.peek() == eof {
		l.emit(kindEof)
		return nil
	}
	if l.peek() == 's' {
		return lexArgsTuple
	}
	if !l.accept(ralpha) {
		if l.accept("|<>") {
			l.emit(kindCmd)
			return lexArg2
		}
		return l.errorf("bad command")
	}
	l.emit(kindCmd)
	switch l.peek() {
	case eof:
		l.emit(kindEof)
		return nil
	default:
		return lexArg
	}
}

func lexOp(l *lexer) statefn {
	ignoreSpaces(l)
	tok := l.peek()
	if tok == eof {
		l.emit(kindEof)
		return nil
	}
	if tok == dollar {
		l.ignore()
		l.inject(item{kindByteOffset, max()})
		return lexAddr
	}
	op := ""
	if l.accept(rop) {
		op = l.String()
		l.emit(kindOp)
	}
	ignoreSpaces(l)
	if l.accept(rdigit + rmod) {
		if op == "" {
			l.inject(l.lastop)
		}
		l.backup()
		l.lastop = item{kindOp, op}
	}
	// use rcmd to det. whether closing addr is injected
	if tok := l.peek(); op != "" && (l.accept(rcmd) || tok == eof) {
		l.inject(item{kindByteOffset, max()})
		l.backup()
	}
	return lexAddr
}

func lexArg(l *lexer) statefn {
	r := string(l.next())
	l.ignore()
	l.acceptUntil(r)
	l.emit(kindArg)
	if !l.accept(string(r)) {
		return l.errorf("bad delimiter")
	}
	l.ignore()
	return lexCmd
}

func lexArg2(l *lexer) statefn {
	l.acceptEOF()
	l.emit(kindArg)
	return lexCmd
}

func lexRegexp(l *lexer) statefn {
	r := l.next()
	if r != '?' && r != '/' {
		return l.errorf("bad regexp delimiter: %q", l)
	}
	l.ignore()
	l.acceptUntil(string(rune(r)))
	if r == '?' {
		l.emit(kindRegexpBack)
	} else {
		l.emit(kindRegexp)
	}
	if !l.accept(string(rune(r))) {
		return l.errorf("bad regexp terminator: %q", l)
	}
	l.ignore()
	return lexOp
}

func (l *lexer) errorf(format string, args ...interface{}) statefn {
	l.items <- item{
		kindErr,
		fmt.Sprintf(format, args...),
	}
	return nil
}
