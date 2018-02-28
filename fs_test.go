package main

import (
	"os"
	"testing"
)

func testput(t *testing.T, fs FS, name string, data []byte, rm bool) {
	t.Helper()
	err := fs.Put(name, data)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !rm {
		return
	}
	err = os.Remove(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func testget(t *testing.T, fs FS, name string, rm bool) (data []byte) {
	t.Helper()

	have, err := fs.Get(name)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if rm {
		if err = os.Remove(name); err != nil {
			t.Log(err)
			t.Fail()
		}
	}

	return have
}

func TestFSPut(t *testing.T) {
	l := &LocalFS{}
	testput(t, l, "fs.test.write", []byte("hello world"), true)
}

func TestFSGet(t *testing.T) {
	l := &LocalFS{}
	name := "fs.test.get"
	want := "take me to your leader"

	testput(t, l, name, []byte(want), false)
	have := testget(t, l, name, false)

	if string(have) != want {
		t.Logf("data mismatch: have %q, want %q\n", have, want)
		t.Fail()
	}
}
