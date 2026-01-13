package static

import (
	"strings"
	"time"

	"github.com/Tesorp1X/chipi-bot/utils"
	"github.com/Tesorp1X/chipi-bot/utils/reader"
)

type CallbackAction string

func (a CallbackAction) String() string {
	return string(a)
}

func (a CallbackAction) DataMatches(data string) bool {
	cringePrefix := "\f" + a.String()
	return data == cringePrefix || strings.HasPrefix(data, cringePrefix+"|")
}

const (
	//CallbackAction<name> CallbackAction = "name"
	CallbackActionSelector CallbackAction = "Selector"
	CallbackActionEditItem CallbackAction = "EditItem"
)

// Represents a record in checks table
type Check struct {
	// Record id in table
	Id int64

	SessionId int64

	// Check name
	Name string
	// Shop name
	OrgName string
	// Who paid for this check
	Owner string

	Total    float64
	TotalPau float64
	TotalLiz float64

	Date *time.Time
}

func (a *Check) IsEqual(b *Check) bool {
	fact := a.Id == b.Id && a.SessionId == b.SessionId
	fact = fact && a.Name == b.Name && a.OrgName == b.OrgName
	fact = fact && a.Total == b.Total && a.TotalLiz == b.TotalLiz
	fact = fact && a.TotalPau == b.TotalPau
	fact = fact && a.Date.Equal(*b.Date)

	return fact
}

// Creates an object of `Check`-type from `reader.CheckData` object.
// Some field are left unfilled.
func CreateCheckFromCheckData(cd *reader.CheckData) *Check {
	return &Check{
		Name:    utils.AssumeCheckName(cd.OrgName),
		OrgName: cd.OrgName,
		Total:   cd.Total,
		Date:    &cd.TimeOfCreation,
	}
}

type Item struct {
	// Record id in table
	Id int64

	CheckId int64
	// Item's name
	Name string

	// Item's owner
	Owner string

	// Cost for 1.0 amount of this item
	Price float64

	Amount float64

	// Price * Amount value
	Subtotal float64
}

func (a *Item) IsEqual(b *Item) bool {
	fact := a.Id == b.Id && a.CheckId == b.CheckId
	fact = fact && a.Name == b.Name && a.Owner == b.Owner
	fact = fact && a.Price == b.Price && a.Amount == b.Amount
	fact = fact && a.Subtotal == b.Subtotal

	return fact
}

func CreateItemsFromCheckData(cd *reader.CheckData) []*Item {
	var items []*Item
	for _, cdItem := range cd.Items {
		items = append(items, &Item{
			Name:     cdItem.Name,
			Price:    cdItem.Price,
			Amount:   cdItem.Amount,
			Subtotal: cdItem.SubTotal,
		})
	}

	return items
}
