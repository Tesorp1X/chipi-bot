package static

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tesorp1X/chipi-bot/utils"
	"github.com/Tesorp1X/chipi-bot/utils/reader"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

type CallbackAction string

func (a CallbackAction) String() string {
	return string(a)
}

func (a CallbackAction) DataMatches(data string) bool {
	cringePrefix := "\f" + a.String()
	return data == cringePrefix || strings.HasPrefix(data, cringePrefix+"|")
}

// Extracts data from [Callback.Data] by removing prefix '\f + CallbackAction + |'
func (a CallbackAction) GetData(data string) string {
	return utils.ExtractCallbackData(data)
}

const (
	//CallbackAction<name> CallbackAction = "name"
	CallbackActionSelector   CallbackAction = "Selector"
	CallbackActionNavigation CallbackAction = "Navigation"

	CallbackActionEditItem        CallbackAction = "EditItem"
	CallbackActionEditUnsavedItem CallbackAction = "EditUnsavedItem"

	CallbackActionEditCheck        CallbackAction = "EditCheck"
	CallbackActionEditUnsavedCheck CallbackAction = "EditUnsavedCheck"

	CallbackActionUnknown CallbackAction = "Unknown"
)

// Returns valid CallbackAction based on raw callback data.
// If unable to match with existing actions, then CallbackActionUnknown will be returned.
func GetCallbackActionFromRawData(rawData string) CallbackAction {
	actions := []CallbackAction{
		CallbackActionSelector, CallbackActionNavigation,
		CallbackActionEditItem, CallbackActionEditUnsavedItem,
		CallbackActionEditCheck, CallbackActionEditUnsavedCheck,
	}

	for _, a := range actions {
		if a.DataMatches(rawData) {
			return a
		}
	}

	return CallbackActionUnknown
}

type actionAndStates struct {
	action CallbackAction
	states []fsm.State
}

type actionsToStates []*actionAndStates

func GetCallbackActionBasedOnState(userState fsm.State) CallbackAction {
	ats := actionsToStates{
		&actionAndStates{
			action: CallbackActionEditUnsavedCheck,
			states: []fsm.State{
				StateWaitingForCheckConfirmationUnsaved, StateEditingCheckUnsaved,
				StateWaitForNewCheckNameUnsaved, StateWaitForCheckCreationDateUnsaved,
				StateWaitForCheckOwnerUnsaved, StateEditingCheckUnsaved,
			},
		},
		&actionAndStates{
			action: CallbackActionEditUnsavedItem,
			states: []fsm.State{
				StateShowingAnItemUnsaved, StateEditingAnItemUnsaved,
			},
		},
	}

	statesToCbActionsMap := make(map[fsm.State]CallbackAction)
	for _, a := range ats {
		for _, state := range a.states {
			statesToCbActionsMap[state] = a.action
		}
	}

	if cbAction, ok := statesToCbActionsMap[userState]; ok {
		return cbAction
	}

	return CallbackActionUnknown
}

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

// Calculates and assigns Total, TotalLiz and TotalPau fields based on a given items slice.
// If there any item without an owner or with invalid one, the method will return an error.
func (c *Check) CalculateTotals(items []*Item) error {
	var total, liz, pau float64

	for _, item := range items {
		total += item.Subtotal
		switch item.Owner {
		case CallbackOwnerLiz:
			liz += item.Subtotal
		case CallbackOwnerPau:
			pau += item.Subtotal
		case CallbackOwnerBoth:
			liz += item.Subtotal / 2
			pau += item.Subtotal / 2
		default:
			return fmt.Errorf("error in Check.CalculateTotals(): invalid owner for item %v", *item)
		}
	}

	c.Total = total
	c.TotalLiz = liz
	c.TotalPau = pau

	return nil
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
