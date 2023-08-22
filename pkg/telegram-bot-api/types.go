package telegram

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result,omitempty"`
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// Error is an error containing extra information returned by the Telegram API.
type Error struct {
	Code    int
	Message string
	ResponseParameters
}

// Error message string.
func (e Error) Error() string {
	return e.Message
}

//
//
//
// Update types
//
//
//

// This object represents an incoming update.
// At most one of the optional parameters can be present in any given update.
type Update struct {
	UpdateID           int                 `json:"update_id"`                      // The update's unique identifier
	Message            *Message            `json:"message,omitempty"`              // Optional. New incoming message
	EditedMessage      *Message            `json:"edited_message,omitempty"`       // Optional. New version of a message that was edited
	ChannelPost        *Message            `json:"channel_post,omitempty"`         // Optional. New incoming channel post
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`  // Optional. New version of a channel post that was edited
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`         // Optional. New incoming inline query
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"` // Optional. Result of an inline query chosen by a user
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`       // Optional. New incoming callback query
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`       // Optional. New incoming shipping query
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`   // Optional. New incoming pre-checkout query
	Poll               *Poll               `json:"poll,omitempty"`                 // Optional. New poll state
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`          // Optional. User changed their answer in a non-anonymous poll
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member,omitempty"`       // Optional. Bot's chat member status was updated in a chat
	ChatMember         *ChatMemberUpdated  `json:"chat_member,omitempty"`          // Optional. Chat member's status was updated in a chat
	ChatJoinRequest    *ChatJoinRequest    `json:"chat_join_request,omitempty"`    // Optional. Request to join the chat has been sent
}

// SentFrom returns the user who sent an update. Can be nil, if Telegram did not provide information
// about the user in the update object.
func (u *Update) SentFrom() *User {
	switch {
	case u.Message != nil:
		return u.Message.From
	case u.EditedMessage != nil:
		return u.EditedMessage.From
	case u.InlineQuery != nil:
		return u.InlineQuery.From
	case u.ChosenInlineResult != nil:
		return u.ChosenInlineResult.From
	case u.CallbackQuery != nil:
		return u.CallbackQuery.From
	case u.ShippingQuery != nil:
		return u.ShippingQuery.From
	case u.PreCheckoutQuery != nil:
		return u.PreCheckoutQuery.From
	default:
		return nil
	}
}

// CallbackData returns the callback query data, if it exists.
func (u *Update) CallbackData() string {
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	}
	return ""
}

// FromChat returns the chat where an update occurred.
func (u *Update) FromChat() *Chat {
	switch {
	case u.Message != nil:
		return u.Message.Chat
	case u.EditedMessage != nil:
		return u.EditedMessage.Chat
	case u.ChannelPost != nil:
		return u.ChannelPost.Chat
	case u.EditedChannelPost != nil:
		return u.EditedChannelPost.Chat
	case u.CallbackQuery != nil:
		return u.CallbackQuery.Message.Chat
	default:
		return nil
	}
}

// UpdatesChannel is the channel for getting updates.
type UpdatesChannel <-chan Update

// Clear discards all unprocessed incoming updates.
func (ch UpdatesChannel) Clear() {
	for len(ch) != 0 {
		<-ch
	}
}

// Describes the current status of a webhook.
type WebhookInfo struct {
	URL                          string   `json:"url"`                                       // Webhook URL
	HasCustomCertificate         bool     `json:"has_custom_certificate"`                    // True, if a custom certificate was provided for webhook certificate checks
	PendingUpdateCount           int      `json:"pending_update_count"`                      // Number of updates awaiting delivery
	IPAddress                    string   `json:"ip_address,omitempty"`                      // Optional. Currently used webhook IP address
	LastErrorDate                int      `json:"last_error_date,omitempty"`                 // Optional. Unix time for the most recent error when delivering an update via webhook
	LastErrorMessage             string   `json:"last_error_message,omitempty"`              // Optional. Error message for the most recent error when delivering an update via webhook
	LastSynchronizationErrorDate int      `json:"last_synchronization_error_date,omitempty"` // Optional. Unix time of the most recent error when synchronizing available updates with Telegram datacenters
	MaxConnections               int      `json:"max_connections,omitempty"`                 // Optional. The maximum allowed number of simultaneous HTTPS connections for update delivery
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`                 // Optional. A list of update types the bot is subscribed to
}

// IsSet returns true if a webhook is currently set.
func (info WebhookInfo) IsSet() bool {
	return info.URL != ""
}

//
//
//
// Basic types
//
//
//

// This object represents a Telegram user or bot.
type User struct {
	ID                      int    `json:"id"`                                    // Unique identifier for this user or bot
	IsBot                   bool   `json:"is_bot"`                                // True, if this user is a bot
	FirstName               string `json:"first_name"`                            // User's or bot's first name
	LastName                string `json:"last_name,omitempty"`                   // Optional. User's or bot's last name
	UserName                string `json:"username,omitempty"`                    // Optional. User's or bot's username
	LanguageCode            string `json:"language_code,omitempty"`               // Optional. IETF language tag of the user's language
	IsPremium               bool   `json:"is_premium,omitempty"`                  // Optional. True, if this user is a Telegram Premium user
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`    // Optional. True, if this user added the bot to the attachment menu
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`             // Optional. True, if the bot can be invited to groups (Returned only in getMe)
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"` // Optional. True, if privacy mode is disabled for the bot (Returned only in getMe)
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`     // Optional. True, if the bot supports inline queries (Returned only in getMe)
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
func (u *User) String() string {
	if u == nil {
		return ""
	}
	if u.UserName != "" {
		return u.UserName
	}

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	return name
}

// This object represents a chat.
type Chat struct {
	ID                                 int              `json:"id"`                                                // Unique identifier for this chat
	Type                               string           `json:"type"`                                              // Type of chat, can be either "private", "group", "supergroup", or "channel"
	Title                              string           `json:"title,omitempty"`                                   // Optional. Title, for supergroups, channels, and group chats
	Username                           string           `json:"username,omitempty"`                                // Optional. Username, for private chats, supergroups, and channels if available
	FirstName                          string           `json:"first_name,omitempty"`                              // Optional. First name of the other party in a private chat
	LastName                           string           `json:"last_name,omitempty"`                               // Optional. Last name of the other party in a private chat
	IsForum                            bool             `json:"is_forum,omitempty"`                                // Optional. True, if the supergroup chat is a forum (has topics enabled)
	Photo                              *ChatPhoto       `json:"photo,omitempty"`                                   // Optional. Chat photo. Returned only in getChat.
	ActiveUsernames                    []string         `json:"active_usernames,omitempty"`                        // Optional. List of all active chat usernames; for private chats, supergroups, and channels. Returned only in getChat.
	EmojiStatusCustomEmojiID           string           `json:"emoji_status_custom_emoji_id,omitempty"`            // Optional. Custom emoji identifier of emoji status of the other party in a private chat. Returned only in getChat.
	Bio                                string           `json:"bio,omitempty"`                                     // Optional. Bio of the other party in a private chat. Returned only in getChat.
	HasPrivateForwards                 bool             `json:"has_private_forwards,omitempty"`                    // Optional. True, if privacy settings of the other party in the private chat allow using tg:// user?id=<user_id> links only in chats with the user. Returned only in getChat.
	HasRestrictedVoiceAndVideoMessages bool             `json:"has_restricted_voice_and_video_messages,omitempty"` // Optional. True, if the privacy settings of the other party restrict sending voice and video note messages in the private chat. Returned only in getChat.
	JoinToSendMessages                 bool             `json:"join_to_send_messages,omitempty"`                   // Optional. True, if users need to join the supergroup before they can send messages. Returned only in getChat.
	JoinByRequest                      bool             `json:"join_by_request,omitempty"`                         // Optional. True, if all users directly joining the supergroup need to be approved by supergroup administrators. Returned only in getChat.
	Description                        string           `json:"description,omitempty"`                             // Optional. Description, for groups, supergroups, and channel chats. Returned only in getChat.
	InviteLink                         string           `json:"invite_link,omitempty"`                             // Optional. Primary invite link, for groups, supergroups, and channel chats. Returned only in getChat.
	PinnedMessage                      *Message         `json:"pinned_message,omitempty"`                          // Optional. The most recent pinned message (by sending date). Returned only in getChat.
	Permissions                        *ChatPermissions `json:"permissions,omitempty"`                             // Optional. Default chat member permissions, for groups and supergroups. Returned only in getChat.
	SlowModeDelay                      int              `json:"slow_mode_delay,omitempty"`                         // Optional. For supergroups, the minimum allowed delay between consecutive messages sent by each unprivileged user; in seconds. Returned only in getChat.
	MessageAutoDeleteTime              int              `json:"message_auto_delete_time,omitempty"`                // Optional. The time after which all messages sent to the chat will be automatically deleted; in seconds. Returned only in getChat.
	HasAggressiveAntiSpamEnabled       bool             `json:"has_aggressive_anti_spam_enabled,omitempty"`        // Optional. True, if aggressive anti-spam checks are enabled in the supergroup. The field is only available to chat administrators. Returned only in getChat.
	HasHiddenMembers                   bool             `json:"has_hidden_members,omitempty"`                      // Optional. True, if non-administrators can only get the list of bots and administrators in the chat. Returned only in getChat.
	HasProtectedContent                bool             `json:"has_protected_content,omitempty"`                   // Optional. True, if messages from the chat can't be forwarded to other chats. Returned only in getChat.
	StickerSetName                     string           `json:"sticker_set_name,omitempty"`                        // Optional. For supergroups, name of the group sticker set. Returned only in getChat.
	CanSetStickerSet                   bool             `json:"can_set_sticker_set,omitempty"`                     // Optional. True, if the bot can change the group sticker set. Returned only in getChat.
	LinkedChatID                       int              `json:"linked_chat_id,omitempty"`                          // Optional. Unique identifier for the linked chat, i.e., the discussion group identifier for a channel and vice versa; for supergroups and channel chats. Returned only in getChat.
	Location                           *ChatLocation    `json:"location,omitempty"`                                // Optional. For supergroups, the location to which the supergroup is connected. Returned only in getChat.
}

// IsPrivate returns if the Chat is a private conversation.
func (c Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns if the Chat is a group.
func (c Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsSuperGroup returns if the Chat is a supergroup.
func (c Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

// IsChannel returns if the Chat is a channel.
func (c Chat) IsChannel() bool {
	return c.Type == "channel"
}

// This object represents a message.
type Message struct {
	MessageID                     int                            `json:"message_id"`                                  // Unique message identifier inside this chat
	MessageThreadID               int                            `json:"message_thread_id,omitempty"`                 // Optional. Unique identifier of a message thread to which the message belongs; for supergroups only
	From                          *User                          `json:"from,omitempty"`                              // Optional. Sender of the message; empty for messages sent to channels
	SenderChat                    *Chat                          `json:"sender_chat,omitempty"`                       // Optional. Sender of the message, sent on behalf of a chat
	Date                          int                            `json:"date"`                                        // Date the message was sent in Unix time
	Chat                          *Chat                          `json:"chat"`                                        // Conversation the message belongs to
	ForwardFrom                   *User                          `json:"forward_from,omitempty"`                      // Optional. For forwarded messages, sender of the original message
	ForwardFromChat               *Chat                          `json:"forward_from_chat,omitempty"`                 // Optional. For messages forwarded from channels or from anonymous administrators, information about the original sender chat
	ForwardFromMessageID          int                            `json:"forward_from_message_id,omitempty"`           // Optional. For messages forwarded from channels, identifier of the original message in the channel
	ForwardSignature              string                         `json:"forward_signature,omitempty"`                 // Optional. For forwarded messages that were originally sent in channels or by an anonymous chat administrator, signature of the message sender if present
	ForwardSenderName             string                         `json:"forward_sender_name,omitempty"`               // Optional. Sender's name for messages forwarded from users who disallow adding a link to their account in forwarded messages
	ForwardDate                   int                            `json:"forward_date,omitempty"`                      // Optional. For forwarded messages, date the original message was sent in Unix time
	IsTopicMessage                bool                           `json:"is_topic_message,omitempty"`                  // Optional. True, if the message is sent to a forum topic
	IsAutomaticForward            bool                           `json:"is_automatic_forward,omitempty"`              // Optional. True, if the message is a channel post that was automatically forwarded to the connected discussion group
	ReplyToMessage                *Message                       `json:"reply_to_message,omitempty"`                  // Optional. For replies, the original message
	ViaBot                        *User                          `json:"via_bot,omitempty"`                           // Optional. Bot through which the message was sent
	EditDate                      int                            `json:"edit_date,omitempty"`                         // Optional. Date the message was last edited in Unix time
	HasProtectedContent           bool                           `json:"has_protected_content,omitempty"`             // Optional. True, if the message can't be forwarded
	MediaGroupID                  string                         `json:"media_group_id,omitempty"`                    // Optional. The unique identifier of a media message group this message belongs to
	AuthorSignature               string                         `json:"author_signature,omitempty"`                  // Optional. Signature of the post author for messages in channels, or the custom title of an anonymous group administrator
	Text                          string                         `json:"text,omitempty"`                              // Optional. For text messages, the actual UTF-8 text of the message
	Entities                      []*MessageEntity               `json:"entities,omitempty"`                          // Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
	Animation                     *Animation                     `json:"animation,omitempty"`                         // Optional. Message is an animation, information about the animation
	Audio                         *Audio                         `json:"audio,omitempty"`                             // Optional. Message is an audio file, information about the file
	Document                      *Document                      `json:"document,omitempty"`                          // Optional. Message is a general file, information about the file
	Photo                         []*PhotoSize                   `json:"photo,omitempty"`                             // Optional. Message is a photo, available sizes of the photo
	Sticker                       *Sticker                       `json:"sticker,omitempty"`                           // Optional. Message is a sticker, information about the sticker
	Video                         *Video                         `json:"video,omitempty"`                             // Optional. Message is a video, information about the video
	VideoNote                     *VideoNote                     `json:"video_note,omitempty"`                        // Optional. Message is a video note, information about the video message
	Voice                         *Voice                         `json:"voice,omitempty"`                             // Optional. Message is a voice message, information about the file
	Caption                       string                         `json:"caption,omitempty"`                           // Optional. Caption for the animation, audio, document, photo, video, or voice
	CaptionEntities               []*MessageEntity               `json:"caption_entities,omitempty"`                  // Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc. that appear in the caption
	HasMediaSpoiler               bool                           `json:"has_media_spoiler,omitempty"`                 // Optional. True, if the message media is covered by a spoiler animation
	Contact                       *Contact                       `json:"contact,omitempty"`                           // Optional. Message is a shared contact, information about the contact
	Dice                          *Dice                          `json:"dice,omitempty"`                              // Optional. Message is a dice with a random value
	Game                          *Game                          `json:"game,omitempty"`                              // Optional. Message is a game, information about the game
	Poll                          *Poll                          `json:"poll,omitempty"`                              // Optional. Message is a native poll, information about the poll
	Venue                         *Venue                         `json:"venue,omitempty"`                             // Optional. Message is a venue, information about the venue
	Location                      *Location                      `json:"location,omitempty"`                          // Optional. Message is a shared location, information about the location
	NewChatMembers                []*User                        `json:"new_chat_members,omitempty"`                  // Optional. New members that were added to the group or supergroup and information about them
	LeftChatMember                *User                          `json:"left_chat_member,omitempty"`                  // Optional. A member was removed from the group, information about them
	NewChatTitle                  string                         `json:"new_chat_title,omitempty"`                    // Optional. A chat title was changed to this value
	NewChatPhoto                  []*PhotoSize                   `json:"new_chat_photo,omitempty"`                    // Optional. A chat photo was changed to this value
	DeleteChatPhoto               bool                           `json:"delete_chat_photo,omitempty"`                 // Optional. Service message: the chat photo was deleted
	GroupChatCreated              bool                           `json:"group_chat_created,omitempty"`                // Optional. Service message: the group has been created
	SupergroupChatCreated         bool                           `json:"supergroup_chat_created,omitempty"`           // Optional. Service message: the supergroup has been created
	ChannelChatCreated            bool                           `json:"channel_chat_created,omitempty"`              // Optional. Service message: the channel has been created
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"` // Optional. Service message: auto-delete timer settings changed in the chat
	MigrateToChatID               int                            `json:"migrate_to_chat_id,omitempty"`                // Optional. The group has been migrated to a supergroup with the specified identifier
	MigrateFromChatID             int                            `json:"migrate_from_chat_id,omitempty"`              // Optional. The supergroup has been migrated from a group with the specified identifier
	PinnedMessage                 *Message                       `json:"pinned_message,omitempty"`                    // Optional. Specified message was pinned
	Invoice                       *Invoice                       `json:"invoice,omitempty"`                           // Optional. Message is an invoice for a payment, information about the invoice
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment,omitempty"`                // Optional. Message is a service message about a successful payment, information about the payment
	UserShared                    *UserShared                    `json:"user_shared,omitempty"`                       // Optional. Service message: a user was shared with the bot
	ChatShared                    *ChatShared                    `json:"chat_shared,omitempty"`                       // Optional. Service message: a chat was shared with the bot
	ConnectedWebsite              string                         `json:"connected_website,omitempty"`                 // Optional. The domain name of the website on which the user has logged in
	WriteAccessAllowed            *WriteAccessAllowed            `json:"write_access_allowed,omitempty"`              // Optional. Service message: the user allowed the bot added to the attachment menu to write messages
	PassportData                  *PassportData                  `json:"passport_data,omitempty"`                     // Optional. Telegram Passport data
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`         // Optional. Service message: A user in the chat triggered another user's proximity alert while sharing Live Location
	ForumTopicCreated             *ForumTopicCreated             `json:"forum_topic_created,omitempty"`               // Optional. Service message: forum topic created
	ForumTopicEdited              *ForumTopicEdited              `json:"forum_topic_edited,omitempty"`                // Optional. Service message: forum topic edited
	ForumTopicClosed              *ForumTopicClosed              `json:"forum_topic_closed,omitempty"`                // Optional. Service message: forum topic closed
	ForumTopicReopened            *ForumTopicReopened            `json:"forum_topic_reopened,omitempty"`              // Optional. Service message: forum topic reopened
	GeneralForumTopicHidden       *GeneralForumTopicHidden       `json:"general_forum_topic_hidden,omitempty"`        // Optional. Service message: the 'General' forum topic hidden
	GeneralForumTopicUnhidden     *GeneralForumTopicUnhidden     `json:"general_forum_topic_unhidden,omitempty"`      // Optional. Service message: the 'General' forum topic unhidden
	VideoChatScheduled            *VideoChatScheduled            `json:"video_chat_scheduled,omitempty"`              // Optional. Service message: video chat scheduled
	VideoChatStarted              *VideoChatStarted              `json:"video_chat_started,omitempty"`                // Optional. Service message: video chat started
	VideoChatEnded                *VideoChatEnded                `json:"video_chat_ended,omitempty"`                  // Optional. Service message: video chat ended
	VideoChatParticipantsInvited  *VideoChatParticipantsInvited  `json:"video_chat_participants_invited,omitempty"`   // Optional. Service message: new participants invited to a video chat
	WebAppData                    *WebAppData                    `json:"web_app_data,omitempty"`                      // Optional. Service message: data sent by a Web App
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup,omitempty"`                      // Optional. Inline keyboard attached to the message. login_url buttons are represented as ordinary url buttons.
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsCommand returns true if message starts with a "bot_command" entity.
func (m *Message) IsCommand() bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
	return entity.Offset == 0 && entity.IsCommand()
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is removed. Use
// CommandWithAt() if you do not want that.
func (m *Message) Command() string {
	command := m.CommandWithAt()

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAt checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is not removed. Use Command()
// if you want that.
func (m *Message) CommandWithAt() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := m.Entities[0]
	return m.Text[1:entity.Length]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
//
// Note: The first character after the command name is omitted:
// - "/foo bar baz" yields "bar baz", not " bar baz"
// - "/foo-bar baz" yields "bar baz", too
// Even though the latter is not a command conforming to the spec, the API
// marks "/foo" as command entity.
func (m *Message) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := m.Entities[0]

	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}

// This object represents a unique message identifier.
type MessageId struct {
	MessageID int `json:"message_id"` // Unique message identifier
}

// This object represents one special entity in a text message. For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	Type          string `json:"type"`                      // Type of the entity. Currently, can be “mention” (@username), “hashtag” (#hashtag), “cashtag” ($USD), “bot_command” (/start@jobs_bot), “url” (https:// telegram.org), “email” (do-not-reply@telegram.org), “phone_number” (+1-212-555-0123), “bold” (bold text), “italic” (italic text), “underline” (underlined text), “strikethrough” (strikethrough text), “spoiler” (spoiler message), “code” (monowidth string), “pre” (monowidth block), “text_link” (for clickable text URLs), “text_mention” (for users without usernames), “custom_emoji” (for inline custom emoji stickers)
	Offset        int    `json:"offset"`                    // Offset in UTF-16 code units to the start of the entity
	Length        int    `json:"length"`                    // Length of the entity in UTF-16 code units
	URL           string `json:"url,omitempty"`             // Optional. URL that will be opened after the user taps on the text (for "text_link" entities)
	User          *User  `json:"user,omitempty"`            // Optional. Mentioned user (for "text_mention" entities)
	Language      string `json:"language,omitempty"`        // Optional. Programming language of the entity text (for "pre" entities)
	CustomEmojiID string `json:"custom_emoji_id,omitempty"` // Optional. Unique identifier of the custom emoji (for "custom_emoji" entities)
}

// ParseURL attempts to parse a URL contained within a MessageEntity.
func (e MessageEntity) ParseURL() (*url.URL, error) {
	if e.URL == "" {
		return nil, fmt.Errorf("bad or empty url")
	}

	return url.Parse(e.URL)
}

// IsMention returns true if the type of the message entity is "mention" (@username).
func (e MessageEntity) IsMention() bool {
	return e.Type == "mention"
}

// IsTextMention returns true if the type of the message entity is "text_mention"
// (At this time, the user field exists, and occurs when tagging a member without a username)
func (e MessageEntity) IsTextMention() bool {
	return e.Type == "text_mention"
}

// IsHashtag returns true if the type of the message entity is "hashtag".
func (e MessageEntity) IsHashtag() bool {
	return e.Type == "hashtag"
}

// IsCommand returns true if the type of the message entity is "bot_command".
func (e MessageEntity) IsCommand() bool {
	return e.Type == "bot_command"
}

// IsURL returns true if the type of the message entity is "url".
func (e MessageEntity) IsURL() bool {
	return e.Type == "url"
}

// IsEmail returns true if the type of the message entity is "email".
func (e MessageEntity) IsEmail() bool {
	return e.Type == "email"
}

// IsBold returns true if the type of the message entity is "bold" (bold text).
func (e MessageEntity) IsBold() bool {
	return e.Type == "bold"
}

// IsItalic returns true if the type of the message entity is "italic" (italic text).
func (e MessageEntity) IsItalic() bool {
	return e.Type == "italic"
}

// IsCode returns true if the type of the message entity is "code" (monowidth string).
func (e MessageEntity) IsCode() bool {
	return e.Type == "code"
}

// IsPre returns true if the type of the message entity is "pre" (monowidth block).
func (e MessageEntity) IsPre() bool {
	return e.Type == "pre"
}

// IsTextLink returns true if the type of the message entity is "text_link" (clickable text URL).
func (e MessageEntity) IsTextLink() bool {
	return e.Type == "text_link"
}

// This object represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	FileID       string `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Width        int    `json:"width"`               // Photo width
	Height       int    `json:"height"`              // Photo height
	FileSize     int    `json:"file_size,omitempty"` // Optional. File size in bytes
}

// This object represents an animation file (GIF or H.264/MPEG-4 AVC video without sound).
type Animation struct {
	FileID       string     `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string     `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Width        int        `json:"width"`               // Video width as defined by the sender
	Height       int        `json:"height"`              // Video height as defined by the sender
	Duration     int        `json:"duration"`            // Duration of the video in seconds as defined by the sender
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"` // Optional. Animation thumbnail as defined by the sender
	FileName     string     `json:"file_name,omitempty"` // Optional. Original animation filename as defined by the sender
	MimeType     string     `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by the sender
	FileSize     int        `json:"file_size,omitempty"` // Optional. File size in bytes. It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
}

// This object represents an audio file to be treated as music by the Telegram clients.
type Audio struct {
	FileID       string     `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string     `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Duration     int        `json:"duration"`            // Duration of the audio in seconds as defined by the sender
	Performer    string     `json:"performer,omitempty"` // Optional. Performer of the audio as defined by the sender or by audio tags
	Title        string     `json:"title,omitempty"`     // Optional. Title of the audio as defined by the sender or by audio tags
	FileName     string     `json:"file_name,omitempty"` // Optional. Original filename as defined by the sender
	MimeType     string     `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by the sender
	FileSize     int        `json:"file_size,omitempty"` // Optional. File size in bytes. It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"` // Optional. Thumbnail of the album cover to which the music file belongs
}

// This object represents a general file (as opposed to photos, voice messages and audio files).
type Document struct {
	FileID       string     `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string     `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"` // Optional. Document thumbnail as defined by the sender
	FileName     string     `json:"file_name,omitempty"` // Optional. Original filename as defined by the sender
	MimeType     string     `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by the sender
	FileSize     int        `json:"file_size,omitempty"` // Optional. File size in bytes. It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
}

// This object represents a video file.
type Video struct {
	FileID       string     `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string     `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Width        int        `json:"width"`               // Video width as defined by the sender
	Height       int        `json:"height"`              // Video height as defined by the sender
	Duration     int        `json:"duration"`            // Duration of the video in seconds as defined by the sender
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"` // Optional. Video thumbnail
	FileName     string     `json:"file_name,omitempty"` // Optional. Original filename as defined by the sender
	MimeType     string     `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by the sender
	FileSize     int        `json:"file_size,omitempty"` // Optional. File size in bytes. It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
}

// This object represents a video message (available in Telegram apps as of v.4.0).
type VideoNote struct {
	FileID       string     `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string     `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Length       int        `json:"length"`              // Video width and height (diameter of the video message) as defined by the sender
	Duration     int        `json:"duration"`            // Duration of the video in seconds as defined by the sender
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"` // Optional. Video thumbnail
	FileSize     int        `json:"file_size,omitempty"` // Optional. File size in bytes
}

// This object represents a voice note.
type Voice struct {
	FileID       string `json:"file_id"`             // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string `json:"file_unique_id"`      // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	Duration     int    `json:"duration"`            // Duration of the audio in seconds as defined by the sender
	MimeType     string `json:"mime_type,omitempty"` // Optional. MIME type of the file as defined by the sender
	FileSize     int    `json:"file_size,omitempty"` // Optional. File size in bytes. It can be bigger than 2^31 and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this value.
}

// This object represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`        // Contact's phone number
	FirstName   string `json:"first_name"`          // Contact's first name
	LastName    string `json:"last_name,omitempty"` // Optional. Contact's last name
	UserID      int    `json:"user_id,omitempty"`   // Optional. Contact's user identifier in Telegram. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.
	VCard       string `json:"vcard,omitempty"`     // Optional. Additional data about the contact in the form of a vCard
}

// This object represents an animated emoji that displays a random value.
type Dice struct {
	Emoji string `json:"emoji"` // Emoji on which the dice throw animation is based
	Value int    `json:"value"` // Value of the dice, 1-6 for "🎲", "🎯" and "🎳" base emoji, 1-5 for "🏀" and "⚽" base emoji, 1-64 for "🎰" base emoji
}

// This object contains information about one answer option in a poll.
type PollOption struct {
	Text       string `json:"text"`        // Option text, 1-100 characters
	VoterCount int    `json:"voter_count"` // Number of users that voted for this option
}

// This object represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID    string `json:"poll_id"`    // Unique poll identifier
	User      User   `json:"user"`       // The user who changed the answer to the poll
	OptionIDs []int  `json:"option_ids"` // 0-based identifiers of answer options chosen by the user. May be empty if the user retracted their vote.
}

// This object contains information about a poll.
type Poll struct {
	ID                    string          `json:"id"`                             // Unique poll identifier
	Question              string          `json:"question"`                       // Poll question, 1-300 characters
	Options               []PollOption    `json:"options"`                        // List of poll options
	TotalVoterCount       int             `json:"total_voter_count"`              // Total number of users that voted in the poll
	IsClosed              bool            `json:"is_closed"`                      // True if the poll is closed
	IsAnonymous           bool            `json:"is_anonymous"`                   // True if the poll is anonymous
	Type                  string          `json:"type"`                           // Poll type, currently can be "regular" or "quiz"
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`        // True if the poll allows multiple answers
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`    // Optional. 0-based identifier of the correct answer option. Available only for quizzes or closed polls sent by the bot or to the private chat with the bot.
	Explanation           string          `json:"explanation,omitempty"`          // Optional. Text that is shown when a user chooses an incorrect answer or taps on the lamp icon in a quiz-style poll, 0-200 characters
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"` // Optional. Special entities like usernames, URLs, bot commands, etc. that appear in the explanation
	OpenPeriod            int             `json:"open_period,omitempty"`          // Optional. Amount of time in seconds the poll will be active after creation
	CloseDate             int             `json:"close_date,omitempty"`           // Optional. Point in time (Unix timestamp) when the poll will be automatically closed
}

// This object represents a point on the map.
type Location struct {
	Longitude            float64 `json:"longitude"`                        // Longitude as defined by sender
	Latitude             float64 `json:"latitude"`                         // Latitude as defined by sender
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`    // Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	LivePeriod           int     `json:"live_period,omitempty"`            // Optional. Time relative to the message sending date, during which the location can be updated; in seconds. For active live locations only.
	Heading              int     `json:"heading,omitempty"`                // Optional. The direction in which the user is moving, in degrees; 1-360. For active live locations only.
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"` // Optional. The maximum distance for proximity alerts about approaching another chat member, in meters. For sent live locations only.
}

// This object represents a venue.
type Venue struct {
	Location        Location `json:"location"`                    // Venue location. Can't be a live location.
	Title           string   `json:"title"`                       // Name of the venue.
	Address         string   `json:"address"`                     // Address of the venue.
	FoursquareID    string   `json:"foursquare_id,omitempty"`     // Optional. Foursquare identifier of the venue.
	FoursquareType  string   `json:"foursquare_type,omitempty"`   // Optional. Foursquare type of the venue. (For example, "arts_entertainment/default", "arts_entertainment/aquarium", or "food/icecream".)
	GooglePlaceID   string   `json:"google_place_id,omitempty"`   // Optional. Google Places identifier of the venue.
	GooglePlaceType string   `json:"google_place_type,omitempty"` // Optional. Google Places type of the venue.
}

// Describes data sent from a Web App to the bot.
type WebAppData struct {
	Data       string `json:"data"`        // The data. Be aware that a bad client can send arbitrary data in this field.
	ButtonText string `json:"button_text"` // Text of the web_app keyboard button from which the Web App was opened. Be aware that a bad client can send arbitrary data in this field.
}

// This object represents the content of a service message, sent whenever a user in the chat triggers a proximity alert set by another user.
type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"` // User that triggered the alert.
	Watcher  User `json:"watcher"`  // User that set the alert.
	Distance int  `json:"distance"` // The distance between the users.
}

// This object represents a service message about a change in auto-delete timer settings.
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"` // New auto-delete time for messages in the chat; in seconds.
}

// This object represents a service message about a new forum topic created in the chat.
type ForumTopicCreated struct {
	Name              string `json:"name"`                           // Name of the topic.
	IconColor         int    `json:"icon_color"`                     // Color of the topic icon in RGB format.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"` // Optional. Unique identifier of the custom emoji shown as the topic icon.
}

// This object represents a service message about a forum topic closed in the chat. Currently holds no information.
type ForumTopicClosed struct {
}

// This object represents a service message about an edited forum topic.
type ForumTopicEdited struct {
	Name              string `json:"name,omitempty"`                 // Optional. New name of the topic, if it was edited.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"` // Optional. New identifier of the custom emoji shown as the topic icon, if it was edited; an empty string if the icon was removed.
}

// This object represents a service message about a forum topic reopened in the chat. Currently holds no information.
type ForumTopicReopened struct {
}

// This object represents a service message about General forum topic hidden in the chat. Currently holds no information.
type GeneralForumTopicHidden struct {
}

// This object represents a service message about General forum topic unhidden in the chat. Currently holds no information.
type GeneralForumTopicUnhidden struct {
}

// This object contains information about the user whose identifier was shared with the bot using a KeyboardButtonRequestUser button.
type UserShared struct {
	RequestID int `json:"request_id"` // Identifier of the request
	UserID    int `json:"user_id"`    // Identifier of the shared user. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier. The bot may not have access to the user and could be unable to use this identifier, unless the user is already known to the bot by some other means.
}

// This object contains information about the chat whose identifier was shared with the bot using a KeyboardButtonRequestChat button.
type ChatShared struct {
	RequestID int `json:"request_id"` // Identifier of the request
	ChatID    int `json:"chat_id"`    // Identifier of the shared chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier. The bot may not have access to the chat and could be unable to use this identifier, unless the chat is already known to the bot by some other means.
}

// This object represents a service message about a user allowing a bot to write messages after adding the bot to the attachment menu or launching a Web App from a link.
type WriteAccessAllowed struct {
	WebAppName string `json:"web_app_name,omitempty"` // Optional. Name of the Web App which was launched from a link.
}

// This object represents a service message about a video chat scheduled in the chat.
type VideoChatScheduled struct {
	StartDate int `json:"start_date"` // Point in time (Unix timestamp) when the video chat is supposed to be started by a chat administrator.
}

// Time converts the scheduled start date into a Time.
func (m *VideoChatScheduled) Time() time.Time {
	return time.Unix(int64(m.StartDate), 0)
}

// This object represents a service message about a video chat started in the chat. Currently holds no information.
type VideoChatStarted struct {
}

// This object represents a service message about a video chat ended in the chat.
type VideoChatEnded struct {
	Duration int `json:"duration"` // Video chat duration in seconds.
}

// This object represents a service message about new members invited to a video chat.
type VideoChatParticipantsInvited struct {
	Users []User `json:"users"` // New members that were invited to the video chat.
}

// This object represent a user's profile pictures.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"` // Total number of profile pictures the target user has.
	Photos     [][]PhotoSize `json:"photos"`      // Requested profile pictures (in up to 4 sizes each).
}

// This object represents a file ready to be downloaded.
// The file can be downloaded via the link https:// api.telegram.org/file/bot<token>/<file_path>.
// It is guaranteed that the link will be valid for at least 1 hour. When the link expires, a new one can be requested by calling getFile.
// The maximum file size to download is 20 MB
type File struct {
	FileID       string `json:"file_id"`        // Identifier for this file, which can be used to download or reuse the file.
	FileUniqueID string `json:"file_unique_id"` // Unique identifier for this file, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.
	FileSize     int    `json:"file_size"`      // Optional. File size in bytes.
	FilePath     string `json:"file_path"`      // Optional. File path. Use https:// api.telegram.org/file/bot<token>/<file_path> to get the file.
}

// Link returns a full path to the download URL for a File.
//
// It requires the Bot token to create the link.
func (f *File) Link(client Client) string {
	return client.fileEndpoint + "/" + f.FilePath
}

// Describes a Web App.
type WebAppInfo struct {
	URL string `json:"url"` // An HTTPS URL of a Web App to be opened with additional data as specified in Initializing Web Apps.
}

// This object represents a custom keyboard with reply options (see Introduction to bots for details and examples).
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`                // Array of button rows, each represented by an Array of KeyboardButton objects.
	IsPersistent          bool               `json:"is_persistent"`           // Optional. Requests clients to always show the keyboard when the regular keyboard is hidden. Defaults to false.
	ResizeKeyboard        bool               `json:"resize_keyboard"`         // Optional. Requests clients to resize the keyboard vertically for optimal fit. Defaults to false.
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`       // Optional. Requests clients to hide the keyboard as soon as it's been used. Defaults to false.
	InputFieldPlaceholder string             `json:"input_field_placeholder"` // Optional. The placeholder to be shown in the input field when the keyboard is active; 1-64 characters.
	Selective             bool               `json:"selective"`               // Optional. Use this parameter if you want to show the keyboard to specific users only. Defaults to false.
}

// This object represents one button of the reply keyboard.
// For simple text buttons, String can be used instead of this object to specify the button text.
// The optional fields web_app, request_user, request_chat, request_contact, request_location, and request_poll are mutually exclusive.
type KeyboardButton struct {
	Text            string                     `json:"text"`                       // Text of the button. If none of the optional fields are used, it will be sent as a message when the button is pressed.
	RequestUser     *KeyboardButtonRequestUser `json:"request_user,omitempty"`     // Optional. If specified, pressing the button will open a list of suitable users. Tapping on any user will send their identifier to the bot in a “user_shared” service message. Available in private chats only.
	RequestChat     *KeyboardButtonRequestChat `json:"request_chat,omitempty"`     // Optional. If specified, pressing the button will open a list of suitable chats. Tapping on a chat will send its identifier to the bot in a “chat_shared” service message. Available in private chats only.
	RequestContact  bool                       `json:"request_contact,omitempty"`  // Optional. If true, the user's phone number will be sent as a contact when the button is pressed. Available in private chats only.
	RequestLocation bool                       `json:"request_location,omitempty"` // Optional. If true, the user's current location will be sent when the button is pressed. Available in private chats only.
	RequestPoll     *KeyboardButtonPollType    `json:"request_poll,omitempty"`     // Optional. If specified, the user will be asked to create a poll and send it to the bot when the button is pressed. Available in private chats only.
	WebApp          *WebAppInfo                `json:"web_app,omitempty"`          // Optional. If specified, the described Web App will be launched when the button is pressed. The Web App will be able to send a “web_app_data” service message. Available in private chats only.
}

// This object defines the criteria used to request a suitable user.
// The identifier of the selected user will be shared with the bot when the corresponding button is pressed.
type KeyboardButtonRequestUser struct {
	RequestID     int  `json:"request_id"`                // Signed 32-bit identifier of the request, which will be received back in the UserShared object. Must be unique within the message.
	UserIsBot     bool `json:"user_is_bot,omitempty"`     // Optional. Pass true to request a bot, pass false to request a regular user. If not specified, no additional restrictions are applied.
	UserIsPremium bool `json:"user_is_premium,omitempty"` // Optional. Pass true to request a premium user, pass false to request a non-premium user. If not specified, no additional restrictions are applied.
}

// This object defines the criteria used to request a suitable chat.
// The identifier of the selected chat will be shared with the bot when the corresponding button is pressed.
type KeyboardButtonRequestChat struct {
	RequestID               int                      `json:"request_id"`                          // Signed 32-bit identifier of the request, which will be received back in the ChatShared object. Must be unique within the message.
	ChatIsChannel           bool                     `json:"chat_is_channel"`                     // Pass true to request a channel chat, pass false to request a group or a supergroup chat.
	ChatIsForum             bool                     `json:"chat_is_forum,omitempty"`             // Optional. Pass true to request a forum supergroup, pass false to request a non-forum chat. If not specified, no additional restrictions are applied.
	ChatHasUsername         bool                     `json:"chat_has_username,omitempty"`         // Optional. Pass true to request a supergroup or a channel with a username, pass false to request a chat without a username. If not specified, no additional restrictions are applied.
	ChatIsCreated           bool                     `json:"chat_is_created,omitempty"`           // Optional. Pass true to request a chat owned by the user. Otherwise, no additional restrictions are applied.
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"` // Optional. A JSON-serialized object listing the required administrator rights of the user in the chat. The rights must be a superset of bot_administrator_rights. If not specified, no additional restrictions are applied.
	BotAdministratorRights  *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`  // Optional. A JSON-serialized object listing the required administrator rights of the bot in the chat. The rights must be a subset of user_administrator_rights. If not specified, no additional restrictions are applied.
	BotIsMember             bool                     `json:"bot_is_member,omitempty"`             // Optional. Pass true to request a chat with the bot as a member. Otherwise, no additional restrictions are applied.
}

// This object represents type of a poll, which is allowed to be created and sent when the corresponding button is pressed.
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"` // Optional. If "quiz" is passed, the user will be allowed to create only polls in the quiz mode. If "regular" is passed, only regular polls will be allowed. Otherwise, the user will be allowed to create a poll of any type.
}

// Upon receiving a message with this object, Telegram clients will remove the current custom keyboard and
// display the default letter-keyboard. By default, custom keyboards are displayed until a new keyboard is sent by a bot.
// An exception is made for one-time keyboards that are hidden immediately after the user presses a button (see ReplyKeyboardMarkup).
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`     // Requests clients to remove the custom keyboard (user will not be able to summon this keyboard; if you want to hide the keyboard from sight but keep it accessible, use one_time_keyboard in ReplyKeyboardMarkup).
	Selective      bool `json:"selective,omitempty"` // Optional. Use this parameter if you want to remove the keyboard for specific users only.
}

// This object represents an inline keyboard that appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"` // Array of button rows
}

// This object represents one button of an inline keyboard. You must use exactly one of the optional fields.
type InlineKeyboardButton struct {
	Text                         string                       `json:"text"`                                       // Label text on the button
	URL                          *string                      `json:"url,omitempty"`                              // Optional. URL to be opened when the button is pressed
	CallbackData                 *string                      `json:"callback_data,omitempty"`                    // Optional. Data to be sent in a callback query to the bot when the button is pressed
	WebApp                       *WebAppInfo                  `json:"web_app,omitempty"`                          // Optional. Description of the Web App that will be launched when the button is pressed
	LoginURL                     *LoginURL                    `json:"login_url,omitempty"`                        // Optional. An HTTPS URL used to automatically authorize the user
	SwitchInlineQuery            *string                      `json:"switch_inline_query,omitempty"`              // Optional. Insert the bot's username and the specified inline query in the input field
	SwitchInlineQueryCurrentChat *string                      `json:"switch_inline_query_current_chat,omitempty"` // Optional. Insert the bot's username and the specified inline query in the current chat's input field
	SwitchInlineQueryChosenChat  *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`  // Optional. Prompt the user to select a chat and insert the bot's username and the specified inline query
	CallbackGame                 *CallbackGame                `json:"callback_game,omitempty"`                    // Optional. Description of the game that will be launched when the button is pressed
	Pay                          bool                         `json:"pay,omitempty"`                              // Optional. Specify True to send a Pay button
}

// This object represents a parameter of the inline keyboard button used to automatically authorize a user.
// Serves as a great replacement for the Telegram Login Widget when the user is coming from Telegram.
// All the user needs to do is tap/click a button and confirm that they want to log in.
// Telegram apps support these buttons as of version 5.7.
type LoginURL struct {
	URL                string `json:"url"`                            // HTTPS URL to be opened with user authorization data added to the query string when the button is pressed
	ForwardText        string `json:"forward_text,omitempty"`         // Optional. New text of the button in forwarded messages
	BotUsername        string `json:"bot_username,omitempty"`         // Optional. Username of a bot used for user authorization
	RequestWriteAccess bool   `json:"request_write_access,omitempty"` // Optional. Request permission for the bot to send messages to the user
}

// This object represents an inline button that switches the current user to inline mode in a chosen chat, with an optional default inline query.
type SwitchInlineQueryChosenChat struct {
	Query             string `json:"query,omitempty"`               // Optional. Default inline query to be inserted in the input field
	AllowUserChats    bool   `json:"allow_user_chats,omitempty"`    // Optional. True if private chats with users can be chosen
	AllowBotChats     bool   `json:"allow_bot_chats,omitempty"`     // Optional. True if private chats with bots can be chosen
	AllowGroupChats   bool   `json:"allow_group_chats,omitempty"`   // Optional. True if group and supergroup chats can be chosen
	AllowChannelChats bool   `json:"allow_channel_chats,omitempty"` // Optional. True if channel chats can be chosen
}

// This object represents an incoming callback query from a callback button in an inline keyboard.
// If the button that originated the query was attached to a message sent by the bot, the field message will be present.
// If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present.
// Exactly one of the fields data or game_short_name will be present.
type CallbackQuery struct {
	ID              string   `json:"id"`                          // Unique identifier for this query
	From            *User    `json:"from"`                        // Sender
	Message         *Message `json:"message,omitempty"`           // Optional. Message with the callback button that originated the query
	InlineMessageID string   `json:"inline_message_id,omitempty"` // Optional. Identifier of the message sent via the bot in inline mode, that originated the query
	ChatInstance    string   `json:"chat_instance"`               // Global identifier corresponding to the chat to which the message with the callback button was sent
	Data            string   `json:"data,omitempty"`              // Optional. Data associated with the callback button
	GameShortName   string   `json:"game_short_name,omitempty"`   // Optional. Short name of a Game to be returned
}

// Upon receiving a message with this object, Telegram clients will display a reply interface to the user
// (act as if the user has selected the bot's message and tapped 'Reply').
// This can be extremely useful if you want to create user-friendly step-by-step interfaces without having to sacrifice privacy mode.
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`                       // Shows reply interface to the user
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"` // Optional. Placeholder to be shown in the input field when the reply is active
	Selective             bool   `json:"selective,omitempty"`               // Optional. Force reply from specific users only
}

// This object represents a chat photo.
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`        // File identifier of small (160x160) chat photo
	SmallFileUniqueID string `json:"small_file_unique_id"` // Unique file identifier of small (160x160) chat photo
	BigFileID         string `json:"big_file_id"`          // File identifier of big (640x640) chat photo
	BigFileUniqueID   string `json:"big_file_unique_id"`   // Unique file identifier of big (640x640) chat photo
}

// Represents an invite link for a chat.
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`                          // The invite link
	Creator                 User   `json:"creator"`                              // Creator of the link
	CreatesJoinRequest      bool   `json:"creates_join_request"`                 // True, if users joining via the link need to be approved by chat administrators
	IsPrimary               bool   `json:"is_primary"`                           // True, if the link is primary
	IsRevoked               bool   `json:"is_revoked"`                           // True, if the link is revoked
	Name                    string `json:"name,omitempty"`                       // Optional. Invite link name
	ExpireDate              int    `json:"expire_date,omitempty"`                // Optional. Point in time (Unix timestamp) when the link will expire or has been expired
	MemberLimit             int    `json:"member_limit,omitempty"`               // Optional. The maximum number of users that can be members of the chat simultaneously after joining via this invite link
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"` // Optional. Number of pending join requests created using this link
}

// Represents the rights of an administrator in a chat.
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous"`                // True, if the user's presence in the chat is hidden
	CanManageChat       bool `json:"can_manage_chat"`             // True, if the administrator can access chat-related information and statistics
	CanDeleteMessages   bool `json:"can_delete_messages"`         // True, if the administrator can delete messages of other users
	CanManageVideoChats bool `json:"can_manage_video_chats"`      // True, if the administrator can manage video chats
	CanRestrictMembers  bool `json:"can_restrict_members"`        // True, if the administrator can restrict, ban or unban chat members
	CanPromoteMembers   bool `json:"can_promote_members"`         // True, if the administrator can add new administrators with a subset of their own privileges or demote administrators that they have promoted
	CanChangeInfo       bool `json:"can_change_info"`             // True, if the user is allowed to change the chat title, photo, and other settings
	CanInviteUsers      bool `json:"can_invite_users"`            // True, if the user is allowed to invite new users to the chat
	CanPostMessages     bool `json:"can_post_messages,omitempty"` // Optional. True, if the administrator can post messages in a channel
	CanEditMessages     bool `json:"can_edit_messages,omitempty"` // Optional. True, if the administrator can edit messages of other users and pin messages in a channel
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`  // Optional. True, if the user is allowed to pin messages in a group or supergroup
	CanManageTopics     bool `json:"can_manage_topics,omitempty"` // Optional. True, if the user is allowed to create, rename, close, and reopen forum topics in a supergroup
}

type ChatMember struct {
	Status                string `json:"status"`                              // The member's status in the chat, can be "creator", "administrator", "member", "restricted", "left", or "kicked"
	User                  User   `json:"user"`                                // Information about the user
	IsAnonymous           bool   `json:"is_anonymous,omitempty"`              // True, if the user's presence in the chat is hidden
	CustomTitle           string `json:"custom_title,omitempty"`              // Optional. Custom title for this user
	CanBeEdited           bool   `json:"can_be_edited,omitempty"`             // True, if the bot is allowed to edit administrator privileges of that user
	CanManageChat         bool   `json:"can_manage_chat,omitempty"`           // True, if the administrator can access the chat event log, chat statistics, message statistics in channels, see channel members, see anonymous administrators in supergroups, and ignore slow mode. Implied by any other administrator privilege
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty"`       // True, if the administrator can delete messages of other users
	CanManageVideoChats   bool   `json:"can_manage_video_chats,omitempty"`    // True, if the administrator can manage video chats
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty"`      // True, if the administrator can restrict, ban or unban chat members
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty"`       //True, if the administrator can add new administrators with a subset of their own privileges or demote administrators that they have promoted, directly or indirectly (promoted by administrators that were appointed by the user)
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`           // True, if the user is allowed to change the chat title, photo, and other settings
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`          // True, if the user is allowed to invite new users to the chat
	CanPostMessages       bool   `json:"can_post_messages,omitempty"`         // True, if the administrator can post in the channel; channels only
	CanEditMessages       bool   `json:"can_edit_messages,omitempty"`         // True, if the administrator can edit messages of other users and can pin messages; channels only
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`          // True, if the user is allowed to pin messages; groups and supergroups only
	CanManageTopics       bool   `json:"can_manage_topics,omitempty"`         // Optional. True, if the user is allowed to pin messages; groups and supergroups only
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`         // True, if the user is allowed to send text messages, contacts, invoices, locations, and venues
	CanSendAudios         bool   `json:"can_send_audios,omitempty"`           // True, if the user is allowed to send audios
	CanSendDocuments      bool   `json:"can_send_documents,omitempty"`        // True, if the user is allowed to send documents
	CanSendPhotos         bool   `json:"can_send_photos,omitempty"`           // True, if the user is allowed to send photos
	CanSendVideos         bool   `json:"can_send_videos,omitempty"`           // True, if the user is allowed to send videos
	CanSendVideoNotes     bool   `json:"can_send_video_notes,omitempty"`      // True, if the user is allowed to send video notes
	CanSendVoiceNotes     bool   `json:"can_send_voice_notes,omitempty"`      // True, if the user is allowed to send voice notes
	CanSendPolls          bool   `json:"can_send_polls,omitempty"`            // True, if the user is allowed to send polls
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty"`   // True, if the user is allowed to send animations, games, stickers, and use inline bots
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty"` // True, if the user is allowed to add web page previews to their messages
	UntilDate             int    `json:"until_date,omitempty"`                // Date when restrictions will be lifted for this user; unix time. If 0, then the user is restricted forever
}

// IsCreator returns if the ChatMember was the creator of the chat.
func (chat ChatMember) IsCreator() bool { return chat.Status == "creator" }

// IsAdministrator returns if the ChatMember is a chat administrator.
func (chat ChatMember) IsAdministrator() bool { return chat.Status == "administrator" }

// HasLeft returns if the ChatMember left the chat.
func (chat ChatMember) HasLeft() bool { return chat.Status == "left" }

// WasKicked returns if the ChatMember was kicked from the chat.
func (chat ChatMember) WasKicked() bool { return chat.Status == "kicked" }

// This object represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`                                  // Chat the user belongs to
	From                    User            `json:"from"`                                  // Performer of the action, which resulted in the change
	Date                    int             `json:"date"`                                  // Date the change was done in Unix time
	OldChatMember           ChatMember      `json:"old_chat_member"`                       // Previous information about the chat member
	NewChatMember           ChatMember      `json:"new_chat_member"`                       // New information about the chat member
	InviteLink              *ChatInviteLink `json:"invite_link,omitempty"`                 // Optional. Chat invite link, which was used by the user to join the chat (for joining by invite link events only)
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link,omitempty"` // Optional. True, if the user joined the chat via a chat folder invite link
}

// Represents a join request sent to a chat.
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`                  // Chat to which the request was sent
	From       User            `json:"from"`                  // User that sent the join request
	UserChatID int             `json:"user_chat_id"`          // Identifier of a private chat with the user who sent the join request
	Date       int             `json:"date"`                  // Date the request was sent in Unix time
	Bio        string          `json:"bio,omitempty"`         // Optional. Bio of the user
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"` // Optional. Chat invite link that was used by the user to send the join request
}

// Describes actions that a non-administrator user is allowed to take in a chat.
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`         // Optional. True, if the user is allowed to send text messages, contacts, invoices, locations, and venues
	CanSendAudios         bool `json:"can_send_audios,omitempty"`           // Optional. True, if the user is allowed to send audios
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`        // Optional. True, if the user is allowed to send documents
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`           // Optional. True, if the user is allowed to send photos
	CanSendVideos         bool `json:"can_send_videos,omitempty"`           // Optional. True, if the user is allowed to send videos
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`      // Optional. True, if the user is allowed to send video notes
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`      // Optional. True, if the user is allowed to send voice notes
	CanSendPolls          bool `json:"can_send_polls,omitempty"`            // Optional. True, if the user is allowed to send polls
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`   // Optional. True, if the user is allowed to send animations, games, stickers, and use inline bots
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"` // Optional. True, if the user is allowed to add web page previews to their messages
	CanChangeInfo         bool `json:"can_change_info,omitempty"`           // Optional. True, if the user is allowed to change the chat title, photo, and other settings. Ignored in public supergroups
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`          // Optional. True, if the user is allowed to invite new users to the chat
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`          // Optional. True, if the user is allowed to pin messages. Ignored in public supergroups
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`         // Optional. True, if the user is allowed to create forum topics. If omitted, defaults to the value of can_pin_messages
}

// Represents a location to which a chat is connected.
type ChatLocation struct {
	Location Location `json:"location"` // The location to which the supergroup is connected. It can't be a live location.
	Address  string   `json:"address"`  // Location address; 1-64 characters, as defined by the chat owner
}

// This object represents a forum topic.
type ForumTopic struct {
	MessageThreadID   int    `json:"message_thread_id"`              // Unique identifier of the forum topic
	Name              string `json:"name"`                           // Name of the topic
	IconColor         int    `json:"icon_color"`                     // Color of the topic icon in RGB format
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"` // Optional. Unique identifier of the custom emoji shown as the topic icon
}

// This object represents a bot command.
type BotCommand struct {
	Command     string `json:"command"`     // Text of the command; 1-32 characters. Can contain only lowercase English letters, digits, and underscores.
	Description string `json:"description"` // Description of the command; 1-256 characters.
}

// This object represents the scope to which bot commands are applied.
type BotCommandScope struct {
	Type   string      `json:"type"`              // Scope type, can be "default", "all_private_chats", "all_group_chats", "all_chat_administrators", "chat", "chat_administrators", "chat_member"
	ChatID interface{} `json:"chat_id,omitempty"` // (chat, chat_administrators, chat_member) Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)
	UserID int         `json:"user_id,omitempty"` // (chat_member) Unique identifier of the target user
}

// This object represents the bot's name.
type BotName struct {
	Name string `json:"name"` // The bot's name
}

// This object represents the bot's description.
type BotDescription struct {
	Description string `json:"description"` // The bot's description
}

// This object represents the bot's short description.
type BotShortDescription struct {
	ShortDescription string `json:"short_description"` // The bot's short description
}

// This object describes the bot's menu button in a private chat.
type MenuButton struct {
	Type   string      `json:"type"`              // Type of the button, can be "commands", "web_app" or "default"
	Text   string      `json:"text,omitempty"`    // (web_app) Text on the button
	WebApp *WebAppInfo `json:"web_app,omitempty"` // (web_app) Description of the Web App that will be launched when the user presses the button. The Web App will be able to send an arbitrary message on behalf of the user using the method answerWebAppQuery.
}

// Describes why a request was unsuccessful.
type ResponseParameters struct {
	MigrateToChatID int `json:"migrate_to_chat_id,omitempty"` // Optional. The group has been migrated to a supergroup with the specified identifier. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	RetryAfter      int `json:"retry_after,omitempty"`        // Optional. In case of exceeding flood control, the number of seconds left to wait before the request can be repeated
}

type InputMediaBase struct {
	Type            string          `json:"type"`                       // Type of the result.
	Media           RequestFileData `json:"media"`                      // File to send
	Caption         string          `json:"caption,omitempty"`          // Optional. Caption of the photo
	ParseMode       string          `json:"parse_mode,omitempty"`       // Optional. Mode for parsing entities in the caption
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"` // Optional. List of special entities in the caption

}

// This object represents the content of a media message to be sent
type InputMediaPhoto struct {
	InputMediaBase      // Type of the result, must be "photo"
	HasSpoiler     bool `json:"has_spoiler,omitempty"` // Optional. Whether the photo should be covered with a spoiler animation
}

// This object represents the content of a media message to be sent
type InputMediaVideo struct {
	InputMediaBase                    // Type of the result, must be "video"
	Thumbnail         RequestFileData `json:"thumbnail,omitempty"`          // Optional. Thumbnail of the video
	Width             int             `json:"width,omitempty"`              // Optional. Video width
	Height            int             `json:"height,omitempty"`             // Optional. Video height
	Duration          int             `json:"duration,omitempty"`           // Optional. Video duration in seconds
	SupportsStreaming bool            `json:"supports_streaming,omitempty"` // Optional. Whether the video is suitable for streaming
	HasSpoiler        bool            `json:"has_spoiler,omitempty"`        // Optional. Whether the video should be covered with a spoiler animation
}

// This object represents the content of a media message to be sent
type InputMediaAnimation struct {
	InputMediaBase                 // Type of the result, must be "animation"
	Thumbnail      RequestFileData `json:"thumbnail,omitempty"`   // Optional. Thumbnail of the animation
	Width          int             `json:"width,omitempty"`       // Optional. Animation width
	Height         int             `json:"height,omitempty"`      // Optional. Animation height
	Duration       int             `json:"duration,omitempty"`    // Optional. Animation duration in seconds
	HasSpoiler     bool            `json:"has_spoiler,omitempty"` // Optional. Whether the animation should be covered with a spoiler animation
}

// This object represents the content of a media message to be sent
type InputMediaAudio struct {
	InputMediaBase                 // Type of the result, must be "audio"
	Thumbnail      RequestFileData `json:"thumbnail,omitempty"` // Optional. Thumbnail of the audio
	Duration       int             `json:"duration,omitempty"`  // Optional. Audio duration in seconds
	Performer      string          `json:"performer,omitempty"` // Optional. Performer of the audio
	Title          string          `json:"title,omitempty"`     // Optional. Title of the audio
}

// This object represents the content of a media message to be sent
type InputMediaDocument struct {
	InputMediaBase                              // Type of the result, must be "document"
	Thumbnail                   RequestFileData `json:"thumbnail,omitempty"`                      // Optional. Thumbnail of the document
	DisableContentTypeDetection bool            `json:"disable_content_type_detection,omitempty"` // Optional. Disables automatic content type detection
}

//
//
//
// Stickers types
//
//
//

// This object represents a sticker.
type Sticker struct {
	FileID           string        `json:"file_id"`                     // Identifier for this file
	FileUniqueID     string        `json:"file_unique_id"`              // Unique identifier for this file
	Type             string        `json:"type"`                        // Type of the sticker: "regular", "mask", or "custom_emoji"
	Width            int           `json:"width"`                       // Sticker width
	Height           int           `json:"height"`                      // Sticker height
	IsAnimated       bool          `json:"is_animated"`                 // True, if the sticker is animated
	IsVideo          bool          `json:"is_video"`                    // True, if the sticker is a video sticker
	Thumbnail        *PhotoSize    `json:"thumbnail,omitempty"`         // Optional. Sticker thumbnail
	Emoji            string        `json:"emoji,omitempty"`             // Optional. Emoji associated with the sticker
	SetName          string        `json:"set_name,omitempty"`          // Optional. Name of the sticker set to which the sticker belongs
	PremiumAnimation *File         `json:"premium_animation,omitempty"` // Optional. Premium animation for premium regular stickers
	MaskPosition     *MaskPosition `json:"mask_position,omitempty"`     // Optional. Position where the mask should be placed for mask stickers
	CustomEmojiID    string        `json:"custom_emoji_id,omitempty"`   // Optional. Unique identifier of the custom emoji for custom emoji stickers
	NeedsRepainting  bool          `json:"needs_repainting,omitempty"`  // Optional. True, if the sticker must be repainted in certain contexts
	FileSize         int           `json:"file_size,omitempty"`         // Optional. File size in bytes
}

// This object represents a sticker set.
type StickerSet struct {
	Name        string     `json:"name"`                // Sticker set name
	Title       string     `json:"title"`               // Sticker set title
	StickerType string     `json:"sticker_type"`        // Type of stickers in the set: "regular", "mask", or "custom_emoji"
	IsAnimated  bool       `json:"is_animated"`         // True, if the sticker set contains animated stickers
	IsVideo     bool       `json:"is_video"`            // True, if the sticker set contains video stickers
	Stickers    []Sticker  `json:"stickers"`            // List of stickers in the set
	Thumbnail   *PhotoSize `json:"thumbnail,omitempty"` // Optional. Sticker set thumbnail
}

// This object describes the position on faces where a mask should be placed by default.
type MaskPosition struct {
	Point  string  `json:"point"`   // The part of the face relative to which the mask should be placed
	XShift float64 `json:"x_shift"` // Shift by X-axis measured in widths of the mask scaled to the face size
	YShift float64 `json:"y_shift"` // Shift by Y-axis measured in heights of the mask scaled to the face size
	Scale  float64 `json:"scale"`   // Mask scaling coefficient
}

// This object describes a sticker to be added to a sticker set.
type InputSticker struct {
	Sticker      RequestFileData `json:"sticker"`                 // The added sticker
	EmojiList    []string        `json:"emoji_list,omitempty"`    // Optional. List of emoji associated with the sticker
	MaskPosition *MaskPosition   `json:"mask_position,omitempty"` // Optional. Position where the mask should be placed for mask stickers
	Keywords     []string        `json:"keywords,omitempty"`      // Optional. List of search keywords for the sticker
}

//
//
//
// Inline types
//
//
//

// This object represents an incoming inline query.
// When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	ID       string    `json:"id"`                  // Unique identifier for this query
	From     *User     `json:"from"`                // Sender
	Query    string    `json:"query"`               // Text of the query (up to 256 characters)
	Offset   string    `json:"offset"`              // Offset of the results to be returned, can be controlled by the bot
	ChatType string    `json:"chat_type,omitempty"` // Optional. Type of the chat from which the inline query was sent
	Location *Location `json:"location,omitempty"`  // Optional. Sender location
}

// This object represents a button to be shown above inline query results. You must use exactly one of the optional fields.
type InlineQueryResultsButton struct {
	Text           string      `json:"text"`                      // Label text on the button
	WebApp         *WebAppInfo `json:"web_app,omitempty"`         // Optional. Description of the Web App to be launched when the button is pressed
	StartParameter string      `json:"start_parameter,omitempty"` // Optional. Deep-linking parameter for the /start message sent to the bot when the button is pressed
}

type InlineQueryResultBase struct {
	Type string `json:"type"` // Type of the result
	ID   string `json:"id"`   // Unique identifier for this result, 1-64 Bytes

}

// Represents a link to an article or web page.
type InlineQueryResultArticle struct {
	InlineQueryResultBase                       // Type of the result, must be article
	Title                 string                `json:"title"`                      // Title of the result
	InputMessageContent   interface{}           `json:"input_message_content"`      // Content of the message to be sent
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`     // Optional. Inline keyboard attached to the message
	URL                   string                `json:"url,omitempty"`              // Optional. URL of the result
	HideURL               bool                  `json:"hide_url,omitempty"`         // Optional. Pass True if you don't want the URL to be shown in the message
	Description           string                `json:"description,omitempty"`      // Optional. Short description of the result
	ThumbnailURL          string                `json:"thumbnail_url,omitempty"`    // Optional. URL of the thumbnail for the result
	ThumbnailWidth        int                   `json:"thumbnail_width,omitempty"`  // Optional. Thumbnail width
	ThumbnailHeight       int                   `json:"thumbnail_height,omitempty"` // Optional. Thumbnail height
}

// Represents a link to a photo. By default, this photo will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultPhoto struct {
	InlineQueryResultBase                       // Type of the result, must be photo
	URL                   string                `json:"photo_url"`                       // A valid URL of the photo
	ThumbnailURL          string                `json:"thumbnail_url"`                   // URL of the thumbnail for the photo
	Width                 int                   `json:"photo_width,omitempty"`           // Optional. Width of the photo
	Height                int                   `json:"photo_height,omitempty"`          // Optional. Height of the photo
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the photo caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the photo
}

// Represents a link to an animated GIF file. By default, this animated GIF file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultGif struct {
	InlineQueryResultBase                       // Type of the result, must be gif
	URL                   string                `json:"gif_url"`                         // A valid URL for the GIF file
	Width                 int                   `json:"gif_width,omitempty"`             // Optional. Width of the GIF
	Height                int                   `json:"gif_height,omitempty"`            // Optional. Height of the GIF
	Duration              int                   `json:"gif_duration,omitempty"`          // Optional. Duration of the GIF in seconds
	ThumbnailURL          string                `json:"thumbnail_url"`                   // URL of the thumbnail for the result
	ThumbnailMimeType     string                `json:"thumbnail_mime_type,omitempty"`   // Optional. MIME type of the thumbnail
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the GIF file to be sent
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the GIF animation
}

// Represents a link to a video animation (H.264/MPEG-4 AVC video without sound).
// By default, this animated MPEG-4 file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultMpeg4Gif struct {
	InlineQueryResultBase                       // Type of the result, must be mpeg4_gif
	URL                   string                `json:"mpeg4_url"`                       // A valid URL for the MPEG4 file
	Width                 int                   `json:"mpeg4_width,omitempty"`           // Optional. Video width
	Height                int                   `json:"mpeg4_height,omitempty"`          // Optional. Video height
	Duration              int                   `json:"mpeg4_duration,omitempty"`        // Optional. Video duration in seconds
	ThumbnailURL          string                `json:"thumbnail_url"`                   // URL of the thumbnail for the result
	ThumbnailMimeType     string                `json:"thumbnail_mime_type,omitempty"`   // Optional. MIME type of the thumbnail
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the MPEG-4 file to be sent
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video animation
}

// Represents a link to a page containing an embedded video player or a video file.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultVideo struct {
	InlineQueryResultBase                       // Type of the result, must be video
	URL                   string                `json:"video_url"`                       // A valid URL for the embedded video player or video file
	MimeType              string                `json:"mime_type"`                       // MIME type of the content of the video URL, "text/html" or "video/mp4"
	ThumbnailURL          string                `json:"thumbnail_url"`                   // URL of the thumbnail (JPEG only) for the video
	Title                 string                `json:"title"`                           // Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the video caption. See formatting options for more details.
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	Width                 int                   `json:"video_width,omitempty"`           // Optional. Video width
	Height                int                   `json:"video_height,omitempty"`          // Optional. Video height
	Duration              int                   `json:"video_duration,omitempty"`        // Optional. Video duration in seconds
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video
}

// Represents a link to an MP3 audio file. By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultAudio struct {
	InlineQueryResultBase                       // Type of the result, must be audio
	URL                   string                `json:"audio_url"`                       // A valid URL for the audio file
	Title                 string                `json:"title"`                           // Title
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the audio caption. See formatting options for more details.
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	Performer             string                `json:"performer,omitempty"`             // Optional. Performer
	Duration              int                   `json:"audio_duration,omitempty"`        // Optional. Audio duration in seconds
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the audio
}

// Represents a link to a voice recording in an .OGG container encoded with OPUS.
// By default, this voice recording will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the the voice message.
type InlineQueryResultVoice struct {
	InlineQueryResultBase                       // Type of the result, must be voice
	URL                   string                `json:"voice_url"`                       // A valid URL for the voice recording
	Title                 string                `json:"title"`                           // Recording title
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the voice message caption. See formatting options for more details.
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	Duration              int                   `json:"voice_duration,omitempty"`        // Optional. Recording duration in seconds
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the voice recording
}

// Represents a link to a file. By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
// Currently, only .PDF and .ZIP files can be sent using this method.
type InlineQueryResultDocument struct {
	InlineQueryResultBase                       // Type of the result, must be document
	Title                 string                `json:"title"`                           // Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the document caption. See formatting options for more details.
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption, which can be specified instead of parse_mode
	URL                   string                `json:"document_url"`                    // A valid URL for the file
	MimeType              string                `json:"mime_type"`                       // MIME type of the content of the file, either "application/pdf" or "application/zip"
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the file
	ThumbnailURL          string                `json:"thumbnail_url,omitempty"`         // Optional. URL of the thumbnail (JPEG only) for the file
	ThumbnailWidth        int                   `json:"thumbnail_width,omitempty"`       // Optional. Thumbnail width
	ThumbnailHeight       int                   `json:"thumbnail_height,omitempty"`      // Optional. Thumbnail height
}

// Represents a location on a map. By default, the location will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the location.
type InlineQueryResultLocation struct {
	InlineQueryResultBase                       // Type of the result, must be "location"
	Latitude              float64               `json:"latitude"`                         // Location latitude in degrees
	Longitude             float64               `json:"longitude"`                        // Location longitude in degrees
	Title                 string                `json:"title"`                            // Location title
	HorizontalAccuracy    float64               `json:"horizontal_accuracy,omitempty"`    // Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	LivePeriod            int                   `json:"live_period,omitempty"`            // Optional. Period in seconds for which the location can be updated, should be between 60 and 86400
	Heading               int                   `json:"heading,omitempty"`                // Optional. For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.
	ProximityAlertRadius  int                   `json:"proximity_alert_radius,omitempty"` // Optional. For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`           // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"`  // Optional. Content of the message to be sent instead of the location
	ThumbnailURL          string                `json:"thumbnail_url,omitempty"`          // Optional. URL of the thumbnail for the result
	ThumbnailWidth        int                   `json:"thumbnail_width,omitempty"`        // Optional. Thumbnail width
	ThumbnailHeight       int                   `json:"thumbnail_height,omitempty"`       // Optional. Thumbnail height
}

// Represents a venue. By default, the venue will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the venue.
type InlineQueryResultVenue struct {
	InlineQueryResultBase                       // Type of the result, must be "venue"
	Latitude              float64               `json:"latitude"`                        // Latitude of the venue location in degrees
	Longitude             float64               `json:"longitude"`                       // Longitude of the venue location in degrees
	Title                 string                `json:"title"`                           // Title of the venue
	Address               string                `json:"address"`                         // Address of the venue
	FoursquareID          string                `json:"foursquare_id,omitempty"`         // Optional. Foursquare identifier of the venue if known
	FoursquareType        string                `json:"foursquare_type,omitempty"`       // Optional. Foursquare type of the venue, if known
	GooglePlaceID         string                `json:"google_place_id,omitempty"`       // Optional. Google Places identifier of the venue
	GooglePlaceType       string                `json:"google_place_type,omitempty"`     // Optional. Google Places type of the venue
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the venue
	ThumbnailURL          string                `json:"thumbnail_url,omitempty"`         // Optional. URL of the thumbnail for the result
	ThumbnailWidth        int                   `json:"thumbnail_width,omitempty"`       // Optional. Thumbnail width
	ThumbnailHeight       int                   `json:"thumbnail_height,omitempty"`      // Optional. Thumbnail height
}

// Represents a contact with a phone number. By default, this contact will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the contact.
type InlineQueryResultContact struct {
	InlineQueryResultBase                       // Type of the result, must be "contact"
	PhoneNumber           string                `json:"phone_number"`                    // Contact's phone number
	FirstName             string                `json:"first_name"`                      // Contact's first name
	LastName              string                `json:"last_name,omitempty"`             // Optional. Contact's last name
	VCard                 string                `json:"vcard,omitempty"`                 // Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the contact
	ThumbnailURL          string                `json:"thumbnail_url,omitempty"`         // Optional. URL of the thumbnail for the result
	ThumbnailWidth        int                   `json:"thumbnail_width,omitempty"`       // Optional. Thumbnail width
	ThumbnailHeight       int                   `json:"thumbnail_height,omitempty"`      // Optional. Thumbnail height
}

// Represents a Game.
type InlineQueryResultGame struct {
	InlineQueryResultBase                       // Type of the result, must be "game"
	GameShortName         string                `json:"game_short_name"`        // Short name of the game
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"` // Optional. Inline keyboard attached to the message
}

// Represents a link to a photo stored on the Telegram servers.
// By default, this photo will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultCachedPhoto struct {
	InlineQueryResultBase                       // Type of the result, must be "photo"
	PhotoFileID           string                `json:"photo_file_id"`                   // A valid file identifier of the photo
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the photo to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the photo caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the photo
}

// Represents a link to an animated GIF file stored on the Telegram servers.
// By default, this animated GIF file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with specified content instead of the animation.
type InlineQueryResultCachedGif struct {
	InlineQueryResultBase                       // Type of the result, must be "gif"
	GifFileID             string                `json:"gif_file_id"`                     // A valid file identifier for the GIF file
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the GIF file to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the GIF animation
}

// Represents a link to a video animation (H.264/MPEG-4 AVC video without sound) stored on the Telegram servers.
// By default, this animated MPEG-4 file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultCachedMpeg4Gif struct {
	InlineQueryResultBase                       // Type of the result, must be "mpeg4_gif"
	Mpeg4FileID           string                `json:"mpeg4_file_id"`                   // A valid file identifier for the MPEG4 file
	Title                 string                `json:"title,omitempty"`                 // Optional. Title for the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the MPEG4 file to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video animation
}

// Represents a link to a sticker stored on the Telegram servers.
// By default, this sticker will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the sticker.
type InlineQueryResultCachedSticker struct {
	InlineQueryResultBase                       // Type of the result, must be "sticker"
	StickerFileID         string                `json:"sticker_file_id"`                 // A valid file identifier of the sticker
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the sticker
}

// Represents a link to a file stored on the Telegram servers.
// By default, this file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the file.
type InlineQueryResultCachedDocument struct {
	InlineQueryResultBase                       // Type of the result, must be "document"
	Title                 string                `json:"title"`                           // Title for the result
	DocumentFileID        string                `json:"document_file_id"`                // A valid file identifier for the file
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the document to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the document caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the file
}

// Represents a link to a video file stored on the Telegram servers.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultCachedVideo struct {
	InlineQueryResultBase                       // Type of the result, must be "video"
	VideoFileID           string                `json:"video_file_id"`                   // A valid file identifier for the video file
	Title                 string                `json:"title"`                           // Title for the result
	Description           string                `json:"description,omitempty"`           // Optional. Short description of the result
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption of the video to be sent, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the video caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video
}

// Represents a link to a voice message stored on the Telegram servers.
// By default, this voice message will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the voice message.
type InlineQueryResultCachedVoice struct {
	InlineQueryResultBase                       // Type of the result, must be "voice"
	VoiceFileID           string                `json:"voice_file_id"`                   // A valid file identifier for the voice message
	Title                 string                `json:"title"`                           // Voice message title
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the voice message caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the voice message
}

// Represents a link to an MP3 audio file stored on the Telegram servers.
// By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultCachedAudio struct {
	InlineQueryResultBase                       // Type of the result, must be "audio"
	AudioFileID           string                `json:"audio_file_id"`                   // A valid file identifier for the audio file
	Caption               string                `json:"caption,omitempty"`               // Optional. Caption, 0-1024 characters after entities parsing
	ParseMode             string                `json:"parse_mode,omitempty"`            // Optional. Mode for parsing entities in the audio caption
	CaptionEntities       []MessageEntity       `json:"caption_entities,omitempty"`      // Optional. List of special entities that appear in the caption
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent   interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the audio
}

// This object represents the content of a message to be sent as a result of an inline query.
type InputTextMessageContent struct {
	MessageText           string          `json:"message_text"`                       // Text of the message to be sent, 1-4096 characters
	ParseMode             string          `json:"parse_mode,omitempty"`               // Optional. Mode for parsing entities in the message text. See formatting options for more details.
	Entities              []MessageEntity `json:"entities,omitempty"`                 // Optional. List of special entities that appear in message text, which can be specified instead of parse_mode
	DisableWebPagePreview bool            `json:"disable_web_page_preview,omitempty"` // Optional. Disables link previews for links in the sent message
}

// This object represents the content of a message to be sent as a result of an inline query.
type InputLocationMessageContent struct {
	Latitude             float64 `json:"latitude"`                         // Latitude of the location in degrees
	Longitude            float64 `json:"longitude"`                        // Longitude of the location in degrees
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`    // Optional. The radius of uncertainty for the location, measured in meters; 0-1500
	LivePeriod           int     `json:"live_period,omitempty"`            // Optional. Period in seconds for which the location can be updated, should be between 60 and 86400.
	Heading              int     `json:"heading,omitempty"`                // Optional. For live locations, a direction in which the user is moving, in degrees. Must be between 1 and 360 if specified.
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"` // Optional. For live locations, a maximum distance for proximity alerts about approaching another chat member, in meters. Must be between 1 and 100000 if specified.
}

// This object represents the content of a message to be sent as a result of an inline query.
type InputVenueMessageContent struct {
	Latitude        float64 `json:"latitude"`                    // Latitude of the venue in degrees
	Longitude       float64 `json:"longitude"`                   // Longitude of the venue in degrees
	Title           string  `json:"title"`                       // Name of the venue
	Address         string  `json:"address"`                     // Address of the venue
	FoursquareID    string  `json:"foursquare_id,omitempty"`     // Optional. Foursquare identifier of the venue, if known
	FoursquareType  string  `json:"foursquare_type,omitempty"`   // Optional. Foursquare type of the venue, if known. (For example, “arts_entertainment/default”, “arts_entertainment/aquarium” or “food/icecream”.)
	GooglePlaceID   string  `json:"google_place_id,omitempty"`   // Optional. Google Places identifier of the venue
	GooglePlaceType string  `json:"google_place_type,omitempty"` // Optional. Google Places type of the venue. (See supported types.)
}

// This object represents the content of a message to be sent as a result of an inline query.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`        // Contact's phone number
	FirstName   string `json:"first_name"`          // Contact's first name
	LastName    string `json:"last_name,omitempty"` // Optional. Contact's last name
	VCard       string `json:"vcard,omitempty"`     // Optional. Additional data about the contact in the form of a vCard, 0-2048 bytes
}

// This object represents the content of a message to be sent as a result of an inline query.
type InputInvoiceMessageContent struct {
	Title                     string         `json:"title"`                                   // Product name, 1-32 characters
	Description               string         `json:"description"`                             // Product description, 1-255 characters
	Payload                   string         `json:"payload"`                                 // Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use for your internal processes.
	ProviderToken             string         `json:"provider_token"`                          // Payment provider token, obtained via @BotFather
	Currency                  string         `json:"currency"`                                // Three-letter ISO 4217 currency code, see more on currencies
	Prices                    []LabeledPrice `json:"prices"`                                  // Price breakdown, a JSON-serialized list of components (e.g. product price, tax, discount, delivery cost, delivery tax, bonus, etc.)
	MaxTipAmount              int            `json:"max_tip_amount,omitempty"`                // Optional. The maximum accepted amount for tips in the smallest units of the currency (integer, not float/double). For example, for a maximum tip of US$ 1.45 pass max_tip_amount = 145. See the exp parameter in currencies.json, it shows the number of digits past the decimal point for each currency (2 for the majority of currencies). Defaults to 0
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts,omitempty"`         // Optional. A JSON-serialized array of suggested amounts of tip in the smallest units of the currency (integer, not float/double). At most 4 suggested tip amounts can be specified. The suggested tip amounts must be positive, passed in a strictly increased order and must not exceed max_tip_amount.
	ProviderData              string         `json:"provider_data,omitempty"`                 // Optional. A JSON-serialized object for data about the invoice, which will be shared with the payment provider. A detailed description of the required fields should be provided by the payment provider.
	PhotoURL                  string         `json:"photo_url,omitempty"`                     // Optional. URL of the product photo for the invoice. Can be a photo of the goods or a marketing image for a service.
	PhotoSize                 int            `json:"photo_size,omitempty"`                    // Optional. Photo size in bytes
	PhotoWidth                int            `json:"photo_width,omitempty"`                   // Optional. Photo width
	PhotoHeight               int            `json:"photo_height,omitempty"`                  // Optional. Photo height
	NeedName                  bool           `json:"need_name,omitempty"`                     // Optional. Pass True if you require the user's full name to complete the order
	NeedPhoneNumber           bool           `json:"need_phone_number,omitempty"`             // Optional. Pass True if you require the user's phone number to complete the order
	NeedEmail                 bool           `json:"need_email,omitempty"`                    // Optional. Pass True if you require the user's email address to complete the order
	NeedShippingAddress       bool           `json:"need_shipping_address,omitempty"`         // Optional. Pass True if you require the user's shipping address to complete the order
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider,omitempty"` // Optional. Pass True if the user's phone number should be sent to provider
	SendEmailToProvider       bool           `json:"send_email_to_provider,omitempty"`        // Optional. Pass True if the user's email address should be sent to provider
	IsFlexible                bool           `json:"is_flexible,omitempty"`                   // Optional. Pass True if the final price depends on the shipping method
}

// Represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`                   // The unique identifier for the result that was chosen
	From            *User     `json:"from"`                        // The user that chose the result
	Location        *Location `json:"location,omitempty"`          // Optional. Sender location (for bots that require user location)
	InlineMessageID string    `json:"inline_message_id,omitempty"` // Optional. Identifier of the sent inline message
	Query           string    `json:"query"`                       // The query that was used to obtain the result
}

// Describes an inline message sent by a Web App on behalf of a user.
type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id,omitempty"` // Optional. Identifier of the sent inline message
}

//
//
//
// Payment types
//
//
//

// This object represents a portion of the price for goods or services.
type LabeledPrice struct {
	Label  string `json:"label"`  // Portion label
	Amount int    `json:"amount"` // Price of the product in the smallest units of the currency
}

// This object contains basic information about an invoice.
type Invoice struct {
	Title          string `json:"title"`           // Product name
	Description    string `json:"description"`     // Product description
	StartParameter string `json:"start_parameter"` // Unique bot deep-linking parameter
	Currency       string `json:"currency"`        // Three-letter ISO 4217 currency code
	TotalAmount    int    `json:"total_amount"`    // Total price in the smallest units of the currency
}

// This object represents a shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"` // Two-letter ISO 3166-1 alpha-2 country code
	State       string `json:"state"`        // State, if applicable
	City        string `json:"city"`         // City
	StreetLine1 string `json:"street_line1"` // First line for the address
	StreetLine2 string `json:"street_line2"` // Second line for the address
	PostCode    string `json:"post_code"`    // Address post code
}

// This object represents information about an order.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`             // Optional. User name
	PhoneNumber     string           `json:"phone_number,omitempty"`     // Optional. User's phone number
	Email           string           `json:"email,omitempty"`            // Optional. User email
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"` // Optional. User shipping address
}

// This object represents one shipping option.
type ShippingOption struct {
	ID     string         `json:"id"`     // Shipping option identifier
	Title  string         `json:"title"`  // Option title
	Prices []LabeledPrice `json:"prices"` // List of price portions
}

// This object contains basic information about a successful payment.
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`                     // Three-letter ISO 4217 currency code
	TotalAmount             int        `json:"total_amount"`                 // Total price in the smallest units of the currency
	InvoicePayload          string     `json:"invoice_payload"`              // Bot specified invoice payload
	ShippingOptionID        string     `json:"shipping_option_id,omitempty"` // Optional. Identifier of the shipping option chosen by the user
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`         // Optional. Order information provided by the user
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`   // Telegram payment identifier
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`   // Provider payment identifier
}

// This object contains information about an incoming shipping query.
type ShippingQuery struct {
	ID              string          `json:"id"`               // Unique query identifier
	From            *User           `json:"from"`             // User who sent the query
	InvoicePayload  string          `json:"invoice_payload"`  // Bot specified invoice payload
	ShippingAddress ShippingAddress `json:"shipping_address"` // User specified shipping address
}

// This object contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`                           // Unique query identifier
	From             *User      `json:"from"`                         // User who sent the query
	Currency         string     `json:"currency"`                     // Three-letter ISO 4217 currency code
	TotalAmount      int        `json:"total_amount"`                 // Total price in the smallest units of the currency
	InvoicePayload   string     `json:"invoice_payload"`              // Bot specified invoice payload
	ShippingOptionID string     `json:"shipping_option_id,omitempty"` // Optional. Identifier of the shipping option chosen by the user
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`         // Optional. Order information provided by the user
}

//
//
//
// Passport types
//
//
//

// Describes Telegram Passport data shared with the bot by the user.
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`        // Array with information about documents and other Telegram Passport elements shared with the bot
	Credentials EncryptedCredentials       `json:"credentials"` // Encrypted credentials required to decrypt the data
}

// This object represents a file uploaded to Telegram Passport.
// Currently all Telegram Passport files are in JPEG format when decrypted and don't exceed 10MB.
type PassportFile struct {
	FileID       string `json:"file_id"`        // Identifier for this file, which can be used to download or reuse the file
	FileUniqueID string `json:"file_unique_id"` // Unique identifier for this file, which is supposed to be the same over time and for different bots
	FileSize     int    `json:"file_size"`      // File size in bytes
	FileDate     int    `json:"file_date"`      // Unix time when the file was uploaded
}

// Describes documents or other Telegram Passport elements shared with the bot by the user.
type EncryptedPassportElement struct {
	Type        string         `json:"type"`                   // Element type
	Data        string         `json:"data,omitempty"`         // Optional. Base64-encoded encrypted Telegram Passport element data
	PhoneNumber string         `json:"phone_number,omitempty"` // Optional. User's verified phone number
	Email       string         `json:"email,omitempty"`        // Optional. User's verified email address
	Files       []PassportFile `json:"files,omitempty"`        // Optional. Array of encrypted files with documents provided by the user
	FrontSide   *PassportFile  `json:"front_side,omitempty"`   // Optional. Encrypted file with the front side of the document
	ReverseSide *PassportFile  `json:"reverse_side,omitempty"` // Optional. Encrypted file with the reverse side of the document
	Selfie      *PassportFile  `json:"selfie,omitempty"`       // Optional. Encrypted file with the user holding a document selfie
	Translation []PassportFile `json:"translation,omitempty"`  // Optional. Array of encrypted files with translated versions of documents
	Hash        string         `json:"hash"`                   // Base64-encoded element hash for using in PassportElementErrorUnspecified
}

// Describes data required for decrypting and authenticating EncryptedPassportElement.
// See the Telegram Passport Documentation for a complete description of the data decryption and authentication processes.
type EncryptedCredentials struct {
	Data   string `json:"data"`   // Base64-encoded encrypted JSON-serialized data with user's payload, data hashes, and secrets
	Hash   string `json:"hash"`   // Base64-encoded data hash for data authentication
	Secret string `json:"secret"` // Base64-encoded secret encrypted with the bot's public RSA key
}

type PassportElementErrorBase struct {
	Source string `json:"source"` // Error source.
	Type   string `json:"type"`   // The section of the user's Telegram Passport which has the error

}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorDataField struct {
	PassportElementErrorBase        // Error source, must be "data"
	FieldName                string `json:"field_name"` // Name of the data field which has the error
	DataHash                 string `json:"data_hash"`  // Base64-encoded data hash
	Message                  string `json:"message"`    // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorFrontSide struct {
	PassportElementErrorBase        // Error source, must be "front_side"
	FileHash                 string `json:"file_hash"` // Base64-encoded hash of the file with the front side of the document
	Message                  string `json:"message"`   // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorReverseSide struct {
	PassportElementErrorBase        // Error source, must be "reverse_side"
	FileHash                 string `json:"file_hash"` // Base64-encoded hash of the file with the reverse side of the document
	Message                  string `json:"message"`   // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorSelfie struct {
	PassportElementErrorBase        // Error source, must be "selfie"
	FileHash                 string `json:"file_hash"` // Base64-encoded hash of the file with the selfie
	Message                  string `json:"message"`   // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorFile struct {
	PassportElementErrorBase        // Error source, must be "file"
	FileHash                 string `json:"file_hash"` // Base64-encoded file hash
	Message                  string `json:"message"`   // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorFiles struct {
	PassportElementErrorBase          // Error source, must be "files"
	FileHashes               []string `json:"file_hashes"` // List of base64-encoded file hashes
	Message                  string   `json:"message"`     // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorTranslationFile struct {
	PassportElementErrorBase        // Error source, must be "translation_file"
	FileHash                 string `json:"file_hash"` // Base64-encoded file hash
	Message                  string `json:"message"`   // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorTranslationFiles struct {
	PassportElementErrorBase          // Error source, must be "translation_files"
	FileHashes               []string `json:"file_hashes"` // List of base64-encoded file hashes
	Message                  string   `json:"message"`     // Error message
}

// This object represents an error in the Telegram Passport element which was submitted that should be resolved by the user.
type PassportElementErrorUnspecified struct {
	PassportElementErrorBase        // Error source, must be "unspecified"
	ElementHash              string `json:"element_hash"` // Base64-encoded element hash
	Message                  string `json:"message"`      // Error message
}

//
//
//
// Game types
//
//
//

// This object represents a game. Use BotFather to create and edit games, their short names will act as unique identifiers.
type Game struct {
	Title        string          `json:"title"`                   // Title of the game
	Description  string          `json:"description"`             // Description of the game
	Photo        []PhotoSize     `json:"photo"`                   // Photo that will be displayed in the game message in chats
	Text         string          `json:"text,omitempty"`          // Optional. Brief description of the game or high scores included in the game message
	TextEntities []MessageEntity `json:"text_entities,omitempty"` // Optional. Special entities that appear in the text
	Animation    *Animation      `json:"animation,omitempty"`     // Optional. Animation that will be displayed in the game message in chats
}

// A placeholder, currently holds no information. Use BotFather to set up your game.
type CallbackGame struct{}

// This object represents one row of the high scores table for a game.
type GameHighScore struct {
	Position int  `json:"position"` // Position in the high score table for the game
	User     User `json:"user"`     // User
	Score    int  `json:"score"`    // Score
}
