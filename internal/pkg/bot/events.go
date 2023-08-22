package bot

import (
	"strconv"
	"strings"
	"telegram-bot-feedback/internal/pkg/database"
	l "telegram-bot-feedback/internal/pkg/logger"
	tg "telegram-bot-feedback/pkg/telegram-bot-api"
	"time"
)

// Button Sets
const (
	UserMain int = iota + 1
	UserStars
	UserClose
	EmplMain
	EmplMainR
	EmplReview
	EmplExit
)

// buttons returns button texts for ReplyKeyboardMarkup
func buttons(key int) []string {
	switch key {
	case UserMain:
		return []string{"⭐Review", "❓Question"}
	case UserStars:
		return []string{"⭐⭐⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐", "⭐⭐", "⭐"}
	case UserClose:
		return []string{"❌Close"}
	case EmplMain:
		return []string{"❓Receive questions", "❓Open questions", "❓Find a question", "⭐Reviews"}
	case EmplMainR:
		return []string{"❓Do not receive questions", "❓Open questions", "❓Find a question", "⭐Reviews"}
	case EmplReview:
		return []string{"📅For a day", "📅For a week", "📅For a month", "📅All (no text)", "↩️Back"}
	case EmplExit:
		return []string{"↩️Back"}
	}
	return []string{}
}

// responserCommand responds to commands
func responserCommand(command string, user *database.User, app *App) error {
	if user.IsEmployee {
		return l.Err(responserCommandEmployee(command, user, app))
	}
	return l.Err(responserCommandUser(command, user, app))
}

// responserCommandUser responds to user commands
func responserCommandUser(command string, user *database.User, app *App) error {
	switch command {
	case "/start":
		message := tg.NewMessage(user.ChatID, "Greetings 👋\nWith my help, you can leave a \"⭐Review\" \nor ask a \"❓Question\"")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(UserMain)...)
		_, err := app.Bot.Send(message)
		if err != nil {
			return l.Err(err)
		}
		err = database.ChangeUserState(SMain, user, app.DB)
		return l.Err(err)
	}
	return nil
}

// responserCommandEmployee responds to employee commands
func responserCommandEmployee(command string, user *database.User, app *App) error {
	switch command {
	case "/start":
		message := tg.NewMessage(user.ChatID, "Greetings 👋\nI implement customer feedback\no receive questions click\n\"❓Receive questions\"")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplMain)...)
		_, err := app.Bot.Send(message)
		if err != nil {
			return l.Err(err)
		}
		err = database.ChangeUserState(SMain, user, app.DB)
		return l.Err(err)
	}
	return nil
}

// responser responds to message
func responser(user *database.User, app *App) error {
	if user.IsEmployee {
		return l.Err(responserEmployee(user, app))
	}
	return l.Err(responserUser(user, app))
}

// responserUser responds to user message
func responserUser(user *database.User, app *App) error {
	switch user.State {
	case SMain:
		message := tg.NewMessage(user.ChatID, "If you have any questions or review, I'm listening carefully")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(UserMain)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SReview:
		message := tg.NewMessage(user.ChatID, "Please rate from 1 to 5")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(UserStars)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SReviewText:
		message := tg.NewMessage(user.ChatID, "Thank you for your review\nYou can also leave a comment\nOr press \"❌Close\"")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(UserClose)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SQuestion:
		message := tg.NewMessage(user.ChatID, "Please ask your question\nOr click \"❌Close\"")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(UserClose)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SQuestionDiscussion:
		question := database.GetOpenQuestionByUser(user, app.DB)
		if question == nil {
			message := tg.NewMessage(user.ChatID, "Please reopen the question")
			_, err := app.Bot.Send(message)
			return l.Err(err)
		}
		id := strconv.Itoa(int(question.ID))
		message := tg.NewMessage(user.ChatID, "Your question #"+id+"\nThank you for your question\nAn available employee will answer you shortly")
		_, err := app.Bot.Send(message)
		return l.Err(err)
	}
	return nil
}

// responserEmployee responds to employee message
func responserEmployee(user *database.User, app *App) error {
	switch user.State {
	case SMain:
		message := tg.NewMessage(user.ChatID, "Choose an action")
		if user.IsReceiver {
			message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplMainR)...)
		} else {
			message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplMain)...)
		}
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SReview:
		message := tg.NewMessage(user.ChatID, "Select Interval")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplReview)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SQuestion:
		message := tg.NewMessage(user.ChatID, "No questions")
		_, err := app.Bot.Send(message)
		if err != nil {
			return l.Err(err)
		}
		err = database.ChangeUserState(SMain, user, app.DB)
		return l.Err(err)
	case SQuestionDiscussion:
		message := tg.NewMessage(user.ChatID, "You have entered a chat with a user")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplExit)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	case SSwitchReceiver:
		err := database.ChangeUserState(SMain, user, app.DB)
		if err != nil {
			return l.Err(err)
		}
		if user.IsReceiver {
			message := tg.NewMessage(user.ChatID, "You no longer receive questions")
			message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplMain)...)
			err = database.ChangeUserIsReceiver(false, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			_, err := app.Bot.Send(message)
			return l.Err(err)
		}
		message := tg.NewMessage(user.ChatID, "Now You receive questions")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplMainR)...)
		err = database.ChangeUserIsReceiver(true, user, app.DB)
		if err != nil {
			return l.Err(err)
		}
		_, err = app.Bot.Send(message)
		return l.Err(err)
	case SSearchQuestion:
		message := tg.NewMessage(user.ChatID, "Enter question number")
		message.ReplyMarkup = newReplyKeyboardMarkup(buttons(EmplExit)...)
		_, err := app.Bot.Send(message)
		return l.Err(err)
	}
	return nil
}

// newReplyKeyboardMarkup returns ReplyKeyboardMarkup by buttons texts
func newReplyKeyboardMarkup(text ...string) tg.ReplyKeyboardMarkup {
	var keyboard [][]tg.KeyboardButton

	for _, row := range text {
		keyboard = append(keyboard, []tg.KeyboardButton{tg.NewKeyboardButton(row)})
	}

	rkm := tg.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
	return rkm
}

// newOneButtonInlineKeyboardMarkup returns InlineKeyboardMarkup with one button by button text
func newOneButtonInlineKeyboardMarkup(text, data string) tg.InlineKeyboardMarkup {
	var keyboard [][]tg.InlineKeyboardButton
	keyboard = append(keyboard, []tg.InlineKeyboardButton{tg.NewInlineKeyboardButtonData(text, data)})
	rkm := tg.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
	return rkm
}

// sendQuestions sends Questions to the chat
func sendQuestions(to *database.User, bot *tg.Client, question []database.Question) error {
	for _, q := range question {
		id := strconv.Itoa(int(q.ID))
		key := strconv.Itoa(CBQuestion) + "-"
		text := "Question #" + id + "\n" + q.Header
		message := tg.NewMessage(to.ChatID, text)
		message.ReplyMarkup = newOneButtonInlineKeyboardMarkup("Take question", key+id)
		_, err := bot.Send(message)
		if err != nil {
			return l.Err(err)
		}
	}
	return nil
}

// sendCorrespondenceFromUser forwarding message from user to employee
func sendCorrespondenceFromUser(question *database.Question, message *tg.Message, bot *tg.Client) error {
	copy := tg.NewForward(question.Answerer.ChatID, question.User.ChatID, message.MessageID)
	_, err := bot.Send(copy)
	return l.Err(err)
}

// sendCorrespondenceFromAnswerer sends copy of message from employee to user
func sendCorrespondenceFromAnswerer(question *database.Question, message *tg.Message, bot *tg.Client) error {
	copy := tg.NewCopyMessage(question.User.ChatID, question.Answerer.ChatID, message.MessageID)
	_, err := bot.Send(copy)
	return l.Err(err)
}

// loadCorrespondence loads Correspondence to the chat by Question ID
func loadCorrespondence(id int, user *database.User, app *App) error {
	question := database.GetNewQuestionById(id, app.DB)
	if question == nil {
		message := tg.NewMessage(user.ChatID, "Question already taken")
		_, err := app.Bot.Send(message)
		return l.Err(err)
	}
	err := database.ChangeQuestionAnswerer(int(user.ID), question, app.DB)
	if err != nil {
		return l.Err(err)
	}
	correspondence := database.GetCorrespondenceByQuestion(question, app.DB)
	for _, corr := range correspondence {
		copy := tg.NewForward(user.ChatID, corr.User.ChatID, corr.MessageID)
		_, err := app.Bot.Send(copy)
		if err != nil {
			return l.Err(err)
		}
	}
	return nil
}

// loadReviews loads Reviews by date interval
func loadReviews(interval int, user *database.User, app *App) {
	fDate := time.Now().UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
	sDate := time.Time{}
	switch interval {
	case RDay:
		sDate = fDate.AddDate(0, 0, -1)
	case RWeek:
		sDate = fDate.AddDate(0, 0, -7)
	case RMonth:
		sDate = fDate.AddDate(0, -1, 0)
	case RAll:
		var text string
		for i, r := range database.GetCountReviewsByRating(app.DB) {
			text = text + ratingInStars(i+1) + " - " + strconv.Itoa(int(r)) + "\n"
		}
		message := tg.NewMessage(user.ChatID, text)
		app.Bot.Send(message)
		return
	}
	reviews := database.GetReviewsInRange(fDate, sDate, app.DB)
	if len(reviews) == 0 {
		return
	}
	for _, r := range reviews {
		message := tg.NewMessage(user.ChatID, ratingInStars(r.Rating)+"\n"+r.Text)
		app.Bot.Send(message)
	}
}

// loadFullQuestionById loads Question to the chat by ID
func loadFullQuestionById(id string, user *database.User, app *App) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		message := tg.NewMessage(user.ChatID, "Wrong format")
		app.Bot.Send(message)
		return
	}
	question := database.GetQuestionById(idInt, app.DB)
	if question == nil {
		message := tg.NewMessage(user.ChatID, "Question not found")
		app.Bot.Send(message)
		return
	}
	message := tg.NewMessage(user.ChatID, question.Header)
	app.Bot.Send(message)
	correspondence := database.GetCorrespondenceByQuestion(question, app.DB)
	for _, corr := range correspondence {
		copy := tg.NewForward(user.ChatID, corr.User.ChatID, corr.MessageID)
		_, err := app.Bot.Send(copy)
		if err != nil {
			return
		}
	}
}

// ratingInStars returns rating as ⭐
func ratingInStars(rating int) string {
	return strings.Repeat("⭐", rating)
}
