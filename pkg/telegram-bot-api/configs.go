package telegram

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
)

// Telegram constants
const (
	// Base Telegram Bot API endpoint
	BaseEndpoint string = "https://api.telegram.org/"
)

// Constant values for ChatActions
const (
	ChatTyping          = "typing"
	ChatUploadPhoto     = "upload_photo"
	ChatRecordVideo     = "record_video"
	ChatUploadVideo     = "upload_video"
	ChatRecordVoice     = "record_voice"
	ChatUploadVoice     = "upload_voice"
	ChatUploadDocument  = "upload_document"
	ChatChooseSticker   = "choose_sticker"
	ChatFindLocation    = "find_location"
	ChatRecordVideoNote = "record_video_note"
	ChatUploadVideoNote = "upload_video_note"
)

// Constant values for ParseMode in MessageConfig
const (
	ModeMarkdown   = "Markdown"
	ModeMarkdownV2 = "MarkdownV2"
	ModeHTML       = "HTML"
)

// Conf is any config type that can be sent.
type Config interface {
	method() string
}

// Conf is any config type that can be sent that includes a file.
type ConfigWithFiles interface {
	Config
	files() []RequestFile
}

// RequestFile represents a file associated with a field name.
type RequestFile struct {
	// The file field name.
	Name string
	// The file data to include.
	Data RequestFileData
}

// RequestFileData represents the data to be used for a file.
type RequestFileData interface {
	// NeedsUpload shows if the file needs to be uploaded.
	NeedsUpload() bool
	// SendData gets the file name and an `io.Reader` for the file to be uploaded or the file data to send (io.Reader = nil)
	// when a file does not need to be uploaded.
	SendData() (string, io.Reader, error)
}

// FileBytes contains information about a set of bytes to upload
// as a File.
type FileBytes struct {
	Name  string
	Bytes []byte
}

func (fb FileBytes) NeedsUpload() bool {
	return true
}

func (fb FileBytes) SendData() (string, io.Reader, error) {
	return fb.Name, bytes.NewReader(fb.Bytes), nil
}

// FileReader contains information about a reader to upload as a File.
type FileReader struct {
	Name   string
	Reader io.Reader
}

func (fr FileReader) NeedsUpload() bool {
	return true
}

func (fr FileReader) SendData() (string, io.Reader, error) {
	return fr.Name, fr.Reader, nil
}

// FilePath is a path to a local file.
type FilePath string

func (fp FilePath) NeedsUpload() bool {
	return true
}

func (fp FilePath) SendData() (string, io.Reader, error) {
	fileHandle, err := os.Open(string(fp))
	if err != nil {
		return "", nil, err
	}

	name := fileHandle.Name()
	return name, fileHandle, err
}

// FileURL is a URL to use as a file for a request.
type FileURL string

func (fu FileURL) NeedsUpload() bool {
	return false
}

func (fu FileURL) SendData() (string, io.Reader, error) {
	return string(fu), nil, nil
}

// FileID is an ID of a file already uploaded to Telegram.
type FileID string

func (fi FileID) NeedsUpload() bool {
	return false
}

func (fi FileID) SendData() (string, io.Reader, error) {
	return string(fi), nil, nil
}

// fileAttach is an internal file type used for processed media groups.
type fileAttach string

func (fa fileAttach) NeedsUpload() bool {
	return false
}

func (fa fileAttach) SendData() (string, io.Reader, error) {
	return string(fa), nil, nil
}

//
//
//
// Update
//
//
//

// GetUpdatesConf contains fields for the getUpdates method. Returns an Array of Update objects.
type GetUpdatesConf struct {
	Offset         int      `json:"offset,omitempty"`          // Optional. Identifier of the first update to be returned.
	Limit          int      `json:"limit,omitempty"`           // Optional. Limits the number of updates to be retrieved.
	Timeout        int      `json:"timeout,omitempty"`         // Optional. Timeout in seconds for long polling.
	AllowedUpdates []string `json:"allowed_updates,omitempty"` // Optional. A list of the update types you want your bot to receive.
}

func (c GetUpdatesConf) method() string {
	return "getUpdates"
}

// SetWebhookConf contains fields for the setWebhook method. Returns True on success.
type SetWebhookConf struct {
	URL                *url.URL        `json:"url"`                            // HTTPS URL to send updates to.
	Certificate        RequestFileData `json:"certificate,omitempty"`          // Optional. Public key certificate.
	IPAddress          string          `json:"ip_address,omitempty"`           // Optional. Fixed IP address.
	MaxConnections     int             `json:"max_connections,omitempty"`      // Optional. Maximum allowed number of simultaneous HTTPS connections.
	AllowedUpdates     []string        `json:"allowed_updates,omitempty"`      // Optional. A list of the update types you want your bot to receive.
	DropPendingUpdates bool            `json:"drop_pending_updates,omitempty"` // Optional. Pass True to drop all pending updates.
	SecretToken        string          `json:"secret_token,omitempty"`         // Optional. A secret token to be sent in a header.
}

func (c SetWebhookConf) method() string {
	return "setWebhook"
}

// DeleteWebhookConf contains fields for the deleteWebhook method. Returns True on success.
type DeleteWebhookConf struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"` // Optional. Pass True to drop all pending updates.
}

func (c DeleteWebhookConf) method() string {
	return "deleteWebhook"
}

//
//
//
// Methods
//
//
//

// ForwardMessageConf contains fields for the forwardMessage method. On success, the sent Message is returned.
type ForwardMessageConf struct {
	ChatID              interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target channel
	MessageThreadID     int         `json:"message_thread_id,omitempty"`    // Optional. Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	FromChatID          interface{} `json:"from_chat_id"`                   // Unique identifier for the chat where the original message was sent
	DisableNotification bool        `json:"disable_notification,omitempty"` // Optional. Sends the message silently
	ProtectContent      bool        `json:"protect_content,omitempty"`      // Optional. Protects the contents of the forwarded message from forwarding and saving
	MessageID           int         `json:"message_id"`                     // Message identifier in the chat specified in from_chat_id
}

func (c ForwardMessageConf) method() string {
	return "forwardMessage"
}

type BaseSend struct {
	ChatID                   interface{} `json:"chat_id"`                               // Unique identifier for the target chat or username of the target channel
	MessageThreadID          int         `json:"message_thread_id,omitempty"`           // Optional. Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	DisableNotification      bool        `json:"disable_notification,omitempty"`        // Optional. Sends the message silently
	ProtectContent           bool        `json:"protect_content,omitempty"`             // Optional. Protects the contents of the sent message from forwarding and saving
	ReplyToMessageID         int         `json:"reply_to_message_id,omitempty"`         // Optional. If the message is a reply, ID of the original message
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply,omitempty"` // Optional. Pass true if the message should be sent even if the specified replied-to message is not found
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`                // Optional. Additional interface options
}

// SendMessageConf contains fields for the sendMessage method. On success, the sent Message is returned.
type SendMessageConf struct {
	BaseSend                              // Unique identifier for the target chat or username of the target channel
	Text                  string          `json:"text"`                               // Text of the message to be sent
	ParseMode             string          `json:"parse_mode,omitempty"`               // Optional. Mode for parsing entities in the message text
	Entities              []MessageEntity `json:"entities,omitempty"`                 // Optional. Special entities that appear in the message text
	DisableWebPagePreview bool            `json:"disable_web_page_preview,omitempty"` // Optional. Disables link previews for links in the message
}

func (c SendMessageConf) method() string {
	return "sendMessage"
}

// CopyMessageConf contains fields for the copyMessage method. Returns the MessageId of the sent message on success.
type CopyMessageConf struct {
	BaseSend                        // Unique identifier for the target chat or username of the target channel
	FromChatID      interface{}     `json:"from_chat_id"`               // Unique identifier for the chat where the original message was sent
	MessageID       int             `json:"message_id"`                 // Message identifier in the chat specified in from_chat_id
	Caption         string          `json:"caption,omitempty"`          // Optional. New caption for media
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the new caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. Special entities that appear in the new caption
}

func (c CopyMessageConf) method() string {
	return "copyMessage"
}

// SendPhotoConf contains fields for the sendPhoto method. On success, the sent Message is returned.
type SendPhotoConf struct {
	BaseSend                        // Unique identifier for the target chat or username of the target channel
	File            RequestFileData `json:"photo"`                      // Photo to send
	Caption         string          `json:"caption,omitempty"`          // Optional. Photo caption
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the photo caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. Special entities that appear in the caption
	HasSpoiler      bool            `json:"has_spoiler,omitempty"`      // Optional. Pass True if the photo needs to be covered with a spoiler animation
}

func (c SendPhotoConf) method() string {
	return "sendPhoto"
}

func (config *SendPhotoConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "photo",
		Data: config.File,
	}}

	return files
}

// SendAudioConf contains fields for the sendAudio method. On success, the sent Message is returned.
type SendAudioConf struct {
	BaseSend                        // Unique identifier for the target chat or username of the target channel
	File            RequestFileData `json:"audio"`                      // Audio file to send
	Caption         string          `json:"caption,omitempty"`          // Optional. Audio caption
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the audio caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. Special entities that appear in the caption
	Duration        int             `json:"duration,omitempty"`         // Optional. Duration of the audio in seconds
	Performer       string          `json:"performer,omitempty"`        // Optional. Performer
	Title           string          `json:"title,omitempty"`            // Optional. Track name
	Thumbnail       RequestFileData `json:"thumbnail,omitempty"`        // Optional. Thumbnail of the file sent
}

func (c SendAudioConf) method() string {
	return "sendAudio"
}

func (config *SendAudioConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "audio",
		Data: config.File,
	}}

	if config.Thumbnail != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumbnail,
		})
	}

	return files
}

// SendDocumentConf contains fields for the sendDocument method. On success, the sent Message is returned.
type SendDocumentConf struct {
	BaseSend                                    // Unique identifier for the target chat or username of the target channel
	File                        RequestFileData `json:"document"`                                 // File to send
	Thumbnail                   RequestFileData `json:"thumbnail,omitempty"`                      // Optional. Thumbnail of the file sent
	Caption                     string          `json:"caption,omitempty"`                        // Optional. Document caption
	ParseMode                   string          `json:"parse_mode,omitempty"`                     // Optional. Mode for parsing entities in the document caption
	CaptionEntities             []MessageEntity `json:"caption_entities,omitempty"`               // Optional. Special entities that appear in the caption
	DisableContentTypeDetection bool            `json:"disable_content_type_detection,omitempty"` // Optional. Disables automatic server-side content type detection for files uploaded using multipart/form-data
}

func (c SendDocumentConf) method() string {
	return "sendDocument"
}

func (config *SendDocumentConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "document",
		Data: config.File,
	}}

	if config.Thumbnail != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumbnail,
		})
	}

	return files
}

// SendVideoConf contains fields for the sendVideo method. On success, the sent Message is returned.
type SendVideoConf struct {
	BaseSend                          // Unique identifier for the target chat or username of the target channel
	File              RequestFileData `json:"video"`                        // Video to send
	Duration          int             `json:"duration,omitempty"`           // Optional. Duration of sent video in seconds
	Width             int             `json:"width,omitempty"`              // Optional. Video width
	Height            int             `json:"height,omitempty"`             // Optional. Video height
	Thumbnail         RequestFileData `json:"thumbnail,omitempty"`          // Optional. Thumbnail of the video
	Caption           string          `json:"caption,omitempty"`            // Optional. Video caption
	ParseMode         string          `json:"parse_mode,omitempty"`         // Optional. Mode for parsing entities in the video caption
	CaptionEntities   []MessageEntity `json:"caption_entities,omitempty"`   // Optional. Special entities that appear in the video caption
	HasSpoiler        bool            `json:"has_spoiler,omitempty"`        // Optional. Pass true if the video needs to be covered with a spoiler animation
	SupportsStreaming bool            `json:"supports_streaming,omitempty"` // Optional. Pass true if the uploaded video is suitable for streaming
}

func (c SendVideoConf) method() string {
	return "sendVideo"
}

func (config *SendVideoConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "video",
		Data: config.File,
	}}

	if config.Thumbnail != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumbnail,
		})
	}

	return files
}

// SendAnimationConf contains fields for the sendAnimation method. On success, the sent Message is returned.
type SendAnimationConf struct {
	BaseSend                        // Unique identifier for the target chat or username of the target channel
	File            RequestFileData `json:"animation"`                  // Animation to send
	Duration        int             `json:"duration,omitempty"`         // Optional. Duration of sent animation in seconds
	Width           int             `json:"width,omitempty"`            // Optional. Animation width
	Height          int             `json:"height,omitempty"`           // Optional. Animation height
	Thumbnail       RequestFileData `json:"thumbnail,omitempty"`        // Optional. Thumbnail of the animation
	Caption         string          `json:"caption,omitempty"`          // Optional. Animation caption
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the animation caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. Special entities that appear in the animation caption
	HasSpoiler      bool            `json:"has_spoiler,omitempty"`      // Optional. Pass true if the animation needs to be covered with a spoiler animation
}

func (c SendAnimationConf) method() string {
	return "sendAnimation"
}

func (config *SendAnimationConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "animation",
		Data: config.File,
	}}

	if config.Thumbnail != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumbnail,
		})
	}

	return files
}

// SendVoiceConf contains fields for the sendVoice method. On success, the sent Message is returned.
type SendVoiceConf struct {
	BaseSend                        // Unique identifier for the target chat or username of the target channel
	File            RequestFileData `json:"voice"`                      // Voice audio file to send
	Caption         string          `json:"caption,omitempty"`          // Optional. Voice message caption
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the voice message caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. Special entities that appear in the voice message caption
	Duration        int             `json:"duration,omitempty"`         // Optional. Duration of the voice message in seconds
}

func (c SendVoiceConf) method() string {
	return "sendVoice"
}

func (config *SendVoiceConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "voice",
		Data: config.File,
	}}

	return files
}

// SendVideoNoteConf contains fields for the sendVideoNote method. On success, the sent Message is returned.
type SendVideoNoteConf struct {
	BaseSend                  // Unique identifier for the target chat or username of the target channel
	File      RequestFileData `json:"video_note"`          // Video note to send
	Duration  int             `json:"duration,omitempty"`  // Optional. Duration of sent video in seconds
	Length    int             `json:"length,omitempty"`    // Optional. Video width and height
	Thumbnail RequestFileData `json:"thumbnail,omitempty"` // Optional. Thumbnail of the file sent
}

func (c SendVideoNoteConf) method() string {
	return "sendVideoNote"
}

func (config *SendVideoNoteConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "video_note",
		Data: config.File,
	}}

	if config.Thumbnail != nil {
		files = append(files, RequestFile{
			Name: "thumbnail",
			Data: config.Thumbnail,
		})
	}

	return files
}

// SendMediaGroupConf contains fields for the sendMediaGroup method. On success, an array of Messages that were sent is returned.
type SendMediaGroupConf struct {
	ChatID                   interface{}   `json:"chat_id"`                               // Unique identifier for the target chat or username of the target channel
	MessageThreadID          int           `json:"message_thread_id,omitempty"`           // Optional. Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	Media                    []interface{} `json:"media"`                                 // A JSON-serialized array describing messages to be sent
	DisableNotification      bool          `json:"disable_notification,omitempty"`        // Optional. Sends messages silently
	ProtectContent           bool          `json:"protect_content,omitempty"`             // Optional. Protects the contents of the sent messages from forwarding and saving
	ReplyToMessageID         int           `json:"reply_to_message_id,omitempty"`         // Optional. If the messages are a reply, ID of the original message
	AllowSendingWithoutReply bool          `json:"allow_sending_without_reply,omitempty"` // Optional. Pass True if the message should be sent even if the specified replied-to message is not found
}

func (c SendMediaGroupConf) method() string {
	return "sendMediaGroup"
}

func (config *SendMediaGroupConf) Files() []RequestFile {
	return prepareMediaGroup(config.Media)
}

func prepareMediaGroup(inputMedia []interface{}) []RequestFile {
	files := []RequestFile{}

	for idx, media := range inputMedia {
		switch m := media.(type) {
		case *InputMediaPhoto:
			if m.Media.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d", idx),
					Data: m.Media,
				})
				m.Media = fileAttach(fmt.Sprintf("attach://file-%d", idx))
			}
		case *InputMediaVideo:
			if m.Media.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d", idx),
					Data: m.Media,
				})
				m.Media = fileAttach(fmt.Sprintf("attach://file-%d", idx))
			}

			if m.Thumbnail != nil && m.Thumbnail.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d-thumbnail", idx),
					Data: m.Thumbnail,
				})
				m.Thumbnail = fileAttach(fmt.Sprintf("attach://file-%d-thumbnail", idx))
			}
		case *InputMediaAnimation:
			if m.Media.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d", idx),
					Data: m.Media,
				})
				m.Media = fileAttach(fmt.Sprintf("attach://file-%d", idx))
			}

			if m.Thumbnail != nil && m.Thumbnail.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d-thumbnail", idx),
					Data: m.Thumbnail,
				})
				m.Thumbnail = fileAttach(fmt.Sprintf("attach://file-%d-thumbnail", idx))
			}
		case *InputMediaDocument:
			if m.Media.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d", idx),
					Data: m.Media,
				})
				m.Media = fileAttach(fmt.Sprintf("attach://file-%d", idx))
			}

			if m.Thumbnail != nil && m.Thumbnail.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d-thumbnail", idx),
					Data: m.Thumbnail,
				})
				m.Thumbnail = fileAttach(fmt.Sprintf("attach://file-%d-thumbnail", idx))
			}
		case *InputMediaAudio:
			if m.Media.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d", idx),
					Data: m.Media,
				})
				m.Media = fileAttach(fmt.Sprintf("attach://file-%d", idx))
			}

			if m.Thumbnail != nil && m.Thumbnail.NeedsUpload() {
				files = append(files, RequestFile{
					Name: fmt.Sprintf("file-%d-thumbnail", idx),
					Data: m.Thumbnail,
				})
				m.Thumbnail = fileAttach(fmt.Sprintf("attach://file-%d-thumbnail", idx))
			}
		}
	}

	return files
}

// SendLocationConf contains fields for the sendLocation method. On success, the sent Message is returned.
type SendLocationConf struct {
	BaseSend                     // Unique identifier for the target chat or username of the target channel
	Latitude             float64 `json:"latitude"`                         // Latitude of the location
	Longitude            float64 `json:"longitude"`                        // Longitude of the location
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`    // Optional. The radius of uncertainty for the location, measured in meters
	LivePeriod           int     `json:"live_period,omitempty"`            // Optional. Period in seconds for which the location will be updated
	Heading              int     `json:"heading,omitempty"`                // Optional. For live locations, a direction in which the user is moving, in degrees
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"` // Optional. For live locations, a maximum distance for proximity alerts about approaching another chat member
}

func (c SendLocationConf) method() string {
	return "sendLocation"
}

// SendVenueConf contains fields for the sendVenue method. On success, the sent Message is returned.
type SendVenueConf struct {
	BaseSend                // Unique identifier for the target chat or username of the target channel
	Latitude        float64 `json:"latitude"`                    // Latitude of the venue
	Longitude       float64 `json:"longitude"`                   // Longitude of the venue
	Title           string  `json:"title"`                       // Name of the venue
	Address         string  `json:"address"`                     // Address of the venue
	FoursquareID    string  `json:"foursquare_id,omitempty"`     // Optional. Foursquare identifier of the venue
	FoursquareType  string  `json:"foursquare_type,omitempty"`   // Optional. Foursquare type of the venue
	GooglePlaceID   string  `json:"google_place_id,omitempty"`   // Optional. Google Places identifier of the venue
	GooglePlaceType string  `json:"google_place_type,omitempty"` // Optional. Google Places type of the venue
}

func (c SendVenueConf) method() string {
	return "sendVenue"
}

// SendContactConf contains fields for the sendContact method. On success, the sent Message is returned.
type SendContactConf struct {
	BaseSend               // Unique identifier for the target chat or username of the target channel
	MessageThreadID int    `json:"message_thread_id,omitempty"` // Optional. Unique identifier for the target message thread of the forum
	PhoneNumber     string `json:"phone_number"`                // Contact's phone number
	FirstName       string `json:"first_name"`                  // Contact's first name
	LastName        string `json:"last_name,omitempty"`         // Optional. Contact's last name
	VCard           string `json:"vcard,omitempty"`             // Optional. Additional data about the contact in the form of a vCard
}

func (c SendContactConf) method() string {
	return "sendContact"
}

// SendPollConf contains fields for the sendPoll method. On success, the sent Message is returned.
type SendPollConf struct {
	BaseSend                              // Unique identifier for the target chat or username of the target channel
	Question              string          `json:"question"`                          // Poll question
	Options               []string        `json:"options"`                           // A list of answer options
	IsAnonymous           bool            `json:"is_anonymous,omitempty"`            // Optional. True, if the poll needs to be anonymous
	Type                  string          `json:"type,omitempty"`                    // Optional. Poll type, "quiz" or "regular"
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"` // Optional. True, if the poll allows multiple answers
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`       // Optional. 0-based identifier of the correct answer option
	Explanation           string          `json:"explanation,omitempty"`             // Optional. Text shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll
	ExplanationParseMode  string          `json:"explanation_parse_mode,omitempty"`  // Optional. Mode for parsing entities in the explanation
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`    // Optional. Special entities that appear in the poll explanation
	OpenPeriod            int             `json:"open_period,omitempty"`             // Optional. Amount of time in seconds the poll will be active after creation
	CloseDate             int             `json:"close_date,omitempty"`              // Optional. Point in time when the poll will be automatically closed
	IsClosed              bool            `json:"is_closed,omitempty"`               // Optional. True, if the poll needs to be immediately closed
}

func (c SendPollConf) method() string {
	return "sendPoll"
}

// SendDiceConf contains fields for the sendDice method. On success, the sent Message is returned.
type SendDiceConf struct {
	BaseSend        // Unique identifier for the target chat or username of the target channel
	Emoji    string `json:"emoji,omitempty"` // Optional. Emoji on which the dice throw animation is based
}

func (c SendDiceConf) method() string {
	return "sendDice"
}

// SendChatActionConf contains fields for the sendChatAction method. Returns True on success.
type SendChatActionConf struct {
	ChatID          interface{} `json:"chat_id"`                     // Unique identifier for the target chat or username of the target channel
	MessageThreadID int         `json:"message_thread_id,omitempty"` // Optional. Unique identifier for the target message thread of the forum
	Action          string      `json:"action"`                      // Type of action to broadcast
}

func (c SendChatActionConf) method() string {
	return "sendChatAction"
}

// GetUserProfilePhotosConf contains fields for the getUserProfilePhotos method. Returns a UserProfilePhotos object.
type GetUserProfilePhotosConf struct {
	UserID int `json:"user_id"`          // Unique identifier of the target user
	Offset int `json:"offset,omitempty"` // Optional. Sequential number of the first photo to be returned
	Limit  int `json:"limit,omitempty"`  // Optional. Limits the number of photos to be retrieved
}

func (c GetUserProfilePhotosConf) method() string {
	return "getUserProfilePhotos"
}

// GetFileConf contains fields for the getFile method. On success, a File object is returned.
type GetFileConf struct {
	FileID string `json:"file_id"` // File identifier to get information about
}

func (c GetFileConf) method() string {
	return "getFile"
}

// BanChatMemberConf contains fields for the banChatMember method. Returns True on success.
type BanChatMemberConf struct {
	ChatID     interface{} `json:"chat_id"`                   // Unique identifier for the target group or username of the target supergroup or channel (in the format @channelusername)
	UserID     int         `json:"user_id"`                   // Unique identifier of the target user
	UntilDate  int         `json:"until_date,omitempty"`      // Optional. Date when the user will be unbanned, unix time
	RevokeMsgs bool        `json:"revoke_messages,omitempty"` // Optional. Pass True to delete all messages from the chat for the user that is being removed
}

func (c BanChatMemberConf) method() string {
	return "banChatMember"
}

// UnbanChatMemberConf contains fields for the unbanChatMember method. Returns True on success.
type UnbanChatMemberConf struct {
	ChatID       interface{} `json:"chat_id"`                  // Unique identifier for the target group or username of the target supergroup or channel (in the format @channelusername)
	UserID       int         `json:"user_id"`                  // Unique identifier of the target user
	OnlyIfBanned bool        `json:"only_if_banned,omitempty"` // Optional. Do nothing if the user is not banned
}

func (c UnbanChatMemberConf) method() string {
	return "unbanChatMember"
}

// RestrictChatMemberConf contains fields for the restrictChatMember method. Returns True on success.
type RestrictChatMemberConf struct {
	ChatID              interface{}     `json:"chat_id"`                                    // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	UserID              int             `json:"user_id"`                                    // Unique identifier of the target user
	Permissions         ChatPermissions `json:"permissions"`                                // A JSON-serialized object for new user permissions
	UseIndependentPerms bool            `json:"use_independent_chat_permissions,omitempty"` // Optional. Pass True if chat permissions are set independently
	UntilDate           int             `json:"until_date,omitempty"`                       // Optional. Date when restrictions will be lifted for the user, unix time
}

func (c RestrictChatMemberConf) method() string {
	return "restrictChatMember"
}

// PromoteChatMemberConf contains fields for the promoteChatMember method. Returns True on success.
type PromoteChatMemberConf struct {
	ChatID              interface{} `json:"chat_id"`                          // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	UserID              int         `json:"user_id"`                          // Unique identifier of the target user
	IsAnonymous         bool        `json:"is_anonymous,omitempty"`           // Optional. Pass True if the administrator's presence in the chat is hidden
	CanManageChat       bool        `json:"can_manage_chat,omitempty"`        // Optional. Pass True if the administrator can access the chat event log, chat statistics, message statistics in channels, see channel members, see anonymous administrators in supergroups and ignore slow mode
	CanPostMessages     bool        `json:"can_post_messages,omitempty"`      // Optional. Pass True if the administrator can create channel posts (channels only)
	CanEditMessages     bool        `json:"can_edit_messages,omitempty"`      // Optional. Pass True if the administrator can edit messages of other users and pin messages (channels only)
	CanDeleteMessages   bool        `json:"can_delete_messages,omitempty"`    // Optional. Pass True if the administrator can delete messages of other users
	CanManageVideoChats bool        `json:"can_manage_video_chats,omitempty"` // Optional. Pass True if the administrator can manage video chats
	CanRestrictMembers  bool        `json:"can_restrict_members,omitempty"`   // Optional. Pass True if the administrator can restrict, ban or unban chat members
	CanPromoteMembers   bool        `json:"can_promote_members,omitempty"`    // Optional. Pass True if the administrator can add new administrators with a subset of their own privileges or demote administrators that they have promoted, directly or indirectly
	CanChangeInfo       bool        `json:"can_change_info,omitempty"`        // Optional. Pass True if the administrator can change chat title, photo, and other settings
	CanInviteUsers      bool        `json:"can_invite_users,omitempty"`       // Optional. Pass True if the administrator can invite new users to the chat
	CanPinMessages      bool        `json:"can_pin_messages,omitempty"`       // Optional. Pass True if the administrator can pin messages (supergroups only)
	CanManageTopics     bool        `json:"can_manage_topics,omitempty"`      // Optional. Pass True if the user is allowed to create, rename, close, and reopen forum topics (supergroups only)
}

func (c PromoteChatMemberConf) method() string {
	return "promoteChatMember"
}

// SetChatAdministratorCustomTitleConf contains fields for the setChatAdministratorCustomTitle method. Returns True on success.
type SetChatAdministratorCustomTitleConf struct {
	ChatID      interface{} `json:"chat_id"`      // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	UserID      int         `json:"user_id"`      // Unique identifier of the target user
	CustomTitle string      `json:"custom_title"` // New custom title for the administrator; 0-16 characters, emoji are not allowed
}

func (c SetChatAdministratorCustomTitleConf) method() string {
	return "setChatAdministratorCustomTitle"
}

// BanChatSenderChatConf contains fields for the banChatSenderChat method. Returns True on success.
type BanChatSenderChatConf struct {
	ChatID       interface{} `json:"chat_id"`        // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	SenderChatID int         `json:"sender_chat_id"` // Unique identifier of the target sender chat
}

func (c BanChatSenderChatConf) method() string {
	return "banChatSenderChat"
}

// UnbanChatSenderChatConf contains fields for the unbanChatSenderChat method. Returns True on success.
type UnbanChatSenderChatConf struct {
	ChatID       interface{} `json:"chat_id"`        // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	SenderChatID int         `json:"sender_chat_id"` // Unique identifier of the target sender chat
}

func (c UnbanChatSenderChatConf) method() string {
	return "unbanChatSenderChat"
}

// SetChatPermissionsConf contains fields for the setChatPermissions method. Returns True on success.
type SetChatPermissionsConf struct {
	ChatID              interface{}     `json:"chat_id"`                                    // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	Permissions         ChatPermissions `json:"permissions"`                                // A JSON-serialized object for new default chat permissions
	UseIndependentPerms bool            `json:"use_independent_chat_permissions,omitempty"` // Optional. Pass True if chat permissions are set independently
}

func (c SetChatPermissionsConf) method() string {
	return "setChatPermissions"
}

// ExportChatInviteLinkConf contains fields for the exportChatInviteLink method. Returns the new invite link as String on success.
type ExportChatInviteLinkConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
}

func (c ExportChatInviteLinkConf) method() string {
	return "exportChatInviteLink"
}

// CreateChatInviteLinkConf contains fields for the createChatInviteLink method. Returns the new invite link as ChatInviteLink object.
type CreateChatInviteLinkConf struct {
	ChatID         interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	Name           string      `json:"name,omitempty"`                 // Optional. Invite link name; 0-32 characters
	ExpireDate     int         `json:"expire_date,omitempty"`          // Optional. Point in time (Unix timestamp) when the link will expire
	MemberLimit    int         `json:"member_limit,omitempty"`         // Optional. The maximum number of users that can be members of the chat simultaneously after joining the chat via this invite link; 1-99999
	CreatesJoinReq bool        `json:"creates_join_request,omitempty"` // Optional. True, if users joining the chat via the link need to be approved by chat administrators. If True, member_limit can't be specified
}

func (c CreateChatInviteLinkConf) method() string {
	return "createChatInviteLink"
}

// EditChatInviteLinkConf contains fields for the editChatInviteLink method. Returns the edited invite link as a ChatInviteLink object.
type EditChatInviteLinkConf struct {
	ChatID         interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	InviteLink     string      `json:"invite_link"`                    // The invite link to edit
	Name           string      `json:"name,omitempty"`                 // Optional. Invite link name; 0-32 characters
	ExpireDate     int         `json:"expire_date,omitempty"`          // Optional. Point in time (Unix timestamp) when the link will expire
	MemberLimit    int         `json:"member_limit,omitempty"`         // Optional. The maximum number of users that can be members of the chat simultaneously after joining the chat via this invite link; 1-99999
	CreatesJoinReq bool        `json:"creates_join_request,omitempty"` // Optional. True, if users joining the chat via the link need to be approved by chat administrators. If True, member_limit can't be specified
}

func (c EditChatInviteLinkConf) method() string {
	return "editChatInviteLink"
}

// RevokeChatInviteLinkConf contains fields for the revokeChatInviteLink method. Returns the revoked invite link as ChatInviteLink object.
type RevokeChatInviteLinkConf struct {
	ChatID     interface{} `json:"chat_id"`     // Unique identifier of the target chat or username of the target channel (in the format @channelusername)
	InviteLink string      `json:"invite_link"` // The invite link to revoke
}

func (c RevokeChatInviteLinkConf) method() string {
	return "revokeChatInviteLink"
}

// ApproveChatJoinRequestConf contains fields for the approveChatJoinRequest method. Returns True on success.
type ApproveChatJoinRequestConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	UserID int         `json:"user_id"` // Unique identifier of the target user
}

func (c ApproveChatJoinRequestConf) method() string {
	return "approveChatJoinRequest"
}

// DeclineChatJoinRequestConf contains fields for the declineChatJoinRequest method. Returns True on success.
type DeclineChatJoinRequestConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	UserID int         `json:"user_id"` // Unique identifier of the target user
}

func (c DeclineChatJoinRequestConf) method() string {
	return "declineChatJoinRequest"
}

// SetChatPhotoConf contains fields for the setChatPhoto method. Returns True on success.
type SetChatPhotoConf struct {
	ChatID interface{}     `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	File   RequestFileData `json:"photo"`   // New chat photo, uploaded using multipart/form-data
}

func (c SetChatPhotoConf) method() string {
	return "setChatPhoto"
}

func (config *SetChatPhotoConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "photo",
		Data: config.File,
	}}

	return files
}

// DeleteChatPhotoConf contains fields for the deleteChatPhoto method. Returns True on success.
type DeleteChatPhotoConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
}

func (c DeleteChatPhotoConf) method() string {
	return "deleteChatPhoto"
}

// SetChatTitleConf contains fields for the setChatTitle method. Returns True on success.
type SetChatTitleConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	Title  string      `json:"title"`   // New chat title, 1-128 characters
}

func (c SetChatTitleConf) method() string {
	return "setChatTitle"
}

// SetChatDescriptionConf contains fields for the setChatDescription method. Returns True on success.
type SetChatDescriptionConf struct {
	ChatID      interface{} `json:"chat_id"`               // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	Description string      `json:"description,omitempty"` // Optional. New chat description, 0-255 characters
}

func (c SetChatDescriptionConf) method() string {
	return "setChatDescription"
}

// PinChatMessageConf contains fields for the pinChatMessage method. Returns True on success.
type PinChatMessageConf struct {
	ChatID              interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	MessageID           int         `json:"message_id"`                     // Identifier of a message to pin
	DisableNotification bool        `json:"disable_notification,omitempty"` // Optional. Pass True if it is not necessary to send a notification to all chat members about the new pinned message. Notifications are always disabled in channels and private chats.
}

func (c PinChatMessageConf) method() string {
	return "pinChatMessage"
}

// UnpinChatMessageConf contains fields for the unpinChatMessage method. Returns True on success.
type UnpinChatMessageConf struct {
	ChatID    interface{} `json:"chat_id"`              // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	MessageID int         `json:"message_id,omitempty"` // Optional. Identifier of a message to unpin. If not specified, the most recent pinned message (by sending date) will be unpinned.
}

func (c UnpinChatMessageConf) method() string {
	return "unpinChatMessage"
}

// UnpinAllChatMessagesConf contains fields for the unpinAllChatMessages method. Returns True on success.
type UnpinAllChatMessagesConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
}

func (c UnpinAllChatMessagesConf) method() string {
	return "unpinAllChatMessages"
}

// LeaveChatConf contains fields for the leaveChat method. Returns True on success.
type LeaveChatConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup or channel (in the format @channelusername)
}

func (c LeaveChatConf) method() string {
	return "leaveChat"
}

// GetChatConf contains fields for the getChat method. Returns a Chat object on success.
type GetChatConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup or channel
}

func (c GetChatConf) method() string {
	return "getChat"
}

// GetChatAdministratorsConf contains fields for the getChatAdministrators method. Returns an Array of ChatMember objects.
type GetChatAdministratorsConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup or channel
}

func (c GetChatAdministratorsConf) method() string {
	return "getChatAdministrators"
}

// GetChatMemberCountConf contains fields for the getChatMemberCount method. Returns Int on success.
type GetChatMemberCountConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup or channel
}

func (c GetChatMemberCountConf) method() string {
	return "getChatMemberCount"
}

// GetChatMemberConf contains fields for the getChatMember method. Returns a ChatMember object on success.
type GetChatMemberConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup or channel
	UserID int         `json:"user_id"` // Unique identifier of the target user
}

func (c GetChatMemberConf) method() string {
	return "getChatMember"
}

// SetChatStickerSetConf contains fields for the setChatStickerSet method. Returns True on success.
type SetChatStickerSetConf struct {
	ChatID         interface{} `json:"chat_id"`          // Unique identifier for the target chat or username of the target supergroup
	StickerSetName string      `json:"sticker_set_name"` // Name of the sticker set to be set as the group sticker set
}

func (c SetChatStickerSetConf) method() string {
	return "setChatStickerSet"
}

// DeleteChatStickerSetConf contains fields for the deleteChatStickerSet method. Returns True on success.
type DeleteChatStickerSetConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup
}

func (c DeleteChatStickerSetConf) method() string {
	return "deleteChatStickerSet"
}

// CreateForumTopicConf contains fields for the createForumTopic method. Returns information about the created topic as a ForumTopic object.
type CreateForumTopicConf struct {
	ChatID            interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target supergroup
	Name              string      `json:"name"`                           // Topic name, 1-128 characters
	IconColor         int         `json:"icon_color,omitempty"`           // Optional. Color of the topic icon in RGB format
	IconCustomEmojiID string      `json:"icon_custom_emoji_id,omitempty"` // Optional. Unique identifier of the custom emoji shown as the topic icon
}

func (c CreateForumTopicConf) method() string {
	return "createForumTopic"
}

// EditForumTopicConf contains fields for the editForumTopic method. Returns True on success.
type EditForumTopicConf struct {
	ChatID          interface{} `json:"chat_id"`                        // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	MessageThreadID int         `json:"message_thread_id"`              // Unique identifier for the target message thread of the forum topic
	Name            string      `json:"name,omitempty"`                 // Optional. New topic name, 0-128 characters. If not specified or empty, the current name of the topic will be kept
	IconCustomEmoji string      `json:"icon_custom_emoji_id,omitempty"` // Optional. New unique identifier of the custom emoji shown as the topic icon. Pass an empty string to remove the icon. If not specified, the current icon will be kept
}

func (c EditForumTopicConf) method() string {
	return "editForumTopic"
}

// CloseForumTopicConf contains fields for the closeForumTopic method. Returns True on success.
type CloseForumTopicConf struct {
	ChatID          interface{} `json:"chat_id"`           // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	MessageThreadID int         `json:"message_thread_id"` // Unique identifier for the target message thread of the forum topic
}

func (c CloseForumTopicConf) method() string {
	return "closeForumTopic"
}

// ReopenForumTopicConf contains fields for the reopenForumTopic method. Returns True on success.
type ReopenForumTopicConf struct {
	ChatID          interface{} `json:"chat_id"`           // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	MessageThreadID int         `json:"message_thread_id"` // Unique identifier for the target message thread of the forum topic
}

func (c ReopenForumTopicConf) method() string {
	return "reopenForumTopic"
}

// DeleteForumTopicConf contains fields for the deleteForumTopic method. Returns True on success.
type DeleteForumTopicConf struct {
	ChatID          interface{} `json:"chat_id"`           // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	MessageThreadID int         `json:"message_thread_id"` // Unique identifier for the target message thread of the forum topic
}

func (c DeleteForumTopicConf) method() string {
	return "deleteForumTopic"
}

// UnpinAllForumTopicMessagesConf contains fields for the unpinAllForumTopicMessages method. Returns True on success.
type UnpinAllForumTopicMessagesConf struct {
	ChatID          interface{} `json:"chat_id"`           // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	MessageThreadID int         `json:"message_thread_id"` // Unique identifier for the target message thread of the forum topic
}

func (c UnpinAllForumTopicMessagesConf) method() string {
	return "unpinAllForumTopicMessages"
}

// EditGeneralForumTopicConf contains fields for the editGeneralForumTopic method. Returns True on success.
type EditGeneralForumTopicConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	Name   string      `json:"name"`    // New topic name, 1-128 characters
}

func (c EditGeneralForumTopicConf) method() string {
	return "editGeneralForumTopic"
}

// CloseGeneralForumTopicConf contains fields for the closeGeneralForumTopic method. Returns True on success.
type CloseGeneralForumTopicConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
}

func (c CloseGeneralForumTopicConf) method() string {
	return "closeGeneralForumTopic"
}

// ReopenGeneralForumTopicConf contains fields for the reopenGeneralForumTopic method. Returns True on success.
type ReopenGeneralForumTopicConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
}

func (c ReopenGeneralForumTopicConf) method() string {
	return "reopenGeneralForumTopic"
}

// HideGeneralForumTopicConf contains fields for the hideGeneralForumTopic method. Returns True on success.
type HideGeneralForumTopicConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
}

func (c HideGeneralForumTopicConf) method() string {
	return "hideGeneralForumTopic"
}

// UnhideGeneralForumTopicConf contains fields for the unhideGeneralForumTopic method. Returns True on success.
type UnhideGeneralForumTopicConf struct {
	ChatID interface{} `json:"chat_id"` // Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
}

func (c UnhideGeneralForumTopicConf) method() string {
	return "unhideGeneralForumTopic"
}

// AnswerCallbackQueryConf contains fields for the answerCallbackQuery method. Returns True on success.
type AnswerCallbackQueryConf struct {
	CallbackQueryID string `json:"callback_query_id"`    // Unique identifier for the query to be answered
	Text            string `json:"text,omitempty"`       // Optional. Text of the notification. If not specified, nothing will be shown to the user, 0-200 characters
	ShowAlert       bool   `json:"show_alert,omitempty"` // Optional. If True, an alert will be shown by the client instead of a notification at the top of the chat screen. Defaults to false.
	URL             string `json:"url,omitempty"`        // Optional. URL that will be opened by the user's client
	CacheTime       int    `json:"cache_time,omitempty"` // Optional. The maximum amount of time in seconds that the result of the callback query may be cached client-side. Defaults to 0.
}

func (c AnswerCallbackQueryConf) method() string {
	return "answerCallbackQuery"
}

// SetMyCommandsConf contains fields for the setMyCommands method. Returns True on success.
type SetMyCommandsConf struct {
	Commands     []BotCommand     `json:"commands"`                // A JSON-serialized list of bot commands to be set as the list of the bot's commands. At most 100 commands can be specified.
	Scope        *BotCommandScope `json:"scope,omitempty"`         // Optional. A JSON-serialized object describing the scope of users for which the commands are relevant. Defaults to BotCommandScopeDefault.
	LanguageCode string           `json:"language_code,omitempty"` // Optional. A two-letter ISO 639-1 language code. If empty, commands will be applied to all users from the given scope, for whom there are no dedicated commands.
}

func (c SetMyCommandsConf) method() string {
	return "setMyCommands"
}

// DeleteMyCommandsConf contains fields for the deleteMyCommands method. Returns True on success.
type DeleteMyCommandsConf struct {
	Scope        *BotCommandScope `json:"scope,omitempty"`         // Optional. Scope of users for which the commands are relevant
	LanguageCode string           `json:"language_code,omitempty"` // Optional. Language code for which the commands are relevant
}

func (c DeleteMyCommandsConf) method() string {
	return "deleteMyCommands"
}

// GetMyCommandsConf contains fields for the getMyCommands method. Returns an Array of BotCommand objects. If commands aren't set, an empty list is returned.
type GetMyCommandsConf struct {
	Scope        *BotCommandScope `json:"scope,omitempty"`         // Optional. Scope of users
	LanguageCode string           `json:"language_code,omitempty"` // Optional. Language code for which the commands are relevant
}

func (c GetMyCommandsConf) method() string {
	return "getMyCommands"
}

// SetMyNameConf contains fields for the setMyName method. Returns True on success.
type SetMyNameConf struct {
	Name         string `json:"name,omitempty"`          // Optional. New bot name
	LanguageCode string `json:"language_code,omitempty"` // Optional. Language code for which the name is relevant
}

func (c SetMyNameConf) method() string {
	return "setMyName"
}

// GetMyNameConf contains fields for the getMyName method. Returns BotName on success.
type GetMyNameConf struct {
	LanguageCode string `json:"language_code,omitempty"` // Optional. Language code for which the name is relevant
}

func (c GetMyNameConf) method() string {
	return "getMyName"
}

// SetMyDescriptionConf contains fields for the setMyDescription method. Returns True on success.
type SetMyDescriptionConf struct {
	Description  string `json:"description,omitempty"`   // Optional. New bot description
	LanguageCode string `json:"language_code,omitempty"` // Optional. Language code for which the description is relevant
}

func (c SetMyDescriptionConf) method() string {
	return "setMyDescription"
}

// GetMyDescriptionConf contains fields for the getMyDescription method. Returns BotDescription on success.
type GetMyDescriptionConf struct {
	LanguageCode string `json:"language_code,omitempty"` // Optional. Language code for which the description is relevant
}

func (c GetMyDescriptionConf) method() string {
	return "getMyDescription"
}

// SetMyShortDescriptionConf contains fields for the setMyShortDescription method. Returns True on success.
type SetMyShortDescriptionConf struct {
	ShortDescription string `json:"short_description,omitempty"` // Optional. New short description for the bot
	LanguageCode     string `json:"language_code,omitempty"`     // Optional. Language code for which the short description is relevant
}

func (c SetMyShortDescriptionConf) method() string {
	return "setMyShortDescription"
}

// GetMyShortDescriptionConf contains fields for the getMyShortDescription method. Returns BotShortDescription on success.
type GetMyShortDescriptionConf struct {
	LanguageCode string `json:"language_code,omitempty"` // Optional. Language code for which the short description is relevant
}

func (c GetMyShortDescriptionConf) method() string {
	return "getMyShortDescription"
}

// SetChatMenuButtonConf contains fields for the setChatMenuButton method. Returns True on success.
type SetChatMenuButtonConf struct {
	ChatID     int         `json:"chat_id,omitempty"`     // Optional. Target private chat ID
	MenuButton *MenuButton `json:"menu_button,omitempty"` // Optional. New menu button for the bot
}

func (c SetChatMenuButtonConf) method() string {
	return "setChatMenuButton"
}

// GetChatMenuButtonConf contains fields for the getChatMenuButton method. Returns MenuButton on success.
type GetChatMenuButtonConf struct {
	ChatID int `json:"chat_id,omitempty"` // Optional. Target private chat ID
}

func (c GetChatMenuButtonConf) method() string {
	return "getChatMenuButton"
}

// SetMyDefaultAdministratorRightsConf contains fields for the setMyDefaultAdministratorRights method. Returns True on success.
type SetMyDefaultAdministratorRightsConf struct {
	Rights      *ChatAdministratorRights `json:"rights,omitempty"`       // Optional. New default administrator rights
	ForChannels bool                     `json:"for_channels,omitempty"` // Optional. Change the default administrator rights for channels
}

func (c SetMyDefaultAdministratorRightsConf) method() string {
	return "setMyDefaultAdministratorRights"
}

// GetMyDefaultAdministratorRightsConf contains fields for the getMyDefaultAdministratorRights method. Returns ChatAdministratorRights on success.
type GetMyDefaultAdministratorRightsConf struct {
	ForChannels bool `json:"for_channels,omitempty"` // Optional. Get the default administrator rights for channels
}

func (c GetMyDefaultAdministratorRightsConf) method() string {
	return "getMyDefaultAdministratorRights"
}

//
//
//
// Messages
//
//
//

// EditMessageTextConf contains fields for the editMessageText method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type EditMessageTextConf struct {
	ChatID                interface{}           `json:"chat_id,omitempty"`                  // Optional. Unique identifier for the target chat or username of the target channel
	MessageID             int                   `json:"message_id,omitempty"`               // Optional. Identifier of the message to edit
	InlineMessageID       string                `json:"inline_message_id,omitempty"`        // Optional. Identifier of the inline message
	Text                  string                `json:"text"`                               // New text of the message
	ParseMode             string                `json:"parse_mode,omitempty"`               // Optional. Mode for parsing entities in the message text
	Entities              []MessageEntity       `json:"entities,omitempty"`                 // Optional. List of special entities that appear in the message text
	DisableWebPagePreview bool                  `json:"disable_web_page_preview,omitempty"` // Optional. Disables link previews for links in this message
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`             // Optional. Inline keyboard markup
}

func (c EditMessageTextConf) method() string {
	return "editMessageText"
}

// EditMessageCaptionConf contains fields for the editMessageCaption method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type EditMessageCaptionConf struct {
	ChatID          interface{}           `json:"chat_id,omitempty"`           // Optional. Unique identifier for the target chat or username of the target channel
	MessageID       int                   `json:"message_id,omitempty"`        // Optional. Identifier of the message to edit
	InlineMessageID string                `json:"inline_message_id,omitempty"` // Optional. Identifier of the inline message
	Caption         string                `json:"caption,omitempty"`           // Optional. New caption of the message
	ParseMode       string                `json:"parse_mode,omitempty"`        // Optional. Mode for parsing entities in the caption
	CaptionEntities []MessageEntity       `json:"caption_entities,omitempty"`  // Optional. List of special entities that appear in the caption
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // Optional. Inline keyboard markup
}

func (c EditMessageCaptionConf) method() string {
	return "editMessageCaption"
}

// EditMessageMediaConf contains fields for the editMessageMedia method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type EditMessageMediaConf struct {
	ChatID          interface{}           `json:"chat_id,omitempty"`           // Optional. Unique identifier for the target chat or username of the target channel
	MessageID       int                   `json:"message_id,omitempty"`        // Optional. Identifier of the message to edit
	InlineMessageID string                `json:"inline_message_id,omitempty"` // Optional. Identifier of the inline message
	Media           interface{}           `json:"media"`                       // A new media content of the message
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // Optional. Inline keyboard markup
}

func (c EditMessageMediaConf) method() string {
	return "editMessageMedia"
}

// EditMessageLiveLocationConf contains fields for the editMessageLiveLocation method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type EditMessageLiveLocationConf struct {
	ChatID               interface{}           `json:"chat_id,omitempty"`                // Optional. Unique identifier for the target chat or username of the target channel
	MessageID            int                   `json:"message_id,omitempty"`             // Optional. Identifier of the message to edit
	InlineMessageID      string                `json:"inline_message_id,omitempty"`      // Optional. Identifier of the inline message
	Latitude             float64               `json:"latitude"`                         // Latitude of the new location
	Longitude            float64               `json:"longitude"`                        // Longitude of the new location
	HorizontalAccuracy   float64               `json:"horizontal_accuracy,omitempty"`    // Optional. The radius of uncertainty for the location, measured in meters
	Heading              int                   `json:"heading,omitempty"`                // Optional. Direction in which the user is moving, in degrees
	ProximityAlertRadius int                   `json:"proximity_alert_radius,omitempty"` // Optional. The maximum distance for proximity alerts about approaching another chat member, in meters
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`           // Optional. Inline keyboard markup
}

func (c EditMessageLiveLocationConf) method() string {
	return "editMessageLiveLocation"
}

// StopMessageLiveLocationConf contains fields for the stopMessageLiveLocation method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type StopMessageLiveLocationConf struct {
	ChatID          interface{}           `json:"chat_id,omitempty"`           // Optional. Unique identifier for the target chat or username of the target channel
	MessageID       int                   `json:"message_id,omitempty"`        // Optional. Identifier of the message with live location to stop
	InlineMessageID string                `json:"inline_message_id,omitempty"` // Optional. Identifier of the inline message
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // Optional. Inline keyboard markup
}

func (c StopMessageLiveLocationConf) method() string {
	return "stopMessageLiveLocation"
}

// EditMessageReplyMarkupConf contains fields for the editMessageReplyMarkup method. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.
type EditMessageReplyMarkupConf struct {
	ChatID          interface{}           `json:"chat_id,omitempty"`           // Optional. Unique identifier for the target chat or username of the target channel
	MessageID       int                   `json:"message_id,omitempty"`        // Optional. Identifier of the message to edit
	InlineMessageID string                `json:"inline_message_id,omitempty"` // Optional. Identifier of the inline message
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`      // Optional. Inline keyboard markup
}

func (c EditMessageReplyMarkupConf) method() string {
	return "editMessageReplyMarkup"
}

// StopPollConf contains fields for the stopPoll method. On success, the stopped Poll is returned.
type StopPollConf struct {
	ChatID      interface{}           `json:"chat_id"`                // Unique identifier for the target chat or username of the target channel
	MessageID   int                   `json:"message_id"`             // Identifier of the original message with the poll
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"` // Optional. Inline keyboard markup for a new message
}

func (c StopPollConf) method() string {
	return "stopPoll"
}

// DeleteMessageConf contains fields for the deleteMessage method. Returns True on success.
type DeleteMessageConf struct {
	ChatID    interface{} `json:"chat_id"`    // Unique identifier for the target chat or username of the target channel
	MessageID int         `json:"message_id"` // Identifier of the message to delete
}

func (c DeleteMessageConf) method() string {
	return "deleteMessage"
}

//
//
//
// Stickers
//
//
//

// SendStickerConf contains fields for the sendSticker method. On success, the sent Message is returned.
type SendStickerConf struct {
	BaseSend                 // Unique identifier for the target chat or username of the target channel
	File     RequestFileData `json:"sticker"`         // Sticker to send
	Emoji    string          `json:"emoji,omitempty"` // Optional. Emoji associated with the sticker
}

func (c SendStickerConf) method() string {
	return "sendSticker"
}

func (config *SendStickerConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "sticker",
		Data: config.File,
	}}

	return files
}

// GetStickerSetConf contains fields for the getStickerSet method. On success, a StickerSet object is returned.
type GetStickerSetConf struct {
	Name string `json:"name"` // Name of the sticker set
}

func (c GetStickerSetConf) method() string {
	return "getStickerSet"
}

// GetCustomEmojiStickersConf contains fields for the getCustomEmojiStickers method. Returns an Array of Sticker objects.
type GetCustomEmojiStickersConf struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"` // List of custom emoji identifiers
}

func (c GetCustomEmojiStickersConf) method() string {
	return "getCustomEmojiStickers"
}

// UploadStickerFileConf contains fields for the uploadStickerFile method. Returns the uploaded File on success.
type UploadStickerFileConf struct {
	UserID        int             `json:"user_id"`        // User identifier of sticker file owner
	File          RequestFileData `json:"sticker"`        // A file with the sticker
	StickerFormat string          `json:"sticker_format"` // Format of the sticker
}

func (c UploadStickerFileConf) method() string {
	return "uploadStickerFile"
}

func (config *UploadStickerFileConf) files() []RequestFile {
	files := []RequestFile{{
		Name: "sticker",
		Data: config.File,
	}}

	return files
}

// CreateNewStickerSetConf contains fields for the createNewStickerSet method. Returns True on success.
type CreateNewStickerSetConf struct {
	UserID          int            `json:"user_id"`                    // User identifier of created sticker set owner
	Name            string         `json:"name"`                       // Short name of sticker set
	Title           string         `json:"title"`                      // Sticker set title
	Stickers        []InputSticker `json:"stickers"`                   // List of initial stickers to be added to the sticker set
	StickerFormat   string         `json:"sticker_format"`             // Format of stickers in the set
	StickerType     string         `json:"sticker_type,omitempty"`     // Optional. Type of stickers in the set
	NeedsRepainting bool           `json:"needs_repainting,omitempty"` // Optional. Pass True if stickers in the sticker set must be repainted based on context
}

func (c CreateNewStickerSetConf) method() string {
	return "createNewStickerSet"
}

// AddStickerToSetConf contains fields for the addStickerToSet method. Returns True on success.
type AddStickerToSetConf struct {
	UserID  int          `json:"user_id"` // User identifier of sticker set owner
	Name    string       `json:"name"`    // Sticker set name
	Sticker InputSticker `json:"sticker"` // Information about the added sticker
}

func (c AddStickerToSetConf) method() string {
	return "addStickerToSet"
}

// SetStickerPositionInSetConf contains fields for the setStickerPositionInSet method. Returns True on success.
type SetStickerPositionInSetConf struct {
	Sticker  string `json:"sticker"`  // File identifier of the sticker
	Position int    `json:"position"` // New sticker position in the set, zero-based
}

func (c SetStickerPositionInSetConf) method() string {
	return "setStickerPositionInSet"
}

// DeleteStickerFromSetConf contains fields for the deleteStickerFromSet method. Returns True on success.
type DeleteStickerFromSetConf struct {
	Sticker string `json:"sticker"` // File identifier of the sticker
}

func (c DeleteStickerFromSetConf) method() string {
	return "deleteStickerFromSet"
}

// SetStickerEmojiListConf contains fields for the setStickerEmojiList method. Returns True on success.
type SetStickerEmojiListConf struct {
	Sticker   string   `json:"sticker"`    // File identifier of the sticker
	EmojiList []string `json:"emoji_list"` // List of emoji associated with the sticker
}

func (c SetStickerEmojiListConf) method() string {
	return "setStickerEmojiList"
}

// SetStickerKeywordsConf contains fields for the setStickerKeywords method. Returns True on success.
type SetStickerKeywordsConf struct {
	Sticker  string   `json:"sticker"`            // File identifier of the sticker
	Keywords []string `json:"keywords,omitempty"` // Optional. List of search keywords for the sticker
}

func (c SetStickerKeywordsConf) method() string {
	return "setStickerKeywords"
}

// SetStickerMaskPositionConf contains fields for the setStickerMaskPosition method. Returns True on success.
type SetStickerMaskPositionConf struct {
	Sticker      string       `json:"sticker"`                 // File identifier of the sticker
	MaskPosition MaskPosition `json:"mask_position,omitempty"` // Optional. Mask position on faces
}

func (c SetStickerMaskPositionConf) method() string {
	return "setStickerMaskPosition"
}

// SetStickerSetTitleConf contains fields for the setStickerSetTitle method. Returns True on success.
type SetStickerSetTitleConf struct {
	Name  string `json:"name"`  // Sticker set name
	Title string `json:"title"` // Sticker set title
}

func (c SetStickerSetTitleConf) method() string {
	return "setStickerSetTitle"
}

// SetStickerSetThumbnailConf contains fields for the setStickerSetThumbnail method. Returns True on success.
type SetStickerSetThumbnailConf struct {
	Name      string          `json:"name"`                // Sticker set name
	UserID    int             `json:"user_id"`             // User identifier of the sticker set owner
	Thumbnail RequestFileData `json:"thumbnail,omitempty"` // Optional. Thumbnail image or animation
}

func (c SetStickerSetThumbnailConf) method() string {
	return "setStickerSetThumbnail"
}

func (config *SetStickerSetThumbnailConf) files() []RequestFile {
	if config.Thumbnail == nil {
		return nil
	}
	files := []RequestFile{{
		Name: "thumbnail",
		Data: config.Thumbnail,
	}}

	return files
}

// SetCustomEmojiStickerSetThumbnailConf contains fields for the setCustomEmojiStickerSetThumbnail method. Returns True on success.
type SetCustomEmojiStickerSetThumbnailConf struct {
	Name          string `json:"name"`                      // Sticker set name
	CustomEmojiID string `json:"custom_emoji_id,omitempty"` // Optional. Custom emoji identifier of a sticker from the set
}

func (c SetCustomEmojiStickerSetThumbnailConf) method() string {
	return "setCustomEmojiStickerSetThumbnail"
}

// DeleteStickerSetConf contains fields for the deleteStickerSet method.
type DeleteStickerSetConf struct {
	Name string `json:"name"` // Sticker set name
}

func (c DeleteStickerSetConf) method() string {
	return "deleteStickerSet"
}

//
//
//
// Inline
//
//
//

// AnswerInlineQueryConf contains fields for the answerInlineQuery method. On success, True is returned. No more than 50 results per query are allowed.
type AnswerInlineQueryConf struct {
	InlineQueryID string                    `json:"inline_query_id"` // Unique identifier for the answered query
	Result        interface{}               `json:"result"`          // A JSON-serialized array of results for the inline query
	CacheTime     int                       `json:"cache_time"`      // Optional. The maximum amount of time in seconds that the result of the inline query may be cached on the server. Defaults to 300.
	IsPersonal    bool                      `json:"is_personal"`     // Optional. Pass True if results may be cached on the server side only for the user that sent the query. By default, results may be returned to any user who sends the same query.
	NextOffset    string                    `json:"next_offset"`     // Optional. Pass the offset that a client should send in the next query with the same text to receive more results. Pass an empty string if there are no more results or if you don't support pagination. Offset length can't exceed 64 bytes.
	Button        *InlineQueryResultsButton `json:"button"`          // Optional. A JSON-serialized object describing a button to be shown above inline query results
}

func (c AnswerInlineQueryConf) method() string {
	return "answerInlineQuery"
}

// AnswerWebAppQueryConf contains fields for the answerWebAppQuery method. On success, a SentWebAppMessage object is returned.
type AnswerWebAppQueryConf struct {
	WebAppQueryID string      `json:"web_app_query_id"` // Unique identifier for the query to be answered
	Result        interface{} `json:"result"`           // A JSON-serialized object describing the message to be sent
}

func (c AnswerWebAppQueryConf) method() string {
	return "answerWebAppQuery"
}

//
//
//
// Payment
//
//
//

// SendInvoiceConf contains fields for the sendInvoice method. On success, the sent Message is returned.
type SendInvoiceConf struct {
	ChatID                    interface{}           `json:"chat_id"`                                 // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	MessageThreadID           int                   `json:"message_thread_id,omitempty"`             // Optional. Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	Title                     string                `json:"title"`                                   // Product name, 1-32 characters
	Description               string                `json:"description"`                             // Product description, 1-255 characters
	Payload                   string                `json:"payload"`                                 // Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use for your internal processes.
	ProviderToken             string                `json:"provider_token"`                          // Payment provider token, obtained via @BotFather
	Currency                  string                `json:"currency"`                                // Three-letter ISO 4217 currency code, see more on currencies
	Prices                    []LabeledPrice        `json:"prices"`                                  // Price breakdown, a JSON-serialized list of components (e.g. product price, tax, discount, delivery cost, delivery tax, bonus, etc.)
	MaxTipAmount              int                   `json:"max_tip_amount,omitempty"`                // Optional. The maximum accepted amount for tips in the smallest units of the currency (integer, not float/double).
	SuggestedTipAmounts       []int                 `json:"suggested_tip_amounts,omitempty"`         // Optional. A JSON-serialized array of suggested amounts of tips in the smallest units of the currency (integer, not float/double).
	StartParameter            string                `json:"start_parameter,omitempty"`               // Optional. Unique deep-linking parameter.
	ProviderData              string                `json:"provider_data,omitempty"`                 // Optional. JSON-serialized data about the invoice, which will be shared with the payment provider.
	PhotoURL                  string                `json:"photo_url,omitempty"`                     // Optional. URL of the product photo for the invoice.
	PhotoSize                 int                   `json:"photo_size,omitempty"`                    // Optional. Photo size in bytes.
	PhotoWidth                int                   `json:"photo_width,omitempty"`                   // Optional. Photo width.
	PhotoHeight               int                   `json:"photo_height,omitempty"`                  // Optional. Photo height.
	NeedName                  bool                  `json:"need_name,omitempty"`                     // Optional. Pass True if you require the user's full name to complete the order.
	NeedPhoneNumber           bool                  `json:"need_phone_number,omitempty"`             // Optional. Pass True if you require the user's phone number to complete the order.
	NeedEmail                 bool                  `json:"need_email,omitempty"`                    // Optional. Pass True if you require the user's email address to complete the order.
	NeedShippingAddress       bool                  `json:"need_shipping_address,omitempty"`         // Optional. Pass True if you require the user's shipping address to complete the order.
	SendPhoneNumberToProvider bool                  `json:"send_phone_number_to_provider,omitempty"` // Optional. Pass True if the user's phone number should be sent to the provider.
	SendEmailToProvider       bool                  `json:"send_email_to_provider,omitempty"`        // Optional. Pass True if the user's email address should be sent to the provider.
	IsFlexible                bool                  `json:"is_flexible,omitempty"`                   // Optional. Pass True if the final price depends on the shipping method.
	DisableNotification       bool                  `json:"disable_notification,omitempty"`          // Optional. Sends the message silently. Users will receive a notification with no sound.
	ProtectContent            bool                  `json:"protect_content,omitempty"`               // Optional. Protects the contents of the sent message from forwarding and saving.
	ReplyToMessageID          int                   `json:"reply_to_message_id,omitempty"`           // Optional. If the message is a reply, ID of the original message.
	AllowSendingWithoutReply  bool                  `json:"allow_sending_without_reply,omitempty"`   // Optional. Pass True if the message should be sent even if the specified replied-to message is not found.
	ReplyMarkup               *InlineKeyboardMarkup `json:"reply_markup,omitempty"`                  // Optional. A JSON-serialized object for an inline keyboard.
}

func (c SendInvoiceConf) method() string {
	return "sendInvoice"
}

// CreateInvoiceLinkConf contains fields for the createInvoiceLink method. Returns the created invoice link as String on success.
type CreateInvoiceLinkConf struct {
	Title                     string         `json:"title"`                                   // Product name, 1-32 characters
	Description               string         `json:"description"`                             // Product description, 1-255 characters
	Payload                   string         `json:"payload"`                                 // Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use for your internal processes.
	ProviderToken             string         `json:"provider_token"`                          // Payment provider token, obtained via BotFather
	Currency                  string         `json:"currency"`                                // Three-letter ISO 4217 currency code, see more on currencies
	Prices                    []LabeledPrice `json:"prices"`                                  // Price breakdown, a JSON-serialized list of components (e.g. product price, tax, discount, delivery cost, delivery tax, bonus, etc.)
	MaxTipAmount              int            `json:"max_tip_amount,omitempty"`                // Optional. The maximum accepted amount for tips in the smallest units of the currency (integer, not float/double).
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts,omitempty"`         // Optional. A JSON-serialized array of suggested amounts of tips in the smallest units of the currency (integer, not float/double).
	ProviderData              string         `json:"provider_data,omitempty"`                 // Optional. JSON-serialized data about the invoice, which will be shared with the payment provider.
	PhotoURL                  string         `json:"photo_url,omitempty"`                     // Optional. URL of the product photo for the invoice.
	PhotoSize                 int            `json:"photo_size,omitempty"`                    // Optional. Photo size in bytes.
	PhotoWidth                int            `json:"photo_width,omitempty"`                   // Optional. Photo width.
	PhotoHeight               int            `json:"photo_height,omitempty"`                  // Optional. Photo height.
	NeedName                  bool           `json:"need_name,omitempty"`                     // Optional. Pass True if you require the user's full name to complete the order.
	NeedPhoneNumber           bool           `json:"need_phone_number,omitempty"`             // Optional. Pass True if you require the user's phone number to complete the order.
	NeedEmail                 bool           `json:"need_email,omitempty"`                    // Optional. Pass True if you require the user's email address to complete the order.
	NeedShippingAddress       bool           `json:"need_shipping_address,omitempty"`         // Optional. Pass True if you require the user's shipping address to complete the order.
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider,omitempty"` // Optional. Pass True if the user's phone number should be sent to the provider.
	SendEmailToProvider       bool           `json:"send_email_to_provider,omitempty"`        // Optional. Pass True if the user's email address should be sent to the provider.
	IsFlexible                bool           `json:"is_flexible,omitempty"`                   // Optional. Pass True if the final price depends on the shipping method.
}

func (c CreateInvoiceLinkConf) method() string {
	return "createInvoiceLink"
}

// AnswerShippingQueryConf contains fields for the answerShippingQuery method. On success, True is returned.
type AnswerShippingQueryConf struct {
	ShippingQueryID string           `json:"shipping_query_id"`          // Unique identifier for the query to be answered
	OK              bool             `json:"ok"`                         // Pass True if delivery to the specified address is possible and False if there are any problems.
	ShippingOptions []ShippingOption `json:"shipping_options,omitempty"` // Optional. A JSON-serialized array of available shipping options.
	ErrorMessage    string           `json:"error_message,omitempty"`    // Optional. Required if ok is False. Error message in human readable form that explains why it is impossible to complete the order.
}

func (c AnswerShippingQueryConf) method() string {
	return "answerShippingQuery"
}

// AnswerPreCheckoutQueryConf contains fields for the answerPreCheckoutQuery method. On success, True is returned.
type AnswerPreCheckoutQueryConf struct {
	PreCheckoutQueryID string `json:"pre_checkout_query_id"`   // Unique identifier for the query to be answered
	OK                 bool   `json:"ok"`                      // Specify True if everything is alright (goods are available, etc.) and the bot is ready to proceed with the order. Use False if there are any problems.
	ErrorMessage       string `json:"error_message,omitempty"` // Optional. Required if ok is False. Error message in human readable form that explains the reason for failure to proceed with the checkout.
}

func (c AnswerPreCheckoutQueryConf) method() string {
	return "answerPreCheckoutQuery"
}

//
//
//
// Passport
//
//
//

// SetPassportDataErrorsConf contains fields for the setPassportDataErrors method. Returns True on success.
type SetPassportDataErrorsConf struct {
	UserID int           `json:"user_id"` // User identifier
	Errors []interface{} `json:"errors"`  // A JSON-serialized array describing the errors
}

func (c SetPassportDataErrorsConf) method() string {
	return "setPassportDataErrors"
}

//
//
//
// Game
//
//
//

// SendGameConf contains fields for the sendGame method. On success, the sent Message is returned.
type SendGameConf struct {
	ChatID                   int                   `json:"chat_id"`                               // Unique identifier for the target chat
	MessageThreadID          int                   `json:"message_thread_id,omitempty"`           // Optional. Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	GameShortName            string                `json:"game_short_name"`                       // Short name of the game, serves as the unique identifier for the game
	DisableNotification      bool                  `json:"disable_notification,omitempty"`        // Optional. Sends the message silently. Users will receive a notification with no sound
	ProtectContent           bool                  `json:"protect_content,omitempty"`             // Optional. Protects the contents of the sent message from forwarding and saving
	ReplyToMessageID         int                   `json:"reply_to_message_id,omitempty"`         // Optional. If the message is a reply, ID of the original message
	AllowSendingWithoutReply bool                  `json:"allow_sending_without_reply,omitempty"` // Optional. Pass True if the message should be sent even if the specified replied-to message is not found
	ReplyMarkup              *InlineKeyboardMarkup `json:"reply_markup,omitempty"`                // Optional. A JSON-serialized object for an inline keyboard. If empty, one 'Play game_title' button will be shown. If not empty, the first button must launch the game.
}

func (c SendGameConf) method() string {
	return "sendGame"
}

// SetGameScoreConf contains fields for the setGameScore method. On success, if the message is not an inline message, the Message is returned, otherwise True is returned. Returns an error, if the new score is not greater than the user's current score in the chat and force is False.
type SetGameScoreConf struct {
	UserID             int    `json:"user_id"`                        // User identifier
	Score              int    `json:"score"`                          // New score, must be non-negative
	Force              bool   `json:"force,omitempty"`                // Optional. Pass True if the high score is allowed to decrease
	DisableEditMessage bool   `json:"disable_edit_message,omitempty"` // Optional. Pass True if the game message should not be automatically edited to include the current scoreboard
	ChatID             int    `json:"chat_id,omitempty"`              // Optional. Required if inline_message_id is not specified. Unique identifier for the target chat
	MessageID          int    `json:"message_id,omitempty"`           // Optional. Required if inline_message_id is not specified. Identifier of the sent message
	InlineMessageID    string `json:"inline_message_id,omitempty"`    // Optional. Required if chat_id and message_id are not specified. Identifier of the inline message
}

func (c SetGameScoreConf) method() string {
	return "setGameScore"
}

// GetGameHighScoresConf contains fields for the getGameHighScores method. Returns an Array of GameHighScore objects.
type GetGameHighScoresConf struct {
	UserID          int    `json:"user_id"`                     // Target user id
	ChatID          int    `json:"chat_id,omitempty"`           // Optional. Required if inline_message_id is not specified. Unique identifier for the target chat
	MessageID       int    `json:"message_id,omitempty"`        // Optional. Required if inline_message_id is not specified. Identifier of the sent message
	InlineMessageID string `json:"inline_message_id,omitempty"` // Optional. Required if chat_id and message_id are not specified. Identifier of the inline message
}

func (c GetGameHighScoresConf) method() string {
	return "getGameHighScores"
}
