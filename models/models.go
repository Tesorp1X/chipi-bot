package models

type Check struct {
	Name  string
	Owner string
}

type Item struct {
	CheckId int64
	Name    string
	Owner   string
	Price   float64
}

type CheckWithItems struct {
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

type SessionTotal struct {
	SessionId int64
	Total     float64
	Recipient string
	Amount    float64
}

type CheckTotal struct {
	Id int64

	OwnerId string

	OwnerTotal float64

	DebtorTotal float64

	Total float64
}
