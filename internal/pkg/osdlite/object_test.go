package osdlite

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	b := testBucket()
	o := testObj(b)
	p := testRandPayload()
	if _, err := o.Write(p); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(o.Payload, p) {
		t.Fatalf("write didn't succeed. Got: %s. Expected: %s", string(o.Payload), string(p))
	}
}

func TestRead(t *testing.T) {
	b := testBucket()
	obj := testObj(b)
	expected := testRandPayload()
	if _, err := obj.Write(expected); err != nil {
		t.Error(err)
	}
	got := make([]byte, len(expected))
	if _, err := obj.Read(got); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(got, expected) {
		t.Fatalf("read didn't succeed. Got: %s. Expected: %s", string(got), string(expected))
	}
}
