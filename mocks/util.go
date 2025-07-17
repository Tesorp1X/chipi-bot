package mocks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	tele "gopkg.in/telebot.v4"
)

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

func (r *HandlerResponse) IsResponseTypeEqualsTo(responseType int) bool {
	return responseType == r.Type
}

func (r *HandlerResponse) IsResponseTextEqualsTo(text string) bool {
	return text == r.Text
}

// Returns true if [ReplyMarkup] stored in context.response is the same to given kb.
func (r *HandlerResponse) IsResponseReplyMarkUpEqualsTo(kb *tele.ReplyMarkup) (bool, error) {

	if r.SendOptions == nil {
		return false, errors.New("SendOption is nil")
	}

	ogKbJson, err := json.Marshal(r.SendOptions.ReplyMarkup)
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

type MockStorage struct {
	s map[string]any
	m *sync.RWMutex
}

func NewMockStorage() *MockStorage {
	ms := new(MockStorage)
	ms.s = make(map[string]any)
	ms.m = new(sync.RWMutex)

	return ms
}

func (s *MockStorage) Set(key string, val any) {
	s.m.Lock()
	s.s[key] = val
	s.m.Unlock()
}

func (s *MockStorage) Get(key string) any {
	s.m.RLock()
	val, ok := s.s[key]
	s.m.RUnlock()
	if !ok {
		return nil
	}

	return val
}

// Clears current storage.
func (s *MockStorage) ClearData() {
	s.m.Lock()
	clear(s.s)
	s.m.Unlock()
}

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
