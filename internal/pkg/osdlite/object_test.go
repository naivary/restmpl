package osdlite

import (
	"bytes"
	"testing"
)

func TestWriteObject(t *testing.T) {
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

func TestReadObject(t *testing.T) {
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

func TestReadAsReader(t *testing.T) {
	var buf bytes.Buffer
	b := testBucket()
	o := testObj(b)
	if _, err := o.Write(testRandPayload()); err != nil {
		t.Error(err)
	}
	if _, err := buf.ReadFrom(o); err != nil {
		t.Error(err)
	}

	t.Log(buf.String())
	t.Log(o.String())
}

func TestReadFrom(t *testing.T) {
	p := testRandPayload()
	r := bytes.NewReader(p)
	b := testBucket()
	o := testObj(b)
	if _, err := o.ReadFrom(r); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(p, o.Payload) {
		t.Fatalf("paylaod is not the same as the expected. Got: %s. Expected: %s", o.PayloadString(), string(p))
	}

}

func TestWriteTo(t *testing.T) {
	p := testRandPayload()
	buf := new(bytes.Buffer)
	b := testBucket()
	o := testObj(b)
	if _, err := o.Write(p); err != nil {
		t.Error(err)
	}
	if _, err := o.WriteTo(buf); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(p, buf.Bytes()) {
		t.Fatalf("bytes are not equal. Got: %s. Expected: %s", buf.String(), string(p))
	}

}
