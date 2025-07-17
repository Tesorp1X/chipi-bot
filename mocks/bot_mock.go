package mocks

import (
	"errors"
	"io"

	tele "gopkg.in/telebot.v4"
)

// MockBotAPI implements the [tele.API] interface with no-op methods and zero values
type MockBotAPI struct {
	// Channel with responses from handler via bot instance or [tele.Context]
	response *HandlerResponse
}

func NewMockBot(resp *HandlerResponse) *MockBotAPI {
	return &MockBotAPI{
		response: resp,
	}
}

// Methods bellow the line are for testing
// --------------------------------------------------------------------

// Methods bellow the line are to satisfy the [tele.API] interface
// --------------------------------------------------------------------
func (m *MockBotAPI) Raw(method string, payload interface{}) ([]byte, error) {
	return nil, nil
}

func (m *MockBotAPI) Accept(query *tele.PreCheckoutQuery, errorMessage ...string) error {
	return nil
}

func (m *MockBotAPI) AddStickerToSet(of tele.Recipient, name string, sticker tele.InputSticker) error {
	return nil
}

func (m *MockBotAPI) AdminsOf(chat *tele.Chat) ([]tele.ChatMember, error) {
	return nil, nil
}

func (m *MockBotAPI) Answer(query *tele.Query, resp *tele.QueryResponse) error {
	return nil
}

func (m *MockBotAPI) AnswerWebApp(query *tele.Query, r tele.Result) (*tele.WebAppMessage, error) {
	return nil, nil
}

func (m *MockBotAPI) ApproveJoinRequest(chat tele.Recipient, user *tele.User) error {
	return nil
}

func (m *MockBotAPI) Ban(chat *tele.Chat, member *tele.ChatMember, revokeMessages ...bool) error {
	return nil
}

func (m *MockBotAPI) BanSenderChat(chat *tele.Chat, sender tele.Recipient) error {
	return nil
}

func (m *MockBotAPI) BusinessConnection(id string) (*tele.BusinessConnection, error) {
	return nil, nil
}

func (m *MockBotAPI) ChatByID(id int64) (*tele.Chat, error) {
	return nil, nil
}

func (m *MockBotAPI) ChatByUsername(name string) (*tele.Chat, error) {
	return nil, nil
}

func (m *MockBotAPI) ChatMemberOf(chat, user tele.Recipient) (*tele.ChatMember, error) {
	return nil, nil
}

func (m *MockBotAPI) Close() (bool, error) {
	return false, nil
}

func (m *MockBotAPI) CloseGeneralTopic(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) CloseTopic(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) Commands(opts ...interface{}) ([]tele.Command, error) {
	return nil, nil
}

func (m *MockBotAPI) Copy(to tele.Recipient, msg tele.Editable, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) CopyMany(to tele.Recipient, msgs []tele.Editable, opts ...*tele.SendOptions) ([]tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) CreateInviteLink(chat tele.Recipient, link *tele.ChatInviteLink) (*tele.ChatInviteLink, error) {
	return nil, nil
}

func (m *MockBotAPI) CreateInvoiceLink(i tele.Invoice) (string, error) {
	return "", nil
}

func (m *MockBotAPI) CreateStickerSet(of tele.Recipient, set *tele.StickerSet) error {
	return nil
}

func (m *MockBotAPI) CreateTopic(chat *tele.Chat, topic *tele.Topic) (*tele.Topic, error) {
	return nil, nil
}

func (m *MockBotAPI) CustomEmojiStickers(ids []string) ([]tele.Sticker, error) {
	return nil, nil
}

func (m *MockBotAPI) DeclineJoinRequest(chat tele.Recipient, user *tele.User) error {
	return nil
}

func (m *MockBotAPI) DefaultRights(forChannels bool) (*tele.Rights, error) {
	return nil, nil
}

func (m *MockBotAPI) Delete(msg tele.Editable) error {
	return nil
}

func (m *MockBotAPI) DeleteCommands(opts ...interface{}) error {
	return nil
}

func (m *MockBotAPI) DeleteGroupPhoto(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) DeleteGroupStickerSet(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) DeleteMany(msgs []tele.Editable) error {
	return nil
}

func (m *MockBotAPI) DeleteSticker(sticker string) error {
	return nil
}

func (m *MockBotAPI) DeleteStickerSet(name string) error {
	return nil
}

func (m *MockBotAPI) DeleteTopic(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) Download(file *tele.File, localFilename string) error {
	return nil
}

func (m *MockBotAPI) Edit(msg tele.Editable, what interface{}, opts ...interface{}) (*tele.Message, error) {
	text, ok := what.(string)
	if !ok {
		return nil, errors.New("expected what of type string")
	}

	m.response.Text = text
	m.response.Type = ResponseTypeEdit
	m.response.SendOptions = extractOptions(opts)

	return nil, nil
}

func (m *MockBotAPI) EditCaption(msg tele.Editable, caption string, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) EditGeneralTopic(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) EditInviteLink(chat tele.Recipient, link *tele.ChatInviteLink) (*tele.ChatInviteLink, error) {
	return nil, nil
}

func (m *MockBotAPI) EditMedia(msg tele.Editable, media tele.Inputtable, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) EditReplyMarkup(msg tele.Editable, markup *tele.ReplyMarkup) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) EditTopic(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) File(file *tele.File) (io.ReadCloser, error) {
	return nil, nil
}

func (m *MockBotAPI) FileByID(fileID string) (tele.File, error) {
	return tele.File{}, nil
}

func (m *MockBotAPI) Forward(to tele.Recipient, msg tele.Editable, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) ForwardMany(to tele.Recipient, msgs []tele.Editable, opts ...*tele.SendOptions) ([]tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) GameScores(user tele.Recipient, msg tele.Editable) ([]tele.GameHighScore, error) {
	return nil, nil
}

func (m *MockBotAPI) HideGeneralTopic(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) InviteLink(chat *tele.Chat) (string, error) {
	return "", nil
}

func (m *MockBotAPI) Leave(chat tele.Recipient) error {
	return nil
}

func (m *MockBotAPI) Len(chat *tele.Chat) (int, error) {
	return 0, nil
}

func (m *MockBotAPI) Logout() (bool, error) {
	return false, nil
}

func (m *MockBotAPI) MenuButton(chat *tele.User) (*tele.MenuButton, error) {
	return nil, nil
}

func (m *MockBotAPI) MyDescription(language string) (*tele.BotInfo, error) {
	return nil, nil
}

func (m *MockBotAPI) MyName(language string) (*tele.BotInfo, error) {
	return nil, nil
}

func (m *MockBotAPI) MyShortDescription(language string) (*tele.BotInfo, error) {
	return nil, nil
}

func (m *MockBotAPI) Notify(to tele.Recipient, action tele.ChatAction, threadID ...int) error {
	return nil
}

func (m *MockBotAPI) Pin(msg tele.Editable, opts ...interface{}) error {
	return nil
}

func (m *MockBotAPI) ProfilePhotosOf(user *tele.User) ([]tele.Photo, error) {
	return nil, nil
}

func (m *MockBotAPI) Promote(chat *tele.Chat, member *tele.ChatMember) error {
	return nil
}

func (m *MockBotAPI) React(to tele.Recipient, msg tele.Editable, r tele.Reactions) error {
	return nil
}

func (m *MockBotAPI) RefundStars(to tele.Recipient, chargeID string) error {
	return nil
}

func (m *MockBotAPI) RemoveWebhook(dropPending ...bool) error {
	return nil
}

func (m *MockBotAPI) ReopenGeneralTopic(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) ReopenTopic(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) ReplaceStickerInSet(of tele.Recipient, stickerSet, oldSticker string, sticker tele.InputSticker) (bool, error) {
	return false, nil
}

func (m *MockBotAPI) Reply(to *tele.Message, what interface{}, opts ...interface{}) (*tele.Message, error) {
	text, ok := what.(string)
	if !ok {
		return nil, errors.New("expected what of type string")
	}

	m.response.Text = text
	m.response.Type = ResponseTypeReply
	m.response.SendOptions = extractOptions(opts)

	return nil, nil
}

func (m *MockBotAPI) Respond(c *tele.Callback, resp ...*tele.CallbackResponse) error {
	return nil
}

func (m *MockBotAPI) Restrict(chat *tele.Chat, member *tele.ChatMember) error {
	return nil
}

func (m *MockBotAPI) RevokeInviteLink(chat tele.Recipient, link string) (*tele.ChatInviteLink, error) {
	return nil, nil
}

func (m *MockBotAPI) Send(to tele.Recipient, what interface{}, opts ...interface{}) (*tele.Message, error) {
	text, ok := what.(string)
	if !ok {
		return nil, errors.New("expected what of type string")
	}

	m.response.Text = text
	m.response.Type = ResponseTypeSend
	m.response.SendOptions = extractOptions(opts)

	return nil, nil
}

func (m *MockBotAPI) SendAlbum(to tele.Recipient, a tele.Album, opts ...interface{}) ([]tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) SendPaid(to tele.Recipient, stars int, a tele.PaidAlbum, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) SetAdminTitle(chat *tele.Chat, user *tele.User, title string) error {
	return nil
}

func (m *MockBotAPI) SetCommands(opts ...interface{}) error {
	return nil
}

func (m *MockBotAPI) SetCustomEmojiStickerSetThumb(name, id string) error {
	return nil
}

func (m *MockBotAPI) SetDefaultRights(rights tele.Rights, forChannels bool) error {
	return nil
}

func (m *MockBotAPI) SetGameScore(user tele.Recipient, msg tele.Editable, score tele.GameHighScore) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) SetGroupDescription(chat *tele.Chat, description string) error {
	return nil
}

func (m *MockBotAPI) SetGroupPermissions(chat *tele.Chat, perms tele.Rights) error {
	return nil
}

func (m *MockBotAPI) SetGroupStickerSet(chat *tele.Chat, setName string) error {
	return nil
}

func (m *MockBotAPI) SetGroupTitle(chat *tele.Chat, title string) error {
	return nil
}

func (m *MockBotAPI) SetMenuButton(chat *tele.User, mb interface{}) error {
	return nil
}

func (m *MockBotAPI) SetMyDescription(desc, language string) error {
	return nil
}

func (m *MockBotAPI) SetMyName(name, language string) error {
	return nil
}

func (m *MockBotAPI) SetMyShortDescription(desc, language string) error {
	return nil
}

func (m *MockBotAPI) SetStickerEmojis(sticker string, emojis []string) error {
	return nil
}

func (m *MockBotAPI) SetStickerKeywords(sticker string, keywords []string) error {
	return nil
}

func (m *MockBotAPI) SetStickerMaskPosition(sticker string, mask tele.MaskPosition) error {
	return nil
}

func (m *MockBotAPI) SetStickerPosition(sticker string, position int) error {
	return nil
}

func (m *MockBotAPI) SetStickerSetThumb(of tele.Recipient, set *tele.StickerSet) error {
	return nil
}

func (m *MockBotAPI) SetStickerSetTitle(s tele.StickerSet) error {
	return nil
}

func (m *MockBotAPI) SetWebhook(w *tele.Webhook) error {
	return nil
}

func (m *MockBotAPI) Ship(query *tele.ShippingQuery, what ...interface{}) error {
	return nil
}

func (m *MockBotAPI) StarTransactions(offset, limit int) ([]tele.StarTransaction, error) {
	return nil, nil
}

func (m *MockBotAPI) StickerSet(name string) (*tele.StickerSet, error) {
	return nil, nil
}

func (m *MockBotAPI) StopLiveLocation(msg tele.Editable, opts ...interface{}) (*tele.Message, error) {
	return nil, nil
}

func (m *MockBotAPI) StopPoll(msg tele.Editable, opts ...interface{}) (*tele.Poll, error) {
	return nil, nil
}

func (m *MockBotAPI) TopicIconStickers() ([]tele.Sticker, error) {
	return nil, nil
}

func (m *MockBotAPI) Unban(chat *tele.Chat, user *tele.User, forBanned ...bool) error {
	return nil
}

func (m *MockBotAPI) UnbanSenderChat(chat *tele.Chat, sender tele.Recipient) error {
	return nil
}

func (m *MockBotAPI) UnhideGeneralTopic(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) Unpin(chat tele.Recipient, messageID ...int) error {
	return nil
}

func (m *MockBotAPI) UnpinAll(chat tele.Recipient) error {
	return nil
}

func (m *MockBotAPI) UnpinAllGeneralTopicMessages(chat *tele.Chat) error {
	return nil
}

func (m *MockBotAPI) UnpinAllTopicMessages(chat *tele.Chat, topic *tele.Topic) error {
	return nil
}

func (m *MockBotAPI) UploadSticker(to tele.Recipient, format tele.StickerSetFormat, f tele.File) (*tele.File, error) {
	return nil, nil
}

func (m *MockBotAPI) UserBoosts(chat, user tele.Recipient) ([]tele.Boost, error) {
	return nil, nil
}

func (m *MockBotAPI) Webhook() (*tele.Webhook, error) {
	return nil, nil
}
