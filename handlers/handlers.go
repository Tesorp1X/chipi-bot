package handlers

import (
	"context"
	"os"
	"strconv"

	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func CancelHandler(c tele.Context, state fsm.Context) error {
	_ = state.Finish(context.TODO(), c.Data() != "")
	return c.Send("Canceled!")
}

func HelloHandler(c tele.Context, state fsm.Context) error {
	return c.Send("Hello, " + strconv.FormatInt(c.Message().Sender.ID, 10))
}

func NewCheckHandler(c tele.Context, state fsm.Context) error {
	if err := state.SetState(context.TODO(), models.StateWaitForCheckName); err != nil {
		return c.Send(models.ErrorSomethingWentWrong)
	}
	val, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		return c.Send(models.ErrorSomethingWentWrong)
	}
	if err := state.Update(context.Background(), models.SESSION_ID, val); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}
	return c.Send("Хорошо, как назовем новый чек?👀")
}

// NewCheckNameResponseHandler is called when user sends a new check name.
// It updates the check name in the database and returns to 'check-scroll' menu with state ShowingChecks.
func NewCheckNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	var check *models.CheckWithItems
	if err := state.Data(context.TODO(), models.CHECK, &check); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := db.EditCheckName(check.Id, msgText); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSavingInDB)
	}

	if err := state.Update(context.TODO(), models.CHECK, nil); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	session, err := getChecksForCurrentSession(c)
	if err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.Update(context.TODO(), models.CHECKS, session); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	var currentIndex int
	if err := state.Data(context.TODO(), models.CURRENT_INDEX_CHECKS, &currentIndex); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		return c.Send(models.ErrorSetState)
	}

	msg := "Название изменено!\n\n" + util.GetCheckWithItemsResponse(*session.Checks[currentIndex], currentIndex)
	kb := models.GetScrollKb()

	return c.EditOrReply(msg, kb)
}

func CheckNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}

	if err := state.Update(context.TODO(), models.CHECK_NAME, msgText); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateWaitForCheckOwner); err != nil {
		return c.Send(models.ErrorSomethingWentWrong)
	}

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "Лиз :3",
			Unique: models.CallbackActionCheckOwner.String(),
			Data:   models.OWNER_LIZ,
		},
		models.Button{
			BtnTxt: "Пау <3",
			Unique: models.CallbackActionCheckOwner.String(),
			Data:   models.OWNER_PAU,
		},
	)

	return c.Send("Хорошо. Кто заплатил?🤑", kb)
}

func ItemNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	if err := state.Update(context.TODO(), models.ITEM_NAME, msgText); err != nil {
		state.Finish(context.Background(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateWaitForItemPrice); err != nil {
		return c.Send(models.ErrorSomethingWentWrong)
	}

	return c.Send("Сколько это стоило?\n<i>Можно указать кол-во товаров вот так: 2*68 (2 товара по 68 руб)</i>")
}

func ItemPriceResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()

	price, err := util.ParsePrice(msgText)

	if err != nil {
		switch err {
		case models.ErrItemPriceMultiplierNotSingleInt:
			return c.Send(models.ErrorAmountOfItemsMustBeANumberMsg)
		case models.ErrItemPriceNotSingleIntOrFloat:
			return c.Send(models.ErrorItemPriceMustBeANumberMsg)
		case models.ErrItemPriceWrongFormat:
			return c.Send(models.AmountOfItemsHelpMsg)
		}
	}

	if err := state.Update(context.TODO(), models.ITEM_PRICE, price); err != nil {
		state.Finish(context.Background(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateWaitForItemOwner); err != nil {
		state.Finish(context.Background(), true)
		return c.Send(models.ErrorSomethingWentWrong)
	}

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "Лиз :3",
			Unique: models.CallbackActionItemOwner.String(),
			Data:   models.OWNER_LIZ,
		},
		models.Button{
			BtnTxt: "Пау <3",
			Unique: models.CallbackActionItemOwner.String(),
			Data:   models.OWNER_PAU,
		},
		models.Button{
			BtnTxt: "Общий",
			Unique: models.CallbackActionItemOwner.String(),
			Data:   models.OWNER_BOTH,
		},
	)
	return c.Send("Хорошо. Чей это товар?😺", kb)
}

// /current -- shows how much both payed and who owns money to whom and how much.
func ShowCurrentTotalCommand(c tele.Context, state fsm.Context) error {

	session, err := getChecksForCurrentSession(c)
	if err != nil {
		return err
	}

	sessionTotal := util.CalculateSessionTotal(session.SessionId, session.Checks)

	msg := util.GetTotalResponse(sessionTotal, true)

	return c.Send(msg)
}

// /finish -- finishes current session and makes a record in totals table.
// Also notifies another person about it.
func FinishSessionCommand(c tele.Context, state fsm.Context) error {

	session, err := getChecksForCurrentSession(c)
	if err != nil {
		return err
	}

	sessionTotal := util.CalculateSessionTotal(session.SessionId, session.Checks)

	msg := util.GetTotalResponse(sessionTotal, false)

	if err := db.CreateTotal(sessionTotal); err != nil {
		return c.Send(models.ErrorSavingInDB)
	}

	if err := db.FinishSession(session.SessionId); err != nil {
		return c.Send(models.ErrorSavingInDB)
	}

	// sending notification to another person
	adminsList := util.ExtractAdminsIDs(os.Getenv("ADMINS"))
	for _, adminId := range adminsList {
		if adminId != c.Sender().ID {
			broadcastMsg := "Сессия была завершена!\n" + msg
			if _, err := c.Bot().Send(&tele.User{ID: adminId}, broadcastMsg); err != nil {
				return err
			}
		}
	}

	return c.Send(msg)
}

// Handler for '/show <arg>' command
func ShowCommand(c tele.Context, state fsm.Context) error {
	args := c.Args()
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}
	switch arg {
	case "checks":

		return showChecks(c, state)

	case "totals":

		return showTotals(c, state)

	default:

		return showHelp(c)

	}
}

func showHelp(c tele.Context) error {
	msg := "<b>Команда /show требует аргумента. Например:</b>\n"
	checksHelp := "<i>/show checks  &#8212; покажет чеки</i>\n\n"
	totalsHelp := "<i>/show totals &#8212; покажет отчеты о прошлых сессиях</>\n\n"

	msg += checksHelp + totalsHelp

	kb := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnChecks := kb.Text("/show checks")
	btnTotals := kb.Text("/show totals")

	kb.Reply(
		kb.Row(btnChecks),
		kb.Row(btnTotals),
	)

	return c.Send(msg, kb)
}

func showChecks(c tele.Context, state fsm.Context) error {
	// Trying to get session from context.
	var session *checksForSession
	if err := state.Data(context.TODO(), models.CHECKS, &session); err != nil || len(session.Checks) == 0 {
		// If nothing is stored in context or slices len is 0, make a request to db.
		session, err = getChecksForCurrentSession(c)
		if err != nil {
			return err
		}
		// length is still zero, then there must be no checks for this session yet.
		if len(session.Checks) == 0 {
			return c.Send("В текущей сессии пока что нет чеков.")
		}
	}

	// Context should be short-lived (few mins).
	// TODO make it short-lived
	if err := state.Update(context.TODO(), models.CHECKS, session); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	var currentIndex int = len(session.Checks) - 1 // Starting from the last one
	if err := state.Update(context.TODO(), models.CURRENT_INDEX_CHECKS, currentIndex); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	kb := models.GetScrollKb()

	return c.Send(util.GetCheckWithItemsResponse(*session.Checks[currentIndex], currentIndex), kb)
}

func showTotals(c tele.Context, state fsm.Context) error {

	totals, err := db.GetAllSessionTotals()
	if err != nil {
		return c.Send(models.ErrorSomethingWentWrong)
	}

	var response string
	var kb *tele.ReplyMarkup

	if len(totals) == 0 {
		// TODO work on that
		response = "No totals"
		state.Finish(context.TODO(), true)
		return c.Send(response, kb)
	}

	// Context should be short-lived (few mins).
	// TODO make it short-lived
	if err := state.Update(context.TODO(), models.SESSION_TOTALS, totals); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	var currentIndex int
	if err := state.Update(context.TODO(), models.CURRENT_INDEX_TOTALS, currentIndex); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if len(totals) <= currentIndex {
		// setting to the last element
		currentIndex = len(totals) - 1
	}

	if currentIndex < 0 {
		// wtf, but resetting index
		currentIndex = 0
	}

	if err := state.SetState(context.TODO(), models.StateShowingTotals); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}
	response = util.GetShowTotalsResponse(totals[currentIndex])

	kb = models.CreateSelectorInlineKb(
		3,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_BACK,
		},
		models.Button{
			BtnTxt: "Чеки",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_CHECKS,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_FORWARD,
		},
	)

	return c.Send(response, kb)
}

type checksForSession struct {
	SessionId int64
	Checks    []*models.CheckWithItems
}

// Helper function. that gets current sessionId and pulls all checks for it from db.
// Can return errors only occurred during [Bot.Send()]
func getChecksForCurrentSession(c tele.Context) (*checksForSession, error) {
	sessionId, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		var err error
		sessionId, err = db.GetSessionId()
		if err != nil {
			return nil, c.Send(models.ErrorSomethingWentWrong)
		}
	}

	checks, err := db.GetAllChecksWithItemsForSessionId(sessionId)
	if err != nil {
		return nil, c.Send(err)
	}

	return &checksForSession{SessionId: sessionId, Checks: checks}, nil
}
