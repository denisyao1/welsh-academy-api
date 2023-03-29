package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSliceHasNoDuplicate(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		input  []string
		output bool
	}{
		{
			input:  []string{"banana", "orange", "banana"},
			output: false,
		},
		{
			input:  []string{"banana", "orange", "apple"},
			output: true,
		},
		{
			input:  []string{"orange", "apple", "banana", "banana"},
			output: false,
		},
	}

	for _, d := range data {
		assert.Equal(d.output, SliceHasNoDuplicate(d.input))
	}

}

func TestContains(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		elm    string
		slice  []string
		output bool
	}{
		{
			elm:    "orange",
			slice:  []string{"Orange", "Banana", "Apple"},
			output: false,
		},
		{
			elm:    "orange",
			slice:  []string{"Orange", "Banana", "Apple", "orange"},
			output: true,
		},
		{
			elm:    "orange",
			slice:  []string{},
			output: false,
		},
	}

	for _, d := range data {
		assert.Equal(d.output, Contains(d.elm, d.slice))
	}
}
