package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// NewMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int, text string) SendMessageConf {
	return SendMessageConf{
		BaseSend: BaseSend{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		DisableWebPagePreview: false,
	}
}

// NewMessageToChannel creates a new Message that is sent to a channel
// by username.
//
// username is the username of the channel, text is the message text,
// and the username should be in the form of `@username`.
func NewMessageToChannel(username string, text string) SendMessageConf {
	return SendMessageConf{
		BaseSend: BaseSend{
			ChatID:           username,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		DisableWebPagePreview: false,
	}
}

// NewDeleteMessage creates a request to delete a message.
func NewDeleteMessage(chatID int, messageID int) DeleteMessageConf {
	return DeleteMessageConf{
		ChatID:    chatID,
		MessageID: messageID,
	}
}

// NewForward creates a new forward.
//
// chatID is where to send it, fromChatID is the source chat,
// and messageID is the ID of the original message.
func NewForward(chatID int, fromChatID int, messageID int) ForwardMessageConf {
	return ForwardMessageConf{
		ChatID:     chatID,
		FromChatID: fromChatID,
		MessageID:  messageID,
	}
}

// NewCopyMessage creates a new copy message.
//
// chatID is where to send it, fromChatID is the source chat,
// and messageID is the ID of the original message.
func NewCopyMessage(chatID int, fromChatID int, messageID int) CopyMessageConf {
	return CopyMessageConf{
		BaseSend:   BaseSend{ChatID: chatID},
		FromChatID: fromChatID,
		MessageID:  messageID,
	}
}

// NewPhoto creates a new sendPhoto request.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
//
// Note that you must send animated GIFs as a document.
func NewPhoto(chatID int, file RequestFileData) SendPhotoConf {
	return SendPhotoConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewPhotoToChannel creates a new photo uploader to send a photo to a channel.
//
// Note that you must send animated GIFs as a document.
func NewPhotoToChannel(username string, file RequestFileData) SendPhotoConf {
	return SendPhotoConf{
		BaseSend: BaseSend{ChatID: username},
		File:     file,
	}
}

// NewAudio creates a new sendAudio request.
func NewAudio(chatID int, file RequestFileData) SendAudioConf {
	return SendAudioConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewDocument creates a new sendDocument request.
func NewDocument(chatID int, file RequestFileData) SendDocumentConf {
	return SendDocumentConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewSticker creates a new sendSticker request.
func NewSticker(chatID int, file RequestFileData) SendStickerConf {
	return SendStickerConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewVideo creates a new sendVideo request.
func NewVideo(chatID int, file RequestFileData) SendVideoConf {
	return SendVideoConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewAnimation creates a new sendAnimation request.
func NewAnimation(chatID int, file RequestFileData) SendAnimationConf {
	return SendAnimationConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewVideoNote creates a new sendVideoNote request.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewVideoNote(chatID int, length int, file RequestFileData) SendVideoNoteConf {
	return SendVideoNoteConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
		Length:   length,
	}
}

// NewVoice creates a new sendVoice request.
func NewVoice(chatID int, file RequestFileData) SendVoiceConf {
	return SendVoiceConf{
		BaseSend: BaseSend{ChatID: chatID},
		File:     file,
	}
}

// NewMediaGroup creates a new media group. Files should be an array of
// two to ten InputMediaPhoto or InputMediaVideo.
func NewMediaGroup(chatID int, files []interface{}) SendMediaGroupConf {
	return SendMediaGroupConf{
		ChatID: chatID,
		Media:  files,
	}
}

// NewInputMediaPhoto creates a new InputMediaPhoto.
func NewInputMediaPhoto(media RequestFileData) InputMediaPhoto {
	return InputMediaPhoto{
		InputMediaBase: InputMediaBase{
			Type:  "photo",
			Media: media,
		},
	}
}

// NewInputMediaVideo creates a new InputMediaVideo.
func NewInputMediaVideo(media RequestFileData) InputMediaVideo {
	return InputMediaVideo{
		InputMediaBase: InputMediaBase{
			Type:  "video",
			Media: media,
		},
	}
}

// NewInputMediaAnimation creates a new InputMediaAnimation.
func NewInputMediaAnimation(media RequestFileData) InputMediaAnimation {
	return InputMediaAnimation{
		InputMediaBase: InputMediaBase{
			Type:  "animation",
			Media: media,
		},
	}
}

// NewInputMediaAudio creates a new InputMediaAudio.
func NewInputMediaAudio(media RequestFileData) InputMediaAudio {
	return InputMediaAudio{
		InputMediaBase: InputMediaBase{
			Type:  "audio",
			Media: media,
		},
	}
}

// NewInputMediaDocument creates a new InputMediaDocument.
func NewInputMediaDocument(media RequestFileData) InputMediaDocument {
	return InputMediaDocument{
		InputMediaBase: InputMediaBase{
			Type:  "document",
			Media: media,
		},
	}
}

// NewContact allows you to send a shared contact.
func NewContact(chatID int, phoneNumber, firstName string) SendContactConf {
	return SendContactConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
}

// NewLocation shares your location.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int, latitude float64, longitude float64) SendLocationConf {
	return SendLocationConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewVenue allows you to send a venue and its location.
func NewVenue(chatID int, title, address string, latitude, longitude float64) SendVenueConf {
	return SendVenueConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
		Title:     title,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewChatAction sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatID is where to send it, action should be set via Chat constants.
func NewChatAction(chatID int, action string) SendChatActionConf {
	return SendChatActionConf{
		ChatID: chatID,
		Action: action,
	}
}

// NewUserProfilePhotos gets user profile photos.
//
// userID is the ID of the user you wish to get profile photos from.
func NewUserProfilePhotos(userID int) GetUserProfilePhotosConf {
	return GetUserProfilePhotosConf{
		UserID: userID,
		Offset: 0,
		Limit:  0,
	}
}

// NewUpdate gets updates since the last Offset.
//
// offset is the last Update ID to include.
// You likely want to set this to the last Update ID plus 1.
func NewUpdate(offset int) GetUpdatesConf {
	return GetUpdatesConf{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}

// NewWebhook creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) (SetWebhookConf, error) {
	u, err := url.Parse(link)

	if err != nil {
		return SetWebhookConf{}, err
	}

	return SetWebhookConf{
		URL: u,
	}, nil
}

// NewWebhookWithCert creates a new webhook with a certificate.
//
// link is the url you wish to get webhooks,
// file contains a string to a file, FileReader, or FileBytes.
func NewWebhookWithCert(link string, file RequestFileData) (SetWebhookConf, error) {
	u, err := url.Parse(link)

	if err != nil {
		return SetWebhookConf{}, err
	}

	return SetWebhookConf{
		URL:         u,
		Certificate: file,
	}, nil
}

// NewInlineQueryResultArticle creates a new inline query article.
func NewInlineQueryResultArticle(id, title, messageText string) InlineQueryResultArticle {
	return InlineQueryResultArticle{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "article",
			ID:   id,
		},
		Title: title,
		InputMessageContent: InputTextMessageContent{
			MessageText: messageText,
		},
	}
}

// NewInlineQueryResultArticleMarkdown creates a new inline query article with Markdown parsing.
func NewInlineQueryResultArticleMarkdown(id, title, messageText string) InlineQueryResultArticle {
	return InlineQueryResultArticle{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "article",
			ID:   id,
		},
		Title: title,
		InputMessageContent: InputTextMessageContent{
			MessageText: messageText,
			ParseMode:   "Markdown",
		},
	}
}

// NewInlineQueryResultArticleMarkdownV2 creates a new inline query article with MarkdownV2 parsing.
func NewInlineQueryResultArticleMarkdownV2(id, title, messageText string) InlineQueryResultArticle {
	return InlineQueryResultArticle{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "article",
			ID:   id,
		},
		Title: title,
		InputMessageContent: InputTextMessageContent{
			MessageText: messageText,
			ParseMode:   "MarkdownV2",
		},
	}
}

// NewInlineQueryResultArticleHTML creates a new inline query article with HTML parsing.
func NewInlineQueryResultArticleHTML(id, title, messageText string) InlineQueryResultArticle {
	return InlineQueryResultArticle{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "article",
			ID:   id,
		},
		Title: title,
		InputMessageContent: InputTextMessageContent{
			MessageText: messageText,
			ParseMode:   "HTML",
		},
	}
}

// NewInlineQueryResultGIF creates a new inline query GIF.
func NewInlineQueryResultGIF(id, url string) InlineQueryResultGif {
	return InlineQueryResultGif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "gif",
			ID:   id,
		},
		URL: url,
	}
}

// NewInlineQueryResultCachedGIF create a new inline query with cached photo.
func NewInlineQueryResultCachedGIF(id, gifID string) InlineQueryResultCachedGif {
	return InlineQueryResultCachedGif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "gif",
			ID:   id,
		},
		GifFileID: gifID,
	}
}

// NewInlineQueryResultMPEG4GIF creates a new inline query MPEG4 GIF.
func NewInlineQueryResultMPEG4GIF(id, url string) InlineQueryResultMpeg4Gif {
	return InlineQueryResultMpeg4Gif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "mpeg4_gif",
			ID:   id,
		},
		URL: url,
	}
}

// NewInlineQueryResultCachedMPEG4GIF create a new inline query with cached MPEG4 GIF.
func NewInlineQueryResultCachedMPEG4GIF(id, Mpeg4FileID string) InlineQueryResultCachedMpeg4Gif {
	return InlineQueryResultCachedMpeg4Gif{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "mpeg4_gif",
			ID:   id,
		},
		Mpeg4FileID: Mpeg4FileID,
	}
}

// NewInlineQueryResultPhoto creates a new inline query photo.
func NewInlineQueryResultPhoto(id, url string) InlineQueryResultPhoto {
	return InlineQueryResultPhoto{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "photo",
			ID:   id,
		},
		URL: url,
	}
}

// NewInlineQueryResultPhotoWithThumb creates a new inline query photo.
func NewInlineQueryResultPhotoWithThumb(id, url, thumb string) InlineQueryResultPhoto {
	return InlineQueryResultPhoto{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "photo",
			ID:   id,
		},
		URL:          url,
		ThumbnailURL: thumb,
	}
}

// NewInlineQueryResultCachedPhoto create a new inline query with cached photo.
func NewInlineQueryResultCachedPhoto(id, photoID string) InlineQueryResultCachedPhoto {
	return InlineQueryResultCachedPhoto{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "photo",
			ID:   id,
		},
		PhotoFileID: photoID,
	}
}

// NewInlineQueryResultVideo creates a new inline query video.
func NewInlineQueryResultVideo(id, url string) InlineQueryResultVideo {
	return InlineQueryResultVideo{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "video",
			ID:   id,
		},
		URL: url,
	}
}

// NewInlineQueryResultCachedVideo create a new inline query with cached video.
func NewInlineQueryResultCachedVideo(id, videoID, title string) InlineQueryResultCachedVideo {
	return InlineQueryResultCachedVideo{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "video",
			ID:   id,
		},
		VideoFileID: videoID,
		Title:       title,
	}
}

// NewInlineQueryResultCachedSticker create a new inline query with cached sticker.
func NewInlineQueryResultCachedSticker(id, stickerID string) InlineQueryResultCachedSticker {
	return InlineQueryResultCachedSticker{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "sticker",
			ID:   id,
		},
		StickerFileID: stickerID,
	}
}

// NewInlineQueryResultAudio creates a new inline query audio.
func NewInlineQueryResultAudio(id, url, title string) InlineQueryResultAudio {
	return InlineQueryResultAudio{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "audio",
			ID:   id,
		},
		URL:   url,
		Title: title,
	}
}

// NewInlineQueryResultCachedAudio create a new inline query with cached photo.
func NewInlineQueryResultCachedAudio(id, audioID string) InlineQueryResultCachedAudio {
	return InlineQueryResultCachedAudio{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "audio",
			ID:   id,
		},
		AudioFileID: audioID,
	}
}

// NewInlineQueryResultVoice creates a new inline query voice.
func NewInlineQueryResultVoice(id, url, title string) InlineQueryResultVoice {
	return InlineQueryResultVoice{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "voice",
			ID:   id,
		},
		URL:   url,
		Title: title,
	}
}

// NewInlineQueryResultCachedVoice create a new inline query with cached photo.
func NewInlineQueryResultCachedVoice(id, voiceID, title string) InlineQueryResultCachedVoice {
	return InlineQueryResultCachedVoice{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "voice",
			ID:   id,
		},
		VoiceFileID: voiceID,
		Title:       title,
	}
}

// NewInlineQueryResultDocument creates a new inline query document.
func NewInlineQueryResultDocument(id, url, title, mimeType string) InlineQueryResultDocument {
	return InlineQueryResultDocument{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "document",
			ID:   id,
		},
		URL:      url,
		Title:    title,
		MimeType: mimeType,
	}
}

// NewInlineQueryResultCachedDocument create a new inline query with cached photo.
func NewInlineQueryResultCachedDocument(id, documentID, title string) InlineQueryResultCachedDocument {
	return InlineQueryResultCachedDocument{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "document",
			ID:   id,
		},
		DocumentFileID: documentID,
		Title:          title,
	}
}

// NewInlineQueryResultLocation creates a new inline query location.
func NewInlineQueryResultLocation(id, title string, latitude, longitude float64) InlineQueryResultLocation {
	return InlineQueryResultLocation{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "location",
			ID:   id,
		},
		Title:     title,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewInlineQueryResultVenue creates a new inline query venue.
func NewInlineQueryResultVenue(id, title, address string, latitude, longitude float64) InlineQueryResultVenue {
	return InlineQueryResultVenue{
		InlineQueryResultBase: InlineQueryResultBase{
			Type: "venue",
			ID:   id,
		},
		Title:     title,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewEditMessageText allows you to edit the text of a message.
func NewEditMessageText(chatID int, messageID int, text string) EditMessageTextConf {
	return EditMessageTextConf{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
	}
}

// NewEditMessageTextAndMarkup allows you to edit the text and reply markup of a message.
func NewEditMessageTextAndMarkup(chatID int, messageID int, text string, replyMarkup InlineKeyboardMarkup) EditMessageTextConf {
	return EditMessageTextConf{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: &replyMarkup,
		Text:        text,
	}
}

// NewEditMessageCaption allows you to edit the caption of a message.
func NewEditMessageCaption(chatID int64, messageID int, caption string) EditMessageCaptionConf {
	return EditMessageCaptionConf{
		ChatID:    chatID,
		MessageID: messageID,
		Caption:   caption,
	}
}

// NewEditMessageReplyMarkup allows you to edit the inline
// keyboard markup.
func NewEditMessageReplyMarkup(chatID int, messageID int, replyMarkup InlineKeyboardMarkup) EditMessageReplyMarkupConf {
	return EditMessageReplyMarkupConf{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: &replyMarkup,
	}
}

// NewRemoveKeyboard hides the keyboard, with the option for being selective
// or hiding for everyone.
func NewRemoveKeyboard(selective bool) ReplyKeyboardRemove {
	return ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      selective,
	}
}

// NewKeyboardButton creates a regular keyboard button.
func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// NewKeyboardButtonWebApp creates a keyboard button with text
// which goes to a WebApp.
func NewKeyboardButtonWebApp(text string, webapp WebAppInfo) KeyboardButton {
	return KeyboardButton{
		Text:   text,
		WebApp: &webapp,
	}
}

// NewKeyboardButtonContact creates a keyboard button that requests
// user contact information upon click.
func NewKeyboardButtonContact(text string) KeyboardButton {
	return KeyboardButton{
		Text:           text,
		RequestContact: true,
	}
}

// NewKeyboardButtonLocation creates a keyboard button that requests
// user location information upon click.
func NewKeyboardButtonLocation(text string) KeyboardButton {
	return KeyboardButton{
		Text:            text,
		RequestLocation: true,
	}
}

// NewKeyboardButtonRow creates a row of keyboard buttons.
func NewKeyboardButtonRow(buttons ...KeyboardButton) []KeyboardButton {
	var row []KeyboardButton

	row = append(row, buttons...)

	return row
}

// NewReplyKeyboard creates a new regular keyboard with sane defaults.
func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	var keyboard [][]KeyboardButton

	keyboard = append(keyboard, rows...)

	return ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

// NewOneTimeReplyKeyboard creates a new one time keyboard.
func NewOneTimeReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	markup := NewReplyKeyboard(rows...)
	markup.OneTimeKeyboard = true
	return markup
}

// NewInlineKeyboardButtonData creates an inline keyboard button with text
// and data for a callback.
func NewInlineKeyboardButtonData(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: &data,
	}
}

// NewInlineKeyboardButtonWebApp creates an inline keyboard button with text
// which goes to a WebApp.
func NewInlineKeyboardButtonWebApp(text string, webapp WebAppInfo) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:   text,
		WebApp: &webapp,
	}
}

// NewInlineKeyboardButtonLoginURL creates an inline keyboard button with text
// which goes to a LoginURL.
func NewInlineKeyboardButtonLoginURL(text string, loginURL LoginURL) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:     text,
		LoginURL: &loginURL,
	}
}

// NewInlineKeyboardButtonURL creates an inline keyboard button with text
// which goes to a URL.
func NewInlineKeyboardButtonURL(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		URL:  &url,
	}
}

// NewInlineKeyboardButtonSwitch creates an inline keyboard button with
// text which allows the user to switch to a chat or return to a chat.
func NewInlineKeyboardButtonSwitch(text, sw string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:              text,
		SwitchInlineQuery: &sw,
	}
}

// NewInlineKeyboardRow creates an inline keyboard row with buttons.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// NewInlineKeyboardMarkup creates a new inline keyboard.
func NewInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton

	keyboard = append(keyboard, rows...)

	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

// NewCallback creates a new callback message.
func NewCallback(id, text string) AnswerCallbackQueryConf {
	return AnswerCallbackQueryConf{
		CallbackQueryID: id,
		Text:            text,
		ShowAlert:       false,
	}
}

// NewCallbackWithAlert creates a new callback message that alerts
// the user.
func NewCallbackWithAlert(id, text string) AnswerCallbackQueryConf {
	return AnswerCallbackQueryConf{
		CallbackQueryID: id,
		Text:            text,
		ShowAlert:       true,
	}
}

// NewInvoice creates a new Invoice request to the user.
func NewInvoice(chatID int, title, description, payload, providerToken, startParameter, currency string, prices []LabeledPrice) SendInvoiceConf {
	return SendInvoiceConf{
		ChatID:         chatID,
		Title:          title,
		Description:    description,
		Payload:        payload,
		ProviderToken:  providerToken,
		StartParameter: startParameter,
		Currency:       currency,
		Prices:         prices}
}

// NewChatTitle allows you to update the title of a chat.
func NewChatTitle(chatID int, title string) SetChatTitleConf {
	return SetChatTitleConf{
		ChatID: chatID,
		Title:  title,
	}
}

// NewChatDescription allows you to update the description of a chat.
func NewChatDescription(chatID int, description string) SetChatDescriptionConf {
	return SetChatDescriptionConf{
		ChatID:      chatID,
		Description: description,
	}
}

// NewChatPhoto allows you to update the photo for a chat.
func NewChatPhoto(chatID int, photo RequestFileData) SetChatPhotoConf {
	return SetChatPhotoConf{
		ChatID: chatID,
		File:   photo,
	}
}

// NewDeleteChatPhoto allows you to delete the photo for a chat.
func NewDeleteChatPhoto(chatID int) DeleteChatPhotoConf {
	return DeleteChatPhotoConf{
		ChatID: chatID,
	}
}

// NewPoll allows you to create a new poll.
func NewPoll(chatID int, question string, options ...string) SendPollConf {
	return SendPollConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
		Question:    question,
		Options:     options,
		IsAnonymous: true, // This is Telegram's default.
	}
}

// NewStopPoll allows you to stop a poll.
func NewStopPoll(chatID int, messageID int) StopPollConf {
	return StopPollConf{
		ChatID:    chatID,
		MessageID: messageID,
	}
}

// NewDice allows you to send a random dice roll.
func NewDice(chatID int) SendDiceConf {
	return SendDiceConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
	}
}

// NewDiceWithEmoji allows you to send a random roll of one of many types.
//
// Emoji may be 🎲 (1-6), 🎯 (1-6), or 🏀 (1-5).
func NewDiceWithEmoji(chatID int, emoji string) SendDiceConf {
	return SendDiceConf{
		BaseSend: BaseSend{
			ChatID: chatID,
		},
		Emoji: emoji,
	}
}

// NewBotCommandScopeDefault represents the default scope of bot commands.
func NewBotCommandScopeDefault() BotCommandScope {
	return BotCommandScope{Type: "default"}
}

// NewBotCommandScopeAllPrivateChats represents the scope of bot commands,
// covering all private chats.
func NewBotCommandScopeAllPrivateChats() BotCommandScope {
	return BotCommandScope{Type: "all_private_chats"}
}

// NewBotCommandScopeAllGroupChats represents the scope of bot commands,
// covering all group and supergroup chats.
func NewBotCommandScopeAllGroupChats() BotCommandScope {
	return BotCommandScope{Type: "all_group_chats"}
}

// NewBotCommandScopeAllChatAdministrators represents the scope of bot commands,
// covering all group and supergroup chat administrators.
func NewBotCommandScopeAllChatAdministrators() BotCommandScope {
	return BotCommandScope{Type: "all_chat_administrators"}
}

// NewBotCommandScopeChat represents the scope of bot commands, covering a
// specific chat.
func NewBotCommandScopeChat(chatID int) BotCommandScope {
	return BotCommandScope{
		Type:   "chat",
		ChatID: chatID,
	}
}

// NewBotCommandScopeChatAdministrators represents the scope of bot commands,
// covering all administrators of a specific group or supergroup chat.
func NewBotCommandScopeChatAdministrators(chatID int) BotCommandScope {
	return BotCommandScope{
		Type:   "chat_administrators",
		ChatID: chatID,
	}
}

// NewBotCommandScopeChatMember represents the scope of bot commands, covering a
// specific member of a group or supergroup chat.
func NewBotCommandScopeChatMember(chatID, userID int) BotCommandScope {
	return BotCommandScope{
		Type:   "chat_member",
		ChatID: chatID,
		UserID: userID,
	}
}

// NewGetMyCommandsWithScope allows you to set the registered commands for a
// given scope.
func NewGetMyCommandsWithScope(scope BotCommandScope) GetMyCommandsConf {
	return GetMyCommandsConf{Scope: &scope}
}

// NewGetMyCommandsWithScopeAndLanguage allows you to set the registered
// commands for a given scope and language code.
func NewGetMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string) GetMyCommandsConf {
	return GetMyCommandsConf{Scope: &scope, LanguageCode: languageCode}
}

// NewSetMyCommands allows you to set the registered commands.
func NewSetMyCommands(commands ...BotCommand) SetMyCommandsConf {
	return SetMyCommandsConf{Commands: commands}
}

// NewSetMyCommandsWithScope allows you to set the registered commands for a given scope.
func NewSetMyCommandsWithScope(scope BotCommandScope, commands ...BotCommand) SetMyCommandsConf {
	return SetMyCommandsConf{Commands: commands, Scope: &scope}
}

// NewSetMyCommandsWithScopeAndLanguage allows you to set the registered commands for a given scope
// and language code.
func NewSetMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string, commands ...BotCommand) SetMyCommandsConf {
	return SetMyCommandsConf{Commands: commands, Scope: &scope, LanguageCode: languageCode}
}

// NewDeleteMyCommands allows you to delete the registered commands.
func NewDeleteMyCommands() DeleteMyCommandsConf {
	return DeleteMyCommandsConf{}
}

// NewDeleteMyCommandsWithScope allows you to delete the registered commands for a given
// scope.
func NewDeleteMyCommandsWithScope(scope BotCommandScope) DeleteMyCommandsConf {
	return DeleteMyCommandsConf{Scope: &scope}
}

// NewDeleteMyCommandsWithScopeAndLanguage allows you to delete the registered commands for a given
// scope and language code.
func NewDeleteMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string) DeleteMyCommandsConf {
	return DeleteMyCommandsConf{Scope: &scope, LanguageCode: languageCode}
}

// ValidateWebAppData validate data received via the Web App
// https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app
func ValidateWebAppData(token, telegramInitData string) (bool, error) {
	initData, err := url.ParseQuery(telegramInitData)
	if err != nil {
		return false, fmt.Errorf("error parsing data %w", err)
	}

	dataCheckString := make([]string, 0, len(initData))
	for k, v := range initData {
		if k == "hash" {
			continue
		}
		if len(v) > 0 {
			dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}

	sort.Strings(dataCheckString)

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(token))

	hHash := hmac.New(sha256.New, secret.Sum(nil))
	hHash.Write([]byte(strings.Join(dataCheckString, "\n")))

	hash := hex.EncodeToString(hHash.Sum(nil))

	if initData.Get("hash") != hash {
		return false, errors.New("hash not equal")
	}

	return true, nil
}
