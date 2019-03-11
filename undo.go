package main

type Doer interface {
	Undo() bool
	Redo() bool
}

func do(undo bool) bool {
	w, ok := actTag.Window.(Doer)
	if !ok {
		return false
	}
	if undo {
		return w.Undo()
	}
	return w.Redo()
}
