package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedInMemoryDB(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	db, err := NewInMemoryDB(true)
	assert.NoError(err)
	file := db.GetFile()
	assert.True(strings.Contains(file, "&cache=shared"))
}

func TestNotSharedInMemoryDB(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	db, err := NewInMemoryDB(false)
	assert.NoError(err)
	file := db.GetFile()
	assert.False(strings.Contains(file, "&cache=shared"))

}
