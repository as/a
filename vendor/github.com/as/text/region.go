package text

import (
	"fmt"
	"os"
)

// Region3 maps r to 1 of 3 positions relative to q
//
// -1: Left of q0
//  0: Between q0 and q1, inclusive
//  1: Right of q1
//
// It must hold that q0 <= q1
func Region3(r, q0, q1 int64) int {
	if r <= q0 {
		return -1
	}
	if r > q1 {
		return 1
	}
	return 0
}

// Region5 maps the Dot x to 1 of 5 positions relative to q
//
// -2: x left of q
// -1: x0 left of q, x1 in q
//  0: x in q or q in x
//  1: x0 in q, x1 right of q
//  2: x right of q
//
// It must hold that r0<=r1 and q0<=q1
func Region5(r0, r1, q0, q1 int64) int {
	return Region3(r0, q0, q1) + Region3(r1, q0, q1)
	// Proof: Using results of Region3
	//          | * |  +
	// -----------------
	// (-1, -1) | 1 | -2
	// (-1,  0) | 0 | -1
	// ( 0,  0) | 0 |  0
	// (-1,  1) |-1 |  0
	// ( 0,  1) | 0 |  1
	// ( 1,  1) | 1 |  2
}

// Coherence returns the dot range q0:q1 should be in after an
// insert or delete operation occurs in the dot range r0:r1. If
// sign < 0, text is deleted, otherwise text is inserted.
func Coherence(sign int, r0, r1, q0, q1 int64) (int64, int64) {
	if sign < 0 {
		// deleting r0:r1 removes r1-r0 chars
		return coDelete(r0, r1, q0, q1)
	}
	return coInsert(r0, r1, q0, q1)
}

// CoherenceM runs Coherence over every dot range in the map q
func CoherenceM(sign int, r0, r1 int64, q map[string][2]int64) {
	if sign < 0 {
		for k, v := range clients {
			if len(v) != 2 {
				panic("!")
			}
			q0, q1 := coDelete(r0, r1, v[0], v[1])
			sel(k, q0, q1)
		}
		return
	}
	if r1-r0 == 0 {
		return
	}
	for k, v := range clients {
		if len(v) != 2 {
			panic("!")
		}
		q0, q1 := coInsert(r0, r1, v[0], v[1])
		sel(k, q0, q1)
	}
}

func coInsert(r0, r1, q0, q1 int64) (int64, int64) {
	dx := r1 - r0
	if dx == 0 {
		return q0, q1
	}
	switch Region3(r0, q0, q1) {
	case -1:
		q0 += dx
		q1 += dx
	case 0:
		q1 += dx
	case 1:
		// nop
	}
	return q0, q1
}

func coDelete(r0, r1, q0, q1 int64) (int64, int64) {
	dx := r1 - r0
	switch Region5(r0, r1, q0, q1) {
	case -2:
		q0 -= dx
		q1 -= dx
	case -1:
		q0 = r0
		q1 -= dx
	case 0:
		if q0 < r0 { // r in q
			q1 -= r1 - q0
		} else { // q in r
			q0 -= q0 - r0
			q1 = q0
		}
	case 1:
		q1 = r0
		//q1 -= q1 - r0
	case 2:
		// nop
	}
	return q0, q1
}

var (
	s       []byte
	clients = make(map[string][2]int64)
)

func insert(p []byte, q0 int64) (n int64) {
	x := append(p, s[q0:]...)
	s = append(s[:q0], x...)
	return int64(len(p))
}
func delete(q0, q1 int64) (n int64) {
	x := []int64{q0, q1}
	copy(s[x[0]:], s[x[1]+1:])
	s = s[:x[1]-x[0]]
	return q1 - q0
}
func sel(name string, q0, q1 int64) {
	clients[name] = [2]int64{q0, q1}
}
func print() {
	fmt.Printf("s=%q len=%d\n", s, len(s))
	for k, v := range clients {
		fmt.Printf("client[%q]=%v, %q\n", k, v, s[v[0]:v[1]+1])
	}
}

func test() {
	insert([]byte("The quick brown"), 0)
	sel("a", 1, 5)
	sel("b", 3, 8)
	sel("c", 7, 14)
	print()

	delete(2, 9)
	CoherenceM(-1, 2, 9, clients)
	print()

	insert([]byte("e quick "), 2)
	CoherenceM(1, 2, 9, clients)
	print()

	os.Exit(0)
}
