package models

import (
	"gitlab.com/thorchain/midgard/internal/clients/thorchain/types"
)

type BondType string

const (
	BondPaid     BondType = `bond_paid`
	BondReturned BondType = `bond_returned`
)

type EventBond struct {
	Event
	Amount   int64   `json:"amount"`
	BondType BondType `json:"bond_type"`
}

func NewBondEvent(bond types.EventBond, event types.Event) EventBond {
	return EventBond{
		Amount:   bond.Amount,
		BondType: BondType(bond.BondType),
		Event:    newEvent(event),
	}
}
