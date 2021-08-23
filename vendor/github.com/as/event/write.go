package event

type Write struct {
	Rec

	// A LHS Delete event ready to become another writer.
	Residue Record
}

func (r Write) Record() (p []byte, err error) {
	r.Kind = 'w'
	return r.Rec.Record()
}

func (e *Write) Coalesce(v Record) Record {
	if e.Residue != nil {
		return e.handleResidue(v)
	}
	switch v := v.(type) {
	case *Insert:
		//e.Residue =v
		return nil
	case *Delete:
		// If we reach here it means that the
		// input was W(D+I)+D. The next event
		// determines whether a write occurs
		if e.Q1 == v.Q0 {
			//e.Q1 = v.Q1
			//e.P = append(e.P, v.P...)
			e.Residue = v
			return e
		}
		e.Residue = nil
		return nil
	case *Write:
		if e.Q1 == v.Q0 {
			e.Q1 = v.Q1
			e.P = append(e.P, v.P...)
			//e.Residue = nil
			return e
		}
		panic("!")
	case interface{}:
		panic("bad interface")
	}
	return nil
}

func (e *Write) handleResidue(v Record) Record {
	if e.Residue == nil {
		panic("handleResidue: internal runtime error: no residue")
	}
	switch v := v.(type) {
	case *Insert:
		if e.Q1 == v.Q0 {
			e.Q1 = v.Q1
			e.P = append(e.P, v.P...)
			e.Residue = nil
			return e
		}
		return nil
	case *Delete:
		x := e.Residue.Coalesce(v)
		if x == nil {
			return nil
		}
		e.Residue = x
		return e
	case interface{}:
		panic("bad interface")
	}
	return nil
}
