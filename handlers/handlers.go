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
		return c.Send(models.ErrorSometingWentWrong)
	}
	val, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		return c.Send(models.ErrorSometingWentWrong)
	}
	if err := state.Update(context.Background(), models.SESSION_ID, val); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}
	return c.Send("Хорошо, как назовем новый чек?👀")
}

func CheckNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.CHECK_NAME, msgText)
	if err := state.SetState(context.TODO(), models.StateWaitForCheckOwner); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	selector := models.CheckOwnershipSelectorInlineKb(
		"Liz :3", models.CallbackActionCheckOwner.String(), models.OWNER_LIZ,
		"Пау <3", models.CallbackActionCheckOwner.String(), models.OWNER_PAU,
	)
	return c.Send("Хорошо. Кто заплатил?🤑", selector)
}

func ItemNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.ITEM_NAME, msgText)

	if err := state.SetState(context.TODO(), models.StateWaitForItemPrice); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	return c.Send("Сколько это столо?")
}

func ItemPriceResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	var (
		price float64
		err   error
	)

	if price, err = strconv.ParseFloat(msgText, 64); err != nil {
		return c.Send(models.ErrorItemPriceMustBeANumberMsg)
	}
	state.Update(context.TODO(), models.ITEM_PRICE, price)

	if err := state.SetState(context.TODO(), models.StateWaitForItemOwner); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	selector := models.ItemOwnershipSelectorInlineKb(
		"Liz :3", models.CallbackActionItemOwner.String(), models.OWNER_LIZ,
		"Пау <3", models.CallbackActionItemOwner.String(), models.OWNER_PAU,
		"Общий", models.CallbackActionItemOwner.String(), models.OWNER_BOTH,
	)
	return c.Send("Хорошо. Чей это товар?😺", selector)
}

// /current -- shows how much both payed and who owns money to whom and how much.
func ShowCurrentTotalCommand(c tele.Context, state fsm.Context) error {
	// sessionId, ok := c.Get(models.SESSION_ID).(int64)
	// if !ok {
	// 	var err error
	// 	sessionId, err = db.GetSessionId()
	// 	if err != nil {
	// 		return c.Send(models.ErrorSometingWentWrong)
	// 	}
	// }

	// checks, err := db.GetAllChecksWithItemsForSesssionId(sessionId)
	// if err != nil {
	// 	return c.Send(err)
	// }

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
	// sessionId, ok := c.Get(models.SESSION_ID).(int64)
	// if !ok {
	// 	var err error
	// 	sessionId, err = db.GetSessionId()
	// 	if err != nil {
	// 		return c.Send(models.ErrorSometingWentWrong)
	// 	}
	// }

	// checks, err := db.GetAllChecksWithItemsForSesssionId(sessionId)
	// if err != nil {
	// 	return c.Send(err)
	// }

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
		msg := "Команда /show требует аргумента. Например:\n/show checks -- покажет чеки\nДругие аргументы пока что в разработке."
		kb := &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
		btnChecks := kb.Text("/show checks")

		kb.Reply(
			kb.Row(btnChecks),
		)

		return c.Send(msg, kb)
	}
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
	}

	var currentIndex int = 0
	// If currentIndex is not stored in context, then it will be just zero.
	state.Data(context.TODO(), models.CURRENT_INDEX, &currentIndex)

	// Context should be short-lived (few mins).
	// TODO make it short-lived
	if err := state.Update(context.TODO(), models.CHECKS, session); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.Update(context.TODO(), models.CURRENT_INDEX, currentIndex); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}
	// set state ShowinChecks
	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BACK,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.FORWARD,
		},
	)
	kb.RemoveKeyboard = true

	return c.Send(util.GetCheckWithItemsResponse(*session.Checks[currentIndex]), kb)
}

func showTotals(c tele.Context, state fsm.Context) error {

	totals, err := db.GetAllSessionTotals()
	if err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	// Context should be short-lived (few mins).
	// TODO make it short-lived
	if err := state.Update(context.TODO(), models.SESSION_TOTALS, totals); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	var currentIndex int
	if err := state.Update(context.TODO(), models.CURRENT_INDEX, currentIndex); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	if err := state.SetState(context.TODO(), models.StateShowingTotals); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BACK,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.FORWARD,
		},
	)

	return c.Send(util.GetShowTotalsResponse(totals[currentIndex]), kb)
}

type checksForSession struct {
	SessionId int64
	Checks    []*models.CheckWithItems
}

// Helper function. that gets current sessionId and pulls all checks for it from db.
// Can return errors only occured during [Bot.Send()]
func getChecksForCurrentSession(c tele.Context) (*checksForSession, error) {
	sessionId, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		var err error
		sessionId, err = db.GetSessionId()
		if err != nil {
			return nil, c.Send(models.ErrorSometingWentWrong)
		}
	}

	checks, err := db.GetAllChecksWithItemsForSesssionId(sessionId)
	if err != nil {
		return nil, c.Send(err)
	}

	return &checksForSession{SessionId: sessionId, Checks: checks}, nil
}
