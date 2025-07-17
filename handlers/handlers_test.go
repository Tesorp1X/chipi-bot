package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/Tesorp1X/chipi-bot/mocks"
	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

var errEmpty = errors.New("")

func assertHandlerError(t testing.TB, wantErr bool, expexted, got error) {
	t.Helper()
	if !wantErr && got == nil {
		return
	}

	if !wantErr && got != nil {
		t.Fatalf("didn't expect an error but got: %v", got)
	}

	if wantErr && got == nil {
		t.Fatalf("expected an error: '%v' but got none", expexted)
	}

	if expexted != got {
		t.Fatalf("expected an error: '%v'', but instead got: '%v'", expexted, got)
	}
}

func assertHandlerResponse(t testing.TB, expected, got *mocks.HandlerResponse) {
	t.Helper()
	if !got.IsResponseTextEqualsTo(expected.Text) {
		t.Fatalf("expected response text: '%s', but got: '%s'", expected.Text, got.Text)
	}

	if !got.IsResponseTypeEqualsTo(expected.Type) {
		t.Fatalf("expected response type: '%d', but got: '%d'", expected.Type, got.Type)
	}

	if expected.SendOptions != nil {
		if ok, err := got.IsResponseReplyMarkUpEqualsTo(expected.SendOptions.ReplyMarkup); !ok {
			if err != nil {
				t.Fatalf("assertResponse error: %v", err)
			}
			t.Fatal("wrong reply markup")
		}
	}

}

// Fails a test if fsmCtx's state doesn't equal to expected.
func assertState(t testing.TB, expected fsm.State, fsmCtx *mocks.MockFsmContext) {
	t.Helper()
	gotState, err := fsmCtx.State(context.Background())
	if err != nil {
		t.Fatalf("assertState error: %v", err)
	}

	if expected != gotState {
		t.Fatalf("expected state: %s, but got insted: %s", expected, gotState)
	}
}

// Returns an update with non-nil Message field. User has ID 1, Name: Test Test and username: @test123.
func makeUpdateWithMessageText(text string) tele.Update {
	return tele.Update{
		ID: 1,
		Message: &tele.Message{
			ID:   1,
			Text: text,
			Sender: &tele.User{
				ID:        1,
				FirstName: "Test",
				LastName:  "Test",
				Username:  "test123",
			},
		},
	}

}

func TestHelloHandler(t *testing.T) {
	response := mocks.HandlerResponse{}
	bot := mocks.NewMockBot(&response)
	botStorage := mocks.NewMockStorage()
	fsmStorage := mocks.NewMockStorage()

	update := makeUpdateWithMessageText("hello")

	teleCtx := mocks.NewMockContext(bot, update, botStorage, &response)
	stateCtx := mocks.NewMockFsmContext(fsmStorage, models.StateDefault)

	expextedResponse := mocks.HandlerResponse{
		Text: "Hello, 1",
		Type: mocks.ResponseTypeSend,
	}

	if err := stateCtx.SetState(context.Background(), models.StateStart); err != nil {
		t.Fatalf("couldn't change state to %s: %v", models.StateStart, err)
	}

	handlerErr := HelloHandler(teleCtx, stateCtx)

	assertHandlerError(t, false, errEmpty, handlerErr)
	assertHandlerResponse(t, &expextedResponse, &response)
}

func TestNewCheckHandler(t *testing.T) {
	t.Run("command-call from default state (no err)", func(t *testing.T) {
		response := mocks.HandlerResponse{}
		bot := mocks.NewMockBot(&response)
		botStorage := mocks.NewMockStorage()
		var sessionId int64 = 1
		botStorage.Set(models.SESSION_ID, sessionId)
		fsmStorage := mocks.NewMockStorage()

		update := makeUpdateWithMessageText("/newcheck")

		teleCtx := mocks.NewMockContext(bot, update, botStorage, &response)
		stateCtx := mocks.NewMockFsmContext(fsmStorage, models.StateDefault)

		expectedResponse := &mocks.HandlerResponse{
			Text: "Хорошо, как назовем новый чек?👀",
			Type: mocks.ResponseTypeSend,
		}

		handlerErr := NewCheckHandler(teleCtx, stateCtx)

		assertHandlerError(t, false, errEmpty, handlerErr)
		assertHandlerResponse(t, expectedResponse, &response)
		assertState(t, models.StateWaitForCheckName, stateCtx)
	})
}
