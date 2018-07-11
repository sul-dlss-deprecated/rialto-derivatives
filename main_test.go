package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordToDoc(t *testing.T) {
	record := "{ \"dc:title\": \"Hello world!\" }"
	doc := recordToResource(record)

	assert.Equal(t, "Hello world!", doc.Title())
}
