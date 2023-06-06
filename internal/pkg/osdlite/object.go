package osdlite

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

type object struct {
	ID           string `json:"id"`
	CreatedAt    int64  `json:"createdAt"`
	LastModified int64  `json:"lastModified"`
	// Version of the file
	Version int `json:"version"`
	// Name of the file
	Name string `json:"name"`
	// Owner of the object
	Owner string `json:"owner"`
	// Bucket id in which the object is stores
	BucketID string `json:"bucketId"`
	// Payload of the object. Only manipulate this
	// if you really know what you are doing.
	// Otherwise use the defined function API.
	Payload []byte `json:"payload"`

	// current reading position
	pos int64
}

func NewObject(b *bucket, name string) *object {
	return &object{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now().Unix(),
		LastModified: time.Now().Unix(),
		BucketID:     b.ID,
		Owner:        b.Owner,
		Name:         name,
		Version:      1,
		Payload:      make([]byte, 0, 10),
	}
}

func (o *object) touch() {
	o.LastModified = time.Now().Unix()
}

func (o object) TableName() string {
	return "objects"
}

func (o *object) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(o, r)
}

func (o object) WriteTo(w io.Writer) (int64, error) {
	if o.Size() <= 0 {
		return 0, io.EOF
	}
	n, err := w.Write(o.Payload)
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

func (o *object) Write(p []byte) (int, error) {
	o.Payload = append(o.Payload, p...)
	o.touch()
	return len(p), nil
}

// Read implements the io.Reader interface.
// It will read the payload of the object,
// without removing the original payload.
func (o *object) Read(p []byte) (int, error) {
	if len(o.Payload) == 0 {
		return 0, io.EOF
	}
	n, err := bytes.NewReader(o.Payload[o.pos:]).Read(p)
	if err != nil {
		return 0, err
	}
	o.pos += int64(n)
	return n, nil
}

func (o object) String() string {
	return fmt.Sprintf("%s_%d", o.Name, o.CreatedAt)
}

// Size of the payload
func (o object) Size() int {
	return len(o.Payload)
}

func (o object) PayloadString() string {
	return string(o.Payload)
}

func (o *object) updateVersion() {
	o.Version += 1
}
