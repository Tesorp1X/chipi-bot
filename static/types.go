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
	CallbackActionKeep CallbackAction = "keep"
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
