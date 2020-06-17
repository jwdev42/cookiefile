/* This file is part of the "cookiefile" library, ©2020 Jörg Walter
 *  This software is licensed under the "GNU Lesser General Public License version 3" */
package domainname

import (
	"fmt"
	"strings"
)

var letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var digits string = "0123456789"

const point = '.'
const hyphen = '-'

const (
	token_eof = iota
	token_hyphen
	token_point
	token_digit
	token_letter
	token_invalid
)

const (
	element_none = iota
	element_label
)

func tok2str(t int) string {
	switch t {
	case token_eof:
		return "EOF"
	case token_hyphen:
		return "HYPHEN"
	case token_point:
		return "POINT"
	case token_digit:
		return "DIGIT"
	case token_letter:
		return "LETTER"
	default:
		return "INVALID"
	}
}

type Validator struct {
	name string
	pos  int
}

func NewValidator(name string) *Validator {
	return &Validator{name: name}
}

func (r *Validator) err_unexpected_token(pos int, token int) error {
	return fmt.Errorf("%q, pos %d, unexpected token: %s", r.name, pos+1, tok2str(token))
}

func (r *Validator) token(index int) int {
	if index >= len(r.name) {
		return token_eof
	} else if r.name[index] == point {
		return token_point
	} else if r.name[index] == hyphen {
		return token_hyphen
	} else if strings.IndexByte(digits, r.name[index]) > -1 {
		return token_digit
	} else if strings.IndexByte(letters, r.name[index]) > -1 {
		return token_letter
	}
	return token_invalid
}

func (r *Validator) peek() int {
	return r.token(r.pos + 1)
}

func (r *Validator) next() int {
	r.pos++
	return r.token(r.pos)
}

func (r *Validator) parse_label() (int, error) {
	var hdl func() (int, error)
	hdl = func() (int, error) { //continues as long as the token is a hyphen, a digit or a letter
		t := r.token(r.pos)
		switch t {
		case token_digit:
			fallthrough
		case token_letter:
			r.next()
			return hdl()
		case token_hyphen:
			nt := r.peek()
			if nt != token_letter && nt != token_digit && nt != token_hyphen {
				return element_none, r.err_unexpected_token(r.pos+1, nt)
			}
			r.next()
			return hdl()
		case token_point:
			return r.parse_subdomain(element_label)
		case token_eof:
			return element_label, nil
		default:
			return element_none, r.err_unexpected_token(r.pos, t)
		}
	}
	t := r.token(r.pos)
	switch t {
	case token_digit:
		fallthrough
	case token_letter:
		r.next()
		return hdl()
	}
	return element_none, r.err_unexpected_token(r.pos, t)
}

func (r *Validator) parse_subdomain(e int) (int, error) {
	t := r.token(r.pos)
	switch e {
	case element_none:
		if t == token_point {
			r.next()
		}
		return r.parse_label()
	case element_label:
		if t != token_point {
			return element_none, r.err_unexpected_token(r.pos, t)
		}
		r.next()
		return r.parse_label()
	}
	panic("function called with an invalid argument")
}

func (r *Validator) Validate() error {
	if _, err := r.parse_subdomain(element_none); err != nil {
		return fmt.Errorf("Domain name validation failed: %v", err)
	}
	return nil
}
