package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type person struct {
	name string
	age  int
}

func TestHash_person(t *testing.T) {
	data := person{
		name: "elton",
		age:  34,
	}

	r, err := Hash(data)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	t.Log(r)
}

func TestHash_empty_person(t *testing.T) {
	data := person{}

	r, err := Hash(data)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	t.Log(r)
}

func TestHash_string(t *testing.T) {
	data := "hello"

	r, err := Hash(data)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	t.Log(r)
}
