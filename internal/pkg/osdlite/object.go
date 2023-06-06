package osdlite

import (
	"bytes"
	"time"

	"github.com/google/uuid"
)

type object struct {
	ID           string `json:"id"`
	CreatedAt    int64  `json:"createdAt"`
	LastModified int64  `json:"lastModified"`
	Version      int    `json:"version"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	BucketID     string `json:"bucketId"`
	Payload      []byte `json:"payload"`

	// current reading position
	pos int
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

func (o *object) Touch() {
	o.LastModified = time.Now().Unix()
}

func (o object) TableName() string {
	return "objects"
}

func (o *object) Write(p []byte) (int, error) {
	o.Payload = append(o.Payload, p...)
	return len(p), nil
}

// Read implements the io.Reader interface.
// It will read the payload of the object,
// without removing the original payload.
func (o *object) Read(p []byte) (int, error) {
	n, err := bytes.NewReader(o.Payload[o.pos:]).Read(p)
	if err != nil {
		return 0, err
	}
	o.pos += n
	return n, nil
}

func (o *object) String() string {
	return string(o.Payload)
}

// Size of the payload
func (o *object) Size() int {
	return len(o.Payload)
}
