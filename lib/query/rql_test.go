package query

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimple(t *testing.T) {
	query := &Query{}
	err := parser.ParseString("", "foobar", query)
	require.NoError(t, err)
	sql := query.ToSql()
	assert.Equal(t, "?? LIKE '%foobar%'", sql)
}

func TestString(t *testing.T) {
	query := &Query{}
	err := parser.ParseString("", "\"foo and bar\"", query)
	require.NoError(t, err)
	sql := query.ToSql()
	assert.Equal(t, "?? LIKE '%foo and bar%'", sql)
}

func TestSimpleOr(t *testing.T) {
	query := &Query{}

	err := parser.ParseString("", "foo OR bar", query)
	require.NoError(t, err)
	sql := query.ToSql()
	assert.Equal(t, "?? LIKE '%foo%' OR ?? LIKE '%bar%'", sql)

	err = parser.ParseString("", "foo OR bar | bla", query)
	require.NoError(t, err)
	sql = query.ToSql()
	assert.Equal(t, "?? LIKE '%foo%' OR ?? LIKE '%bar%' OR ?? LIKE '%bla%'", sql)
}

func TestSimpleOrAnd(t *testing.T) {
	query := &Query{}

	err := parser.ParseString("", "foo OR bar AND bla", query)
	require.NoError(t, err)
	sql := query.ToSql()
	assert.Equal(t, "?? LIKE '%foo%' OR ?? LIKE '%bar%' AND ?? LIKE '%bla%'", sql)

	err = parser.ParseString("", "foo AND bar OR bla", query)
	require.NoError(t, err)
	sql = query.ToSql()
	assert.Equal(t, "?? LIKE '%foo%' AND ?? LIKE '%bar%' OR ?? LIKE '%bla%'", sql)
}

func TestSimpleNot(t *testing.T) {
	query := &Query{}

	err := parser.ParseString("", "!foo", query)
	require.NoError(t, err)
	sql := query.ToSql()
	assert.Equal(t, "NOT ?? LIKE '%foo%'", sql)

	err = parser.ParseString("", "!foo OR !bla", query)
	require.NoError(t, err)
	sql = query.ToSql()
	assert.Equal(t, "NOT ?? LIKE '%foo%' OR NOT ?? LIKE '%bla%'", sql)
}
