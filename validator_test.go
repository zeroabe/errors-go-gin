package ginerrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapWalk(t *testing.T) {
	sl1 := []FieldName{"foo", "bar", "baz"}
	sl2 := []FieldName{"foo", "bar", "bazz"}

	m := make(map[FieldName]interface{})
	er := fmt.Errorf("some err")

	c1 := mapWalk(m, sl1, er)
	c2 := mapWalk(m, sl2, er)

	_, ok := c1[FieldName("baz")]
	assert.True(t, true, ok)

	_, ok = c2[FieldName("bazz")]
	assert.True(t, true, ok)
}
