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

func TestServerClient(t *testing.T) {
	srv, err := Serve("tcp", "localhost:0")
	if err != nil {
		t.Logf("serve: %s\n", err)
		t.Fail()
	}
	defer srv.Close()

	addr := srv.fd.Addr()

	client, err := Dial(addr.Network(), addr.String())
	if err != nil {
		t.Logf("dial: %s\n", err)
		t.Fail()
	}

	name := "fs.net.test.get"
	want := "take me to your leader"

	testput(t, client, name, []byte(want), false)
	have := testget(t, client, name, true)

	if string(have) != want {
		t.Logf("data mismatch: have %q, want %q\n", have, want)
		t.Fail()
	}
}
