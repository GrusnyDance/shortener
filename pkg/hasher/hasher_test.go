package hasher_test

import (
	"github.com/stretchr/testify/assert"
	"shortener/pkg/hasher"
	"testing"
)

func TestLen(t *testing.T) {
	linkToHash := "hahaha"
	hashedLink := hasher.Apply(linkToHash)
	assert.Len(t, hashedLink, 10, "len is not 10, but %v", len(hashedLink))
}

func TestLenDot(t *testing.T) {
	linkToHash := "."
	hashedLink := hasher.Apply(linkToHash)
	assert.Len(t, hashedLink, 10, "len is not 10, but %v", len(hashedLink))
}

func TestLenEmptyString(t *testing.T) {
	linkToHash := ""
	hashedLink := hasher.Apply(linkToHash)
	assert.Len(t, hashedLink, 10, "len is not 10, but %v", len(hashedLink))
}

func TestFromStringToHash(t *testing.T) {
	linkToHash := "lalala"
	hashedLink := hasher.Apply(linkToHash)
	assert.Equal(t, hashedLink, "8ja9JS8xBK")
}

func TestIfConsistent(t *testing.T) {
	linkToHash := "hahaha"
	hashedLink1 := hasher.Apply(linkToHash)
	hashedLink2 := hasher.Apply(linkToHash)
	hashedLink3 := hasher.Apply(linkToHash)
	hashedLink4 := hasher.Apply(linkToHash)

	assert.Equal(t, hashedLink1, "P8HO4BKzxV")
	assert.Equal(t, hashedLink2, "P8HO4BKzxV")
	assert.Equal(t, hashedLink3, "P8HO4BKzxV")
	assert.Equal(t, hashedLink4, "P8HO4BKzxV")
}
