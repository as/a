// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package win32

import (
	"fmt"
	"syscall"
	"unicode/utf16"

	"github.com/as/shiny/event/key"
	"github.com/as/shiny/screen"
)

type Key = key.Event

var KeyEvent func(hwnd syscall.Handle, e key.Event)
var keyboardLayout = GetKeyboardLayout(0)

func changeLanguage(h syscall.Handle, m uint32, charset, localeID uintptr) {
}

func readRune(vKey uint32, scanCode uint8) rune {
	var (
		keystate [256]byte
		buf      [4]uint16
	)
	if err := GetKeyboardState(&keystate[0]); err != nil {
		panic(fmt.Sprintf("win32: %v", err))
	}
	ret := ToUnicodeEx(vKey, uint32(scanCode), &keystate[0], &buf[0], int32(len(buf)), 0, keyboardLayout)
	if ret < 1 {
		return -1
	}
	return utf16.Decode(buf[:ret])[0]
}

func keyModifiers() (m key.Modifiers) {
	down := func(x int32) bool {
		// GetKeyState gets the key state at the time of the message, so this is what we want.
		return GetKeyState(x)&0x80 != 0
	}

	if down(VkControl) {
		m |= key.ModControl
	}
	if down(VkMenu) {
		m |= key.ModAlt
	}
	if down(VkShift) {
		m |= key.ModShift
	}
	if down(VkLwin) || down(VkRwin) {
		m |= key.ModMeta
	}
	return m
}

type ktab [256]key.Code

func (k *ktab) sendDown(h syscall.Handle, m uint32, w, l uintptr) uintptr {
	const prev = 1 << 30
	dir := key.DirNone
	if l&prev != prev {
		dir = key.DirPress
	}
	screen.Dev.Key <- key.Event{
		Rune:      readRune(uint32(w), byte(l>>16)),
		Code:      keytab[byte(w)],
		Modifiers: keyModifiers(),
		Direction: dir,
	}
	return 0
}
func (k *ktab) sendUp(h syscall.Handle, m uint32, w, l uintptr) uintptr {
	screen.Dev.Key <- key.Event{
		Rune:      readRune(uint32(w), byte(l>>16)),
		Code:      keytab[byte(w)],
		Modifiers: keyModifiers(),
		Direction: key.DirRelease,
	}
	return 0
}

var keytab = ktab{
	0x08: key.CodeDeleteBackspace,
	0x09: key.CodeTab,
	0x0D: key.CodeReturnEnter,
	0x10: key.CodeLeftShift,
	0x11: key.CodeLeftControl,
	0x12: key.CodeLeftAlt,
	0x14: key.CodeCapsLock,
	0x1B: key.CodeEscape,
	0x20: key.CodeSpacebar,
	0x21: key.CodePageUp,
	0x22: key.CodePageDown,
	0x23: key.CodeEnd,
	0x24: key.CodeHome,
	0x25: key.CodeLeftArrow,
	0x26: key.CodeUpArrow,
	0x27: key.CodeRightArrow,
	0x28: key.CodeDownArrow,
	0x2E: key.CodeDeleteForward,
	0x2F: key.CodeHelp,
	0x30: key.Code0,
	0x31: key.Code1,
	0x32: key.Code2,
	0x33: key.Code3,
	0x34: key.Code4,
	0x35: key.Code5,
	0x36: key.Code6,
	0x37: key.Code7,
	0x38: key.Code8,
	0x39: key.Code9,
	0x41: key.CodeA,
	0x42: key.CodeB,
	0x43: key.CodeC,
	0x44: key.CodeD,
	0x45: key.CodeE,
	0x46: key.CodeF,
	0x47: key.CodeG,
	0x48: key.CodeH,
	0x49: key.CodeI,
	0x4A: key.CodeJ,
	0x4B: key.CodeK,
	0x4C: key.CodeL,
	0x4D: key.CodeM,
	0x4E: key.CodeN,
	0x4F: key.CodeO,
	0x50: key.CodeP,
	0x51: key.CodeQ,
	0x52: key.CodeR,
	0x53: key.CodeS,
	0x54: key.CodeT,
	0x55: key.CodeU,
	0x56: key.CodeV,
	0x57: key.CodeW,
	0x58: key.CodeX,
	0x59: key.CodeY,
	0x5A: key.CodeZ,
	0x5B: key.CodeLeftGUI,
	0x5C: key.CodeRightGUI,
	0x60: key.CodeKeypad0,
	0x61: key.CodeKeypad1,
	0x62: key.CodeKeypad2,
	0x63: key.CodeKeypad3,
	0x64: key.CodeKeypad4,
	0x65: key.CodeKeypad5,
	0x66: key.CodeKeypad6,
	0x67: key.CodeKeypad7,
	0x68: key.CodeKeypad8,
	0x69: key.CodeKeypad9,
	0x6A: key.CodeKeypadAsterisk,
	0x6B: key.CodeKeypadPlusSign,
	0x6D: key.CodeKeypadHyphenMinus,
	0x6E: key.CodeFullStop,
	0x6F: key.CodeKeypadSlash,
	0x70: key.CodeF1,
	0x71: key.CodeF2,
	0x72: key.CodeF3,
	0x73: key.CodeF4,
	0x74: key.CodeF5,
	0x75: key.CodeF6,
	0x76: key.CodeF7,
	0x77: key.CodeF8,
	0x78: key.CodeF9,
	0x79: key.CodeF10,
	0x7A: key.CodeF11,
	0x7B: key.CodeF12,
	0x7C: key.CodeF13,
	0x7D: key.CodeF14,
	0x7E: key.CodeF15,
	0x7F: key.CodeF16,
	0x80: key.CodeF17,
	0x81: key.CodeF18,
	0x82: key.CodeF19,
	0x83: key.CodeF20,
	0x84: key.CodeF21,
	0x85: key.CodeF22,
	0x86: key.CodeF23,
	0x87: key.CodeF24,
	0x90: key.CodeKeypadNumLock,
	0xA0: key.CodeLeftShift,
	0xA1: key.CodeRightShift,
	0xA2: key.CodeLeftControl,
	0xA3: key.CodeRightControl,
	0xAD: key.CodeMute,
	0xAE: key.CodeVolumeDown,
	0xAF: key.CodeVolumeUp,
	0xBA: key.CodeSemicolon,
	0xBB: key.CodeEqualSign,
	0xBC: key.CodeComma,
	0xBD: key.CodeHyphenMinus,
	0xBE: key.CodeFullStop,
	0xBF: key.CodeSlash,
	0xC0: key.CodeGraveAccent,
	0xDB: key.CodeLeftSquareBracket,
	0xDC: key.CodeBackslash,
	0xDD: key.CodeRightSquareBracket,
	0xDE: key.CodeApostrophe,
	0xDF: key.CodeUnknown,
}
