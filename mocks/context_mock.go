package mocks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	tele "gopkg.in/telebot.v4"
)

type MockStorage struct {
	s map[string]any
	m *sync.Mutex
}

func (s *MockStorage) Set(key string, val any) {
	s.m.Lock()
	s.s[key] = val
	s.m.Unlock()
}

func (s *MockStorage) Get(key string) any {
	val, ok := s.s[key]
	if !ok {
		return nil
	}

	return val
}

// what was sent to a user
type HandlerResponse struct {

	// Text of a displayed message or Text field of [tele.tele.CallbackResponse].
	Text string
	// In which way message was sent (Send, Reply, Edit, EditOrReply, Respond).
	// Supported options are defined as iota-constants.
	Type int
	// Which [SendOptions] were used with response.
	SendOptions *tele.SendOptions
}

const (
	// c is tele.Context
	// if handler called c.Send()
	ResponseTypeSend = iota
	// if handler called c.Reply()
	ResponseTypeReply
	// if handler called c.Edit()
	ResponseTypeEdit
	// if handler called c.EditOrReply()
	ResponseTypeEditOrReply
	// if handler called c.EditOrSend()
	ResponseTypeEditOrSend
	// if handler responded to a callback query
	ResponseTypeCallbackResponse
	// if handler responded to a callback query with alert
	ResponseTypeCallbackResponseWithAlert
)

func copyReplyMarkUp(r *tele.ReplyMarkup) *tele.ReplyMarkup {
	cp := *r

	if len(r.ReplyKeyboard) > 0 {
		cp.ReplyKeyboard = make([][]tele.ReplyButton, len(r.ReplyKeyboard))
		for i, row := range r.ReplyKeyboard {
			cp.ReplyKeyboard[i] = make([]tele.ReplyButton, len(row))
			copy(cp.ReplyKeyboard[i], row)
		}
	}

	if len(r.InlineKeyboard) > 0 {
		cp.InlineKeyboard = make([][]tele.InlineButton, len(r.InlineKeyboard))
		for i, row := range r.InlineKeyboard {
			cp.InlineKeyboard[i] = make([]tele.InlineButton, len(row))
			copy(cp.InlineKeyboard[i], row)
		}
	}

	return &cp
}

func copySendOptions(og *tele.SendOptions) *tele.SendOptions {
	cp := *og
	if cp.ReplyMarkup != nil {
		cp.ReplyMarkup = copyReplyMarkUp(cp.ReplyMarkup)
	}
	return &cp
}

func extractOptions(how []interface{}) *tele.SendOptions {
	opts := &tele.SendOptions{
		ParseMode: tele.ModeHTML,
	}

	for _, prop := range how {
		switch opt := prop.(type) {
		case *tele.SendOptions:
			opts = copySendOptions(opt)
		case *tele.ReplyMarkup:
			if opt != nil {
				opts.ReplyMarkup = copyReplyMarkUp(opt)
			}
		case *tele.ReplyParams:
			opts.ReplyParams = opt
		case *tele.Topic:
			opts.ThreadID = opt.ThreadID
		case tele.Option:
			switch opt {
			case tele.NoPreview:
				opts.DisableWebPagePreview = true
			case tele.Silent:
				opts.DisableNotification = true
			case tele.AllowWithoutReply:
				opts.AllowWithoutReply = true
			case tele.ForceReply:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &tele.ReplyMarkup{}
				}
				opts.ReplyMarkup.ForceReply = true
			case tele.OneTimeKeyboard:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &tele.ReplyMarkup{}
				}
				opts.ReplyMarkup.OneTimeKeyboard = true
			case tele.RemoveKeyboard:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &tele.ReplyMarkup{}
				}
				opts.ReplyMarkup.RemoveKeyboard = true
			case tele.Protected:
				opts.Protected = true
			default:
				panic("telebot: unsupported flag-option")
			}
		case tele.ParseMode:
			opts.ParseMode = opt
		case tele.Entities:
			opts.Entities = opt
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}

// MockContext mocks the original tele.Context for testing purpouses.
type MockContext struct {
	bot     tele.API
	update  tele.Update
	storage *MockStorage

	response *HandlerResponse
}

func NewMockContext(bot tele.API, update tele.Update, storage *MockStorage) *MockContext {
	return &MockContext{
		bot:     bot,
		update:  update,
		storage: storage,
	}
}

// Methods bellow the line are for testing
// --------------------------------------------------------------------

func (c *MockContext) IsResponseTypeEqualsTo(responseType int) bool {
	return responseType == c.response.Type
}

func (c *MockContext) IsResponseTextEqualsTo(text string) bool {
	return text == c.response.Text
}

// Returns true if [ReplyMarkup] stored in context.response is the same to given kb.
func (c *MockContext) IsResponseReplyMarkUpEqualsTo(kb *tele.ReplyMarkup) (bool, error) {

	ogKbJson, err := json.Marshal(c.response.SendOptions.ReplyMarkup)
	if err != nil {
		return false, fmt.Errorf("IsResponseReplyMarkUpEqualsTo: couldn't marshal original ReplyMarkup: %w", err)
	}

	givenKbJson, err := json.Marshal(kb)
	if err != nil {
		return false, fmt.Errorf("IsResponseReplyMarkUpEqualsTo: couldn't marshal given ReplyMarkup: %w", err)
	}

	if !bytes.Equal(ogKbJson, givenKbJson) {
		return false, nil
	}

	return true, nil
}

// Methods bellow the line are to satisfy the tele.Context interface
// --------------------------------------------------------------------

// Bot returns the bot instance.
func (c *MockContext) Bot() tele.API {
	return c.bot
}

// // Update returns the original update.
func (c *MockContext) Update() tele.Update {
	return c.update
}

// // Message returns stored message if such presented.
func (c *MockContext) Message() *tele.Message {
	if c.update.Message != nil {
		return c.update.Message
	}
	return nil
}

// // Callback returns stored callback if such presented.
func (c *MockContext) Callback() *tele.Callback {
	if c.update.Callback != nil {
		return c.update.Callback
	}
	return nil
}

// // Query returns stored query if such presented.
func (c *MockContext) Query() *tele.Query {
	return &tele.Query{}
}

// // InlineResult returns stored inline result if such presented.
func (c *MockContext) InlineResult() *tele.InlineResult {
	return &tele.InlineResult{}
}

// // ShippingQuery returns stored shipping query if such presented.
func (c *MockContext) ShippingQuery() *tele.ShippingQuery {
	return &tele.ShippingQuery{}
}

// // PreCheckoutQuery returns stored pre checkout query if such presented.
func (c *MockContext) PreCheckoutQuery() *tele.PreCheckoutQuery {
	return &tele.PreCheckoutQuery{}
}

// // Payment returns payment instance.
func (c *MockContext) Payment() *tele.Payment {
	return &tele.Payment{}
}

// // Poll returns stored poll if such presented.
func (c *MockContext) Poll() *tele.Poll {
	return &tele.Poll{}
}

// // PollAnswer returns stored poll answer if such presented.
func PollAnswer() *tele.PollAnswer {
	return &tele.PollAnswer{}
}

// // ChatMember returns chat member changes.
func (c *MockContext) ChatMember() *tele.ChatMemberUpdate {
	return &tele.ChatMemberUpdate{}
}

// // ChatJoinRequest returns the chat join request.
func (c *MockContext) ChatJoinRequest() *tele.ChatJoinRequest {
	return &tele.ChatJoinRequest{}
}

// // Migration returns both migration from and to chat IDs.
func (c *MockContext) Migration() (int64, int64) {
	return 0, 0
}

// // Topic returns the topic changes.
func (c *MockContext) Topic() *tele.Topic {
	return &tele.Topic{}
}

// // Boost returns the boost instance.
func (c *MockContext) Boost() *tele.BoostUpdated {
	return &tele.BoostUpdated{}
}

// // BoostRemoved returns the boost removed from a chat instance.
func (c *MockContext) BoostRemoved() *tele.BoostRemoved {
	return &tele.BoostRemoved{}
}

// // Sender returns the current recipient, depending on the context type.
// // Returns nil if user is not presented.
func (c *MockContext) Sender() *tele.User {
	return &tele.User{}
}

// // Chat returns the current chat, depending on the context type.
// // Returns nil if chat is not presented.
func (c *MockContext) Chat() *tele.Chat {
	return &tele.Chat{}
}

// // Recipient combines both Sender and Chat functions. If there is no user
// // the chat will be returned. The native context cannot be without sender,
// // but it is useful in the case when the context created intentionally
// // by the NewContext constructor and have only Chat field inside.
func (c *MockContext) Recipient() tele.Recipient {
	return &tele.User{}
}

// // Text returns the message text, depending on the context type.
// // In the case when no related data presented, returns an empty string.
func (c *MockContext) Text() string {
	if c.update.Message != nil {
		return c.update.Message.Text
	}
	return ""
}

// // Entities returns the message entities, whether it's media caption's or the text's.
// // In the case when no entities presented, returns a nil.
func (c *MockContext) Entities() tele.Entities {
	return tele.Entities{}
}

// // Data returns the current data, depending on the context type.
// // If the context contains command, returns its arguments string.
// // If the context contains payment, returns its payload.
// // In the case when no related data presented, returns an empty string.
func (c *MockContext) Data() string {
	return ""
}

// // Args returns a raw slice of command or callback arguments as strings.
// // The message arguments split by space, while the callback's ones by a "|" symbol.
func (c *MockContext) Args() []string {
	return []string{}
}

// // Send sends a message to the current recipient.
// // See Send from bot.go.
func (c *MockContext) Send(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeSend
	c.response.SendOptions = extractOptions(opts)

	return nil
}

// // SendAlbum sends an album to the current recipient.
// // See SendAlbum from bot.go.
func (c *MockContext) SendAlbum(a tele.Album, opts ...interface{}) error {
	return nil
}

// // Reply replies to the current message.
// // See Reply from bot.go.
func (c *MockContext) Reply(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeReply
	c.response.SendOptions = extractOptions(opts)

	return nil
}

// // Forward forwards the given message to the current recipient.
// // See Forward from bot.go.
func (c *MockContext) Forward(msg tele.Editable, opts ...interface{}) error {
	return nil
}

// // ForwardTo forwards the current message to the given recipient.
// // See Forward from bot.go
func (c *MockContext) ForwardTo(to tele.Recipient, opts ...interface{}) error {
	return nil
}

// // Edit edits the current message.
// // See Edit from bot.go.
func (c *MockContext) Edit(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeEdit
	c.response.SendOptions = extractOptions(opts)

	return nil
}

// // EditCaption edits the caption of the current message.
// // See EditCaption from bot.go.
func (c *MockContext) EditCaption(caption string, opts ...interface{}) error {
	return nil
}

// // EditOrSend edits the current message if the update is callback,
// // otherwise the content is sent to the chat as a separate message.
func (c *MockContext) EditOrSend(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeEditOrSend
	c.response.SendOptions = extractOptions(opts)

	return nil
}

// // EditOrReply edits the current message if the update is callback,
// // otherwise the content is replied as a separate message.
func (c *MockContext) EditOrReply(what interface{}, opts ...interface{}) error {
	text, ok := what.(string)
	if !ok {
		return errors.New("expected what of type string")
	}

	c.response.Text = text
	c.response.Type = ResponseTypeEditOrReply
	c.response.SendOptions = extractOptions(opts)

	return nil
}

// // Delete removes the current message.
// // See Delete from bot.go.
func (c *MockContext) Delete() error {
	return nil
}

// // DeleteAfter waits for the duration to elapse and then removes the
// // message. It handles an error automatically using b.OnError callback.
// // It returns a Timer that can be used to cancel the call using its Stop method.
func (c *MockContext) DeleteAfter(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}

// // Notify updates the chat action for the current recipient.
// // See Notify from bot.go.
func (c *MockContext) Notify(action tele.ChatAction) error {
	return nil
}

// // Ship replies to the current shipping query.
// // See Ship from bot.go.
func (c *MockContext) Ship(what ...interface{}) error {
	return nil
}

// // Accept finalizes the current deal.
// // See Accept from bot.go.
func (c *MockContext) Accept(errorMessage ...string) error {
	return nil
}

// // Answer sends a response to the current inline query.
// // See Answer from bot.go.
func (c *MockContext) Answer(resp *tele.QueryResponse) error {
	return nil
}

// // Respond sends a response for the current callback query.
// // See Respond from bot.go.
func (c *MockContext) Respond(resp ...*tele.CallbackResponse) error {
	c.response.Text = resp[0].Text
	if resp[0].ShowAlert {
		c.response.Type = ResponseTypeCallbackResponseWithAlert
	} else {
		c.response.Type = ResponseTypeCallbackResponse
	}

	return nil
}

// // RespondText sends a popup response for the current callback query.
func (c *MockContext) RespondText(text string) error {
	return nil
}

// // RespondAlert sends an alert response for the current callback query.
func (c *MockContext) RespondAlert(text string) error {
	return nil
}

// // Get retrieves data from the context.
func (c *MockContext) Get(key string) interface{} {
	return c.storage.Get(key)
}

// // Set saves data in the context.
func (c *MockContext) Set(key string, val interface{}) {
	c.storage.Set(key, val)
}
