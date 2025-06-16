package models

import (
	"time"
)

// Represents a row in 'checks' table.
type Check struct {
	Id    int64
	Name  string
	Owner string
}

// Represents a row in 'items' table with a [Check] obj.
type Item struct {
	Id      int64
	CheckId int64
	Name    string
	Owner   string
	Price   float64
}

type CheckWithItems struct {
	// TODO defer to use check.Id istead of Id
	Id    int64
	check Check
	items []Item
}

func (c *CheckWithItems) GetCheckName() string {
	return c.check.Name
}

func (c *CheckWithItems) GetCheckOwner() string {
	return c.check.Owner
}

// Returns a copy of check stroed in c.
func (c *CheckWithItems) GetCheck() Check {
	return c.check
}

func (c *CheckWithItems) GetItems() []Item {
	return c.items
}

func (c *CheckWithItems) SetCheckName(name string) {
	c.check.Name = name
}

func (c *CheckWithItems) SetCheckOwner(owner string) {
	c.check.Owner = owner
}

func (c *CheckWithItems) SetCheck(check *Check) {
	c.check.Name = check.Name
	c.check.Owner = check.Owner
}

func (c *CheckWithItems) SetItems(items []Item) {
	c.items = items
}

// Represents a row in 'sessions' table
type Session struct {
	// seesion_id
	Id int64
	// opened_at (format [time.DateTime])
	OpenedAt *time.Time
	// closed_at (format [time.DateTime])
	ClosedAt *time.Time
	// is_open
	IsOpen bool
}

// Represents a row in 'totals' table with a [Session] obj.
type SessionTotal struct {
	// TODO get rid of SessionId usage
	SessionId int64
	Total     float64
	Recipient string
	Amount    float64
	// Can be null...
	session *Session
}

func (st *SessionTotal) GetSessionId() int64 {
	return st.session.Id
}

func (st *SessionTotal) GetOpenedAtTime() *time.Time {
	return st.session.OpenedAt
}

func (st *SessionTotal) GetClosedAtTime() *time.Time {
	return st.session.ClosedAt
}

func (st *SessionTotal) SetSession(s *Session) {
	if s != nil {
		st.session = s
	} else {
		st.session = new(Session)
	}

}

type CheckTotal struct {
	Id int64

	OwnerId string

	OwnerTotal float64

	DebtorTotal float64

	Total float64
}
