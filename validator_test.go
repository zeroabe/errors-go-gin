package ginerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapWalk(t *testing.T) {
	sl1 := []FieldName{"foo", "bar", "baz"}
	sl2 := []FieldName{"foo", "bar", "bazz"}

	m := make(map[FieldName]interface{})

	c1 := mapWalk(m, sl1)
	c2 := mapWalk(m, sl2)

	_, ok := c1[FieldName("baz")]
	assert.True(t, true, ok)

	_, ok = c2[FieldName("bazz")]
	assert.True(t, true, ok)
}
