package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/Tesorp1X/chipi-bot/mocks"
	"github.com/Tesorp1X/chipi-bot/models"
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
	storage := mocks.NewMockStorage()

	update := makeUpdateWithMessageText("hello")

	teleCtx := mocks.NewMockContext(bot, update, storage, &response)
	stateCtx := mocks.NewMockFsmContext(storage, models.StateDefault)

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
