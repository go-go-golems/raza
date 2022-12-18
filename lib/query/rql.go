package query

import (
	"github.com/alecthomas/participle/v2"
	"strings"
)

// raza query language

type Query struct {
	Disjunction *Disjunction `@@`
}

type Disjunction struct {
	Conjunction *Conjunction `@@`
	Op          string       `[ @( "OR" | "|" )`
	Next        *Disjunction `  @@ ]`
}

type Conjunction struct {
	Unary *Unary       `@@`
	Op    string       `[ @( "AND" | "&" )`
	Next  *Conjunction `  @@ ]`
}

type Unary struct {
	Op    string `  ( @( "!" | "-" | "NOT" )`
	Unary *Unary `    @@ )`
	Term  *Term  `| @@`
}

type Term struct {
	Simple   *string `  @Ident`
	String   *string `| @String`
	SubQuery *Query  `| "(" @@ ")" `
}

var parser = participle.MustBuild(&Query{}, participle.UseLookahead(2))

func (d *Disjunction) ToSql() string {
	var s string
	if d.Conjunction != nil {
		s = d.Conjunction.ToSql()
	}
	if d.Op != "" {
		s += " OR "
	}
	if d.Next != nil {
		s += d.Next.ToSql()
	}
	return s
}

func (c *Conjunction) ToSql() string {
	var s string
	if c.Unary != nil {
		s = c.Unary.ToSql()
	}
	if c.Op != "" {
		s += " AND "
	}
	if c.Next != nil {
		s += c.Next.ToSql()
	}
	return s
}

func (u *Unary) ToSql() string {
	var s string
	if u.Op != "" {
		s += "NOT "
		if u.Unary != nil {
			s += u.Unary.ToSql()
		}
		return s
	}

	if u.Term != nil {
		return u.Term.ToSql()
	}

	return s
}

func (t *Term) ToSql() string {
	if t.Simple != nil {
		return "?? LIKE '%" + *t.Simple + "%'"
	}
	if t.String != nil {
		return "?? LIKE '%" + strings.Trim(*t.String, "\"") + "%'"
	}
	if t.SubQuery != nil {
		return "(" + t.SubQuery.ToSql() + ")"
	}
	return ""
}

func (q *Query) ToSql() string {
	if q.Disjunction != nil {
		return q.Disjunction.ToSql()
	}
	return ""
}
