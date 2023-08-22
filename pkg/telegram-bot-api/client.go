package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client allows you to interact with the Telegram Bot API.
type Client struct {
	Host            string     // Telegram Bot API Host
	Token           string     // Telegram Bot API Token
	Debug           bool       // If true, enable debug logging
	Buffer          int        // Buffer size (default 100)
	Self            User       // Bot info from method getMe
	Client          HTTPClient //HTTP client
	botEndpoint     string     // Endpoint format: https://api.telegram.org/bot<token>
	fileEndpoint    string     // Endpoint format: https://api.telegram.org/file/bot<token>
	shutdownChannel chan interface{}
}

// New creates a new Client instance.
//
// It requires a token, provided by @BotFather on Telegram.
func New(token string) (*Client, error) {
	return NewWithClient(token, BaseEndpoint, &http.Client{})
}

// NewWithHost creates a new Client instance
// and allows you to pass API endpoint.
//
// Format: "https://api.telegram.org/"
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewWithHost(token, host string) (*Client, error) {
	return NewWithClient(token, host, &http.Client{})
}

// NewWithClient creates a new Client instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewWithClient(token, host string, client HTTPClient) (*Client, error) {
	bot := &Client{
		Host:            host,
		Token:           token,
		Client:          client,
		Buffer:          100,
		botEndpoint:     strings.TrimSuffix(host, "/") + "/bot" + token,
		fileEndpoint:    strings.TrimSuffix(host, "/") + "/file/bot" + token,
		shutdownChannel: make(chan interface{}),
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = *self

	return bot, nil
}

// UpdateEndpoint set new bot and file endpoints
// Always use a UpdateEndpoint if you change host.
func (client *Client) UpdateEndpoints() {
	client.botEndpoint = strings.TrimSuffix(client.Host, "/") + "/bot" + client.Token
	client.fileEndpoint = strings.TrimSuffix(client.Host, "/") + "/file/bot" + client.Token
}

// MakeRequest creates a request to send data.
// The transfer type is application/json, not suitable for file transfer. Accepts any struct with JSON tags.
func (client *Client) MakeRequest(method string, data interface{}) (*APIResponse, error) {
	if client.Debug {
		slog.Debug("Method: %s, data: %v\n", method, data)
	}

	url := client.botEndpoint + "/" + strings.TrimPrefix(method, "/")

	values, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(values))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := client.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if client.Debug {
		slog.Debug("Method: %s, response: %s\n", method, string(bytes))
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// MakeRequestWithFiles creates a request to send data.
// The transfer type is multipart/form-data, suitable for file transfer. Accepts any struct with JSON tags.
func (client *Client) MakeRequestWithFiles(method string, data interface{}, files []RequestFile) (*APIResponse, error) {
	values, err := structToMap(data)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		for _, val := range files {
			delete(values, val.Name)
		}

		for field, value := range values {
			if err := m.WriteField(field, value); err != nil {
				w.CloseWithError(err)
				return
			}
		}

		for _, file := range files {
			if file.Data.NeedsUpload() {
				name, reader, err := file.Data.SendData()
				if err != nil {
					w.CloseWithError(err)
					return
				}

				part, err := m.CreateFormFile(file.Name, name)
				if err != nil {
					w.CloseWithError(err)
					return
				}

				if _, err := io.Copy(part, reader); err != nil {
					w.CloseWithError(err)
					return
				}

				if closer, ok := reader.(io.ReadCloser); ok {
					if err = closer.Close(); err != nil {
						w.CloseWithError(err)
						return
					}
				}
			} else {
				value, _, _ := file.Data.SendData()

				if err := m.WriteField(file.Name, value); err != nil {
					w.CloseWithError(err)
					return
				}
			}
		}
	}()

	if client.Debug {
		slog.Debug("Method: %s, data: %v, with %d files\n", method, data, len(files))
	}

	url := client.botEndpoint + "/" + strings.TrimPrefix(method, "/")

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := client.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if client.Debug {
		slog.Debug("Method: %s, response: %s\n", method, string(bytes))
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// decodeAPIResponse decode response and return slice of bytes if debug enabled.
// If debug disabled, just decode http.Response.Body stream to APIResponse struct
// for efficient memory usage
func (client *Client) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) ([]byte, error) {
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func structToMap(data interface{}) (map[string]string, error) {
	result := make(map[string]string)

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct")
	}

	requestFileDataType := reflect.TypeOf((*RequestFileData)(nil)).Elem()

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonTag = strings.Split(jsonTag, ",")[0]
			if value.Kind() == reflect.Struct {
				if reflect.TypeOf(value).Implements(requestFileDataType) {
					continue
				}
				nestedData, err := json.Marshal(value.Interface())
				if err != nil {
					return nil, err
				}
				result[jsonTag] = string(nestedData)
			} else {
				result[jsonTag] = fmt.Sprintf("%v", value.Interface())
			}
		}
	}

	return result, nil
}

// Request sends a Config to Telegram, and returns the APIResponse.
func (client *Client) Request(c Config) (*APIResponse, error) {
	if t, ok := c.(ConfigWithFiles); ok {
		files := t.files()

		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFilesNeedingUpload(files) {
			return client.MakeRequestWithFiles(t.method(), c, files)
		}
	}

	return client.MakeRequest(c.method(), c)
}

func hasFilesNeedingUpload(files []RequestFile) bool {
	for _, file := range files {
		if file.Data.NeedsUpload() {
			return true
		}
	}

	return false
}

// RequestOK sends a Config to Telegram, and returns True on success.
//
// Use for all methods that return only True on success.
func (client *Client) RequestOK(c Config) (bool, error) {
	_, err := client.Request(c)
	if err != nil {
		return false, err
	}

	return true, nil
}

// EscapeText takes an input text and escape Telegram markup symbols.
// In this way we can send a text without being afraid of having to escape the characters manually.
// Note that you don't have to include the formatting style in the input text, or it will be escaped too.
// If there is an error, an empty string will be returned.
//
// parseMode is the text formatting mode (ModeMarkdown, ModeMarkdownV2 or ModeHTML)
// text is the input string that will be escaped
func EscapeText(parseMode string, text string) string {
	var replacer *strings.Replacer

	if parseMode == ModeHTML {
		replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	} else if parseMode == ModeMarkdown {
		replacer = strings.NewReplacer("_", "\\_", "*", "\\*", "`", "\\`", "[", "\\[")
	} else if parseMode == ModeMarkdownV2 {
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	} else {
		return ""
	}

	return replacer.Replace(text)
}

//
//
//
// Update
//
//
//

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, Timeout, and AllowedUpdates are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests, so you can get updates
// instantly instead of having to wait between requests.
func (client *Client) GetUpdates(config GetUpdatesConf) ([]Update, error) {
	resp, err := client.Request(config)
	if err != nil {
		return nil, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (client *Client) GetWebhookInfo() (*WebhookInfo, error) {
	resp, err := client.MakeRequest("getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var info *WebhookInfo
	if resp.Result != nil {
		info = new(WebhookInfo)
		err = json.Unmarshal(resp.Result, info)
		if err != nil {
			return nil, err
		}
	}

	return info, nil
}

// GetUpdatesChan starts and returns a channel for getting updates.
func (client *Client) GetUpdatesChan(config GetUpdatesConf) UpdatesChannel {
	ch := make(chan Update, client.Buffer)

	go func() {
		for {
			select {
			case <-client.shutdownChannel:
				close(ch)
				return
			default:
			}

			updates, err := client.GetUpdates(config)
			if err != nil {
				slog.Error(err.Error())
				slog.Info("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}

// StopReceivingUpdates stops the go routine which receives updates
func (client *Client) StopReceivingUpdates() {
	if client.Debug {
		slog.Debug("stopping the update receiver routine...")
	}
	close(client.shutdownChannel)
}

// ListenForWebhook registers a http handler for a webhook.
func (client *Client) ListenForWebhook(pattern string) UpdatesChannel {
	ch := make(chan Update, client.Buffer)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		update, err := client.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	})

	return ch
}

// ListenForWebhookRespReqFormat registers a http handler for a single incoming webhook.
func (client *Client) ListenForWebhookRespReqFormat(w http.ResponseWriter, r *http.Request) UpdatesChannel {
	ch := make(chan Update, client.Buffer)

	func(w http.ResponseWriter, r *http.Request) {
		defer close(ch)

		update, err := client.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	}(w, r)

	return ch
}

// HandleUpdate parses and returns update received via webhook
func (client *Client) HandleUpdate(r *http.Request) (*Update, error) {
	if r.Method != http.MethodPost {
		err := fmt.Errorf("wrong HTTP method required POST")
		return nil, err
	}

	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

//
//
//
// Methods
//
//
//

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (client *Client) GetMe() (*User, error) {
	resp, err := client.MakeRequest("getMe", nil)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Use this method to log out from the cloud Bot API server before launching the bot locally.
// You must log out the bot before running it locally,
// otherwise there is no guarantee that the bot will receive updates.
// After a successful call, you can immediately log in on a local server,
// but will not be able to log in back to the cloud Bot API server for 10 minutes.
// Returns True on success. Requires no parameters.
func (client *Client) LogOut() (bool, error) {
	_, err := client.MakeRequest("logOut", nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Use this method to close the bot instance before moving it from one local server to another.
// You need to delete the webhook before calling this method to ensure that
// the bot isn't launched again after server restart. The method will return error 429 in the first 10 minutes
// after the bot is launched. Returns True on success. Requires no parameters.
func (client *Client) Close() (bool, error) {
	_, err := client.MakeRequest("close", nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

// IsMessageToMe returns true if message directed to this bot.
//
// It requires the Message.
func (client *Client) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+client.Self.UserName)
}

// Send will send a Config item to Telegram and provides the returned Message.
//
// Use for all methods that return only Message on success.
func (client *Client) Send(c Config) (*Message, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// CopyMessage copy messages of any kind. The method is analogous to the method
// forwardMessage, but the copied message doesn't have a link to the original
// message. Returns the MessageID of the sent message on success.
func (client *Client) CopyMessage(c CopyMessageConf) (*MessageId, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var messageId MessageId
	err = json.Unmarshal(resp.Result, &messageId)
	if err != nil {
		return nil, err
	}

	return &messageId, nil
}

// SendMediaGroup sends a media group and returns the resulting messages.
func (client *Client) SendMediaGroup(c SendMediaGroupConf) ([]Message, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var messages []Message
	err = json.Unmarshal(resp.Result, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetUserProfilePhotos gets a user's profile photos.
//
// It requires UserID.
// Offset and Limit are optional.
func (client *Client) GetUserProfilePhotos(c GetUserProfilePhotosConf) (*UserProfilePhotos, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var photo UserProfilePhotos
	err = json.Unmarshal(resp.Result, &photo)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

// GetFile returns a File which can download a file from Telegram.
//
// Requires FileID.
func (client *Client) GetFile(c GetFileConf) (*File, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var file File
	err = json.Unmarshal(resp.Result, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// ExportChatInviteLink returns the generated a new primary invite link for a chat.
//
// Requires ChatID.
func (client *Client) ExportChatInviteLink(c ExportChatInviteLinkConf) (string, error) {
	resp, err := client.Request(c)
	if err != nil {
		return "", err
	}

	return string(resp.Result), nil
}

// ExportChatInviteLink returns the new invite link for a chat.
//
// Requires ChatID.
func (client *Client) CreateChatInviteLink(c CreateChatInviteLinkConf) (*ChatInviteLink, error) {
	return client.chatInviteLink(c)
}

// RevokeChatInviteLink returns the edited invite link created by the bot.
//
// Requires ChatID and InviteLink
func (client *Client) EditChatInviteLink(c EditChatInviteLinkConf) (*ChatInviteLink, error) {
	return client.chatInviteLink(c)
}

// RevokeChatInviteLink returns the revoked invite link created by the bot.
//
// Requires ChatID and InviteLink
func (client *Client) RevokeChatInviteLink(c RevokeChatInviteLinkConf) (*ChatInviteLink, error) {
	return client.chatInviteLink(c)
}

func (client *Client) chatInviteLink(c Config) (*ChatInviteLink, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var link ChatInviteLink
	err = json.Unmarshal(resp.Result, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// GetChat gets information about a chat.
func (client *Client) GetChat(c GetChatConf) (*Chat, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var chat Chat
	err = json.Unmarshal(resp.Result, &chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// GetChatAdministrators gets a list of administrators in the chat.
//
// If none have been appointed, only the creator will be returned.
// Bots are not shown, even if they are an administrator.
func (client *Client) GetChatAdministrators(c GetChatAdministratorsConf) ([]ChatMember, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var members []ChatMember
	err = json.Unmarshal(resp.Result, &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// GetChatMembersCount gets the number of users in a chat.
func (client *Client) GetChatMemberCount(c GetChatMemberCountConf) (int, error) {
	resp, err := client.Request(c)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(resp.Result))
}

// GetChatMember gets a specific chat member.
func (client *Client) GetChatMember(c GetChatMemberConf) (*ChatMember, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var member ChatMember
	err = json.Unmarshal(resp.Result, &member)
	if err != nil {
		return nil, err
	}

	return &member, nil
}

// GetForumTopicIconStickers gets a custom emoji stickers, which can be used as a forum topic icon by any user.
func (client *Client) GetForumTopicIconStickers() ([]Sticker, error) {
	resp, err := client.MakeRequest("getForumTopicIconStickers", nil)
	if err != nil {
		return nil, err
	}

	var stickers []Sticker
	err = json.Unmarshal(resp.Result, &stickers)
	if err != nil {
		return nil, err
	}

	return stickers, nil
}

// CreateForumTopic create a topic in a forum supergroup chat.
func (client *Client) CreateForumTopic(c CreateForumTopicConf) (*ForumTopic, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var topic ForumTopic
	err = json.Unmarshal(resp.Result, &topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// GetMyCommands gets the currently registered commands.
//
// Returns nil if no commands.
func (client *Client) GetMyCommands(c GetMyCommandsConf) ([]BotCommand, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var commands []BotCommand
	err = json.Unmarshal(resp.Result, &commands)
	if err != nil {
		return nil, nil
	}

	return commands, nil
}

// GetMyName returns bot name.
func (client *Client) GetMyName(c GetMyNameConf) (*BotName, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var name BotName
	err = json.Unmarshal(resp.Result, &name)
	if err != nil {
		return nil, err
	}

	return &name, nil
}

// GetMyDescription returns bot description.
func (client *Client) GetMyDescription(c GetMyDescriptionConf) (*BotDescription, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var description BotDescription
	err = json.Unmarshal(resp.Result, &description)
	if err != nil {
		return nil, err
	}

	return &description, nil
}

// GetMyShortDescription returns bot short description.
func (client *Client) GetMyShortDescription(c GetMyShortDescriptionConf) (*BotShortDescription, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var description BotShortDescription
	err = json.Unmarshal(resp.Result, &description)
	if err != nil {
		return nil, err
	}

	return &description, nil
}

// GetChatMenuButton gets the current value of the bot's menu button in a private chat, or the default menu button.
func (client *Client) GetChatMenuButton(c GetChatMenuButtonConf) (*MenuButton, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var button MenuButton
	err = json.Unmarshal(resp.Result, &button)
	if err != nil {
		return nil, err
	}

	return &button, nil
}

// GetMyDefaultAdministratorRights gets the current default administrator rights of the bot.
func (client *Client) GetMyDefaultAdministratorRights(c GetMyDefaultAdministratorRightsConf) (*ChatAdministratorRights, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var rights ChatAdministratorRights
	err = json.Unmarshal(resp.Result, &rights)
	if err != nil {
		return nil, err
	}

	return &rights, nil
}

//
//
//
// Messages
//
//
//

// On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
//
// Use for all EditMessage methods.
func (client *Client) EditMessage(c Config) (*Message, bool, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, resp.Ok, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)
	if err != nil {
		return nil, resp.Ok, nil
	}

	return &message, resp.Ok, nil
}

// StopPoll stops a poll and returns the result.
func (client *Client) StopPoll(c StopPollConf) (*Poll, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var poll Poll
	err = json.Unmarshal(resp.Result, &poll)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

//
//
//
// Stickers
//
//
//

// GetStickerSet returns a StickerSet.
func (client *Client) GetStickerSet(c GetStickerSetConf) (*StickerSet, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var stickerSet StickerSet
	err = json.Unmarshal(resp.Result, &stickerSet)
	if err != nil {
		return nil, err
	}

	return &stickerSet, nil
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their identifiers.
func (client *Client) GetCustomEmojiStickers(c GetCustomEmojiStickersConf) ([]Sticker, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var stickers []Sticker
	err = json.Unmarshal(resp.Result, &stickers)
	if err != nil {
		return nil, err
	}

	return stickers, nil
}

// UploadStickerFile upload a file with a sticker for later use in the
// createNewStickerSet and addStickerToSet methods (the file can be used multiple times)
func (client *Client) UploadStickerFile(c UploadStickerFileConf) (*File, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var file File
	err = json.Unmarshal(resp.Result, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

//
//
//
// Inline
//
//
//

// AnswerWebAppQuery set the result of an interaction with a Web App and
// send a corresponding message on behalf of the user to the chat from which the query originated.
func (client *Client) AnswerWebAppQuery(c AnswerWebAppQueryConf) (*SentWebAppMessage, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var message SentWebAppMessage
	err = json.Unmarshal(resp.Result, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

//
//
//
// Payment
//
//
//

// CreateInvoiceLink create a link for an invoice.
func (client *Client) CreateInvoiceLink(c CreateInvoiceLinkConf) (string, error) {
	resp, err := client.Request(c)
	if err != nil {
		return "", err
	}
	return string(resp.Result), nil
}

//
//
//
// Game
//
//
//

// SetGameScore set the score of the specified user in a game message.
//
// On success, if the message is not an inline message, the Message is returned,
// otherwise True is returned. Returns an error, if the new score is not greater
// than the user's current score in the chat and force is False.
func (client *Client) SetGameScore(c SetGameScoreConf) (*Message, bool, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, resp.Ok, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)
	if err != nil {
		return nil, resp.Ok, nil
	}

	return &message, resp.Ok, nil
}

// GetGameHighScores allows you to get the high scores for a game.
func (client *Client) GetGameHighScores(c GetGameHighScoresConf) ([]GameHighScore, error) {
	resp, err := client.Request(c)
	if err != nil {
		return nil, err
	}

	var score []GameHighScore
	err = json.Unmarshal(resp.Result, &score)
	if err != nil {
		return nil, err
	}

	return score, nil
}
