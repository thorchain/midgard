package common

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	BNBSymbol     = Symbol("BNB")
	Rune67CSymbol = Symbol("RUNE-67C")
	RuneB1ASymbol = Symbol("RUNE-B1A")
)

var isAlpha = regexp.MustCompile(`^[A-Za-z0-9-]+$`).MatchString

type Symbol string

func NewSymbol(input string) (Symbol, error) {
	if !isAlpha(input) {
		return "", fmt.Errorf("Invalid symbol")
	}
	return Symbol(input), nil
}

func (s Symbol) Ticker() Ticker {
	parts := strings.Split(s.String(), "-")
	ticker, _ := NewTicker(parts[0])
	return ticker
}

func (s Symbol) Equals(s2 Symbol) bool {
	return strings.EqualFold(s.String(), s2.String())
}

func (s Symbol) IsEmpty() bool {
	return strings.TrimSpace(s.String()) == ""
}

func (s Symbol) String() string {
	// uppercasing again just in case someone created a ticker via Chain("rune")
	return strings.ToUpper(string(s))
}

func (s Symbol) IsMiniToken() bool {
	parts := strings.Split(s.String(), "-")
	if len(parts) == 1 {
		return false
	}
	return len(parts[1]) == 4 && strings.HasSuffix(s.String(), "M")
}

func IsBNBSymbol(s Symbol) bool {
	return s.Equals(BNBSymbol)
}

func IsRuneSymbol(s Symbol) bool {
	return s.Equals(RuneB1ASymbol) || s.Equals(Rune67CSymbol)
}
