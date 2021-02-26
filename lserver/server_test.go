package lserver

import (
	"bytes"
	"testing"
)

type Result struct {
	advance int
	data    []byte
	err     error
}

type ArgsSent struct {
	payload []byte
	eof     bool
}

var argsSent = []ArgsSent{
	{[]byte("+++idspispopd_____ABCAAAAAAAAAAAAAAZ789102jsAUoek99"), false},
}

var expected = []Result{
	{1, []byte("ABCAAAAAAAAAAAAAAZ789102jsAUoek99"), error(nil)},
}

func TestLSplit(t *testing.T) {
	// add tests
	for i, args := range argsSent {
		r0, r1, r2 := lSplit(args.payload, args.eof)

		exp0, exp1, exp2 := expected[i].advance, expected[i].data, expected[i].err

		if r0 != exp0 || !bytes.Equal(r1, exp1) || r2 != exp2 {
			t.Error(r0, string(r1), r2, " != ", exp0, string(exp1), exp2)
		}
	}
}
