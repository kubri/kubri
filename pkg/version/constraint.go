package version

import (
	"errors"
	"strings"

	"golang.org/x/mod/semver"
)

type operator uint8

const (
	equal operator = iota
	notEqual
	greaterThan
	lessThan
	greaterThanEqual
	lessThanEqual
)

const separator = ","

func parseOperator(v string) (operator, string, bool) {
	var i int
	for ; i < len(v); i++ {
		if v[i] != ' ' {
			break
		}
	}

	if len(v[i:]) < 2 {
		return 0, v, false
	}

	switch v[i] {
	case 'v':
		return equal, v[i:], true
	case '=':
		return equal, v[i+1:], true
	case '!':
		if v[i+1] == '=' {
			return notEqual, v[i+2:], true
		}
		return 0, v, false
	case '>':
		if v[i+1] == '=' {
			return greaterThanEqual, v[i+2:], true
		}
		return greaterThan, v[i+1:], true
	case '<':
		if v[i+1] == '=' {
			return lessThanEqual, v[i+2:], true
		}
		return lessThan, v[i+1:], true
	}

	if v[i] >= '0' && v[i] <= '9' {
		return equal, v[i:], true
	}

	return 0, v, false
}

type Constraints []Constraint

func NewConstraint(v string) (Constraints, error) {
	if v == "" {
		return nil, nil
	}

	res := make([]Constraint, 0, strings.Count(v, separator)+1)
	var ok bool
	var c string

	for {
		c, v, _ = strings.Cut(v, separator)
		if c == "" {
			break
		}

		var op operator
		op, c, ok = parseOperator(c)
		if !ok {
			return nil, errors.New("invalid constraint: " + c)
		}

		res = append(res, Constraint{op: op, v: clean(c)})
	}

	return res, nil
}

func (c Constraints) Check(v string) bool {
	for _, c := range c {
		if !c.Check(v) {
			return false
		}
	}
	return true
}

type Constraint struct {
	op operator
	v  string
}

func (c Constraint) Check(v string) bool {
	v = clean(v)
	switch c.op {
	case equal:
		return semver.Compare(c.v, v) == 0
	case notEqual:
		return semver.Compare(c.v, v) != 0
	case greaterThan:
		return semver.Compare(c.v, v) < 0
	case lessThan:
		return semver.Compare(c.v, v) > 0
	case greaterThanEqual:
		return semver.Compare(c.v, v) <= 0
	case lessThanEqual:
		return semver.Compare(c.v, v) >= 0
	default:
		return false // Does not happen.
	}
}

func Check(constraint, version string) bool {
	c, err := NewConstraint(constraint)
	if err != nil {
		return false
	}
	return c.Check(version)
}

func clean(v string) string {
	i := 0
	for i < len(v) && (v[i] == ' ') {
		i++
	}
	n := len(v)
	for n > i && (v[n-1] == ' ') {
		n--
	}

	if v[i] != 'v' {
		return "v" + v[i:n]
	}

	return v[i:n]
}
