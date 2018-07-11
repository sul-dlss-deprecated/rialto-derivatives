package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordToDoc(t *testing.T) {
	record := "{ \"foo_si\": \"Hello world!\" }"
	doc := recordToDoc(record)

	assert.Equal(t, "Hello world!", doc.Get("foo_si"))
}
