package bot

import (
	"strconv"
	"strings"
	"telegram-bot-feedback/internal/pkg/database"
	l "telegram-bot-feedback/internal/pkg/logger"
	tg "telegram-bot-feedback/pkg/telegram-bot-api"
)

// User states
const (
	SNew int = iota + 1
	SMain
	SReview
	SReviewText
	SQuestion
	SQuestionDiscussion
	SSwitchReceiver
	SSearchQuestion
)

// Callback data types
const (
	CBQuestion int = iota + 1
)

// Date intervals
const (
	RDay int = iota + 1
	RWeek
	RMonth
	RAll
)

// parseUpdate parse bot Update
func parseUpdate(update *tg.Update, app *App) (err error) {
	if update.Message != nil {
		err = parseMessage(update.Message, app)
		if err != nil {
			l.Err(err)
		}
	}
	if update.CallbackQuery != nil {
		err = parseCallback(update.CallbackQuery, app)
		if err != nil {
			l.Err(err)
		}
	}
	if err == nil {
		app.Conf.Set("offset", update.UpdateID+1)
		err = app.Conf.WriteConfig()
	}
	return l.Err(err)
}

// parseMessage parse Message
func parseMessage(message *tg.Message, app *App) (err error) {
	if isCommand, err := parseCommand(message, app); isCommand {
		return l.Err(err)
	}
	user := database.GetUserByChatID(message.From.ID, app.DB)
	if user == nil {
		return l.Err(l.NewError("User " + strconv.Itoa(int(message.From.ID)) + " is not found"))
	}
	if user.IsEmployee {
		return l.Err(parseMessageEmployee(user, message, app))
	}
	return l.Err(parseMessageUser(user, message, app))
}

// parseMessageUser parse Message from user
func parseMessageUser(user *database.User, message *tg.Message, app *App) (err error) {
	switch user.State {
	case SNew:
		return database.ChangeUserState(SMain, user, app.DB)
	case SMain:
		switch message.Text {
		case "⭐Review":
			err := database.ChangeUserState(SReview, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SMain, user, app.DB)
			}
			return l.Err(err)
		case "❓Question":
			err := database.ChangeUserState(SQuestion, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SMain, user, app.DB)
			}
			return l.Err(err)
		default:
			return nil
		}
	case SReview:
		switch message.Text {
		case "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐", "1", "2", "3", "4", "5":
			err := parseReview(message.Text, user, app)
			if err != nil {
				return l.Err(err)
			}
			err = database.ChangeUserState(SReviewText, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SReview, user, app.DB)
			}
			return l.Err(err)
		default:
			return nil
		}
	case SReviewText:
		switch message.Text {
		case "❌Close":
			err := database.ChangeTextReviewByUser("-", user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SReviewText, user, app.DB)
			}
			return l.Err(err)
		default:
			err := database.ChangeTextReviewByUser(message.Text, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SReviewText, user, app.DB)
			}
			return l.Err(err)
		}
	case SQuestion:
		switch message.Text {
		case "❌Close":
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SQuestion, user, app.DB)
			}
			return l.Err(err)
		default:
			question, err := database.AddQuestion(message.Text, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			questions := []database.Question{*question}
			receivers := database.GetReceivers(app.DB)
			for _, receiver := range receivers {
				sendQuestions(&receiver, app.Bot, questions)
			}
			err = database.ChangeUserState(SQuestionDiscussion, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SQuestion, user, app.DB)
			}
			return l.Err(err)
		}
	case SQuestionDiscussion:
		switch message.Text {
		case "❌Close":
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			question := database.GetOpenQuestionByUser(user, app.DB)
			if question != nil {
				_, err = database.AddCorrespondence(user, message.MessageID, app.DB)
				if err != nil {
					return l.Err(err)
				}
				err = database.ChangeQuestionIsClosed(true, question, app.DB)
				if err != nil {
					return l.Err(err)
				}
				if question.Answerer.ID != 0 {
					err = sendCorrespondenceFromUser(question, message, app.Bot)
					if err != nil {
						return l.Err(err)
					}
				}
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SQuestionDiscussion, user, app.DB)
			}
			return l.Err(err)
		default:
			question := database.GetOpenQuestionByUser(user, app.DB)
			if question == nil {
				return nil
			}
			if question.Answerer.ID != 0 {
				err = sendCorrespondenceFromUser(question, message, app.Bot)
				if err != nil {
					return l.Err(err)
				}
			}
			err = database.ChangeQuestionHaveAnswer(false, question, app.DB)
			if err != nil {
				return l.Err(err)
			}
			_, err = database.AddCorrespondence(user, message.MessageID, app.DB)
			return l.Err(err)
		}
	default:
		return nil
	}
}

// parseMessageUser parse Message from employee
func parseMessageEmployee(user *database.User, message *tg.Message, app *App) (err error) {
	switch user.State {
	case SNew:
		return database.ChangeUserState(SMain, user, app.DB)
	case SMain:
		switch message.Text {
		case "❓Receive questions":
			err = database.ChangeUserState(SSwitchReceiver, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			return l.Err(err)
		case "❓Do not receive questions":
			err = database.ChangeUserState(SSwitchReceiver, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			return l.Err(err)
		case "❓Open questions":
			questions := database.GetNewQuestions(app.DB)
			if len(questions) == 0 {
				err = database.ChangeUserState(SQuestion, user, app.DB)
				if err != nil {
					return l.Err(err)
				}
				return l.Err(responser(user, app))
			}
			sendQuestions(user, app.Bot, questions)
			return l.Err(err)
		case "⭐Reviews":
			err := database.ChangeUserState(SReview, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SMain, user, app.DB)
			}
			return l.Err(err)
		case "❓Find a question":
			err := database.ChangeUserState(SSearchQuestion, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SMain, user, app.DB)
			}
			return l.Err(err)
		default:
			return nil
		}
	case SReview:
		switch message.Text {
		case "📅For a day":
			loadReviews(RDay, user, app)
			return nil
		case "📅For a week":
			loadReviews(RWeek, user, app)
			return nil
		case "📅For a month":
			loadReviews(RMonth, user, app)
			return nil
		case "📅All (no text)":
			loadReviews(RAll, user, app)
			return nil
		case "↩️Back":
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SQuestionDiscussion, user, app.DB)
			}
			return l.Err(err)
		default:
			return nil
		}
	case SQuestionDiscussion:
		switch message.Text {
		case "↩️Back":
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			question := database.GetOpenQuestionByAnswerer(user, app.DB)
			if question != nil {
				err = database.ChangeQuestionAnswerer(0, question, app.DB)
				if err != nil {
					return l.Err(err)
				}
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SQuestionDiscussion, user, app.DB)
			}
			return l.Err(err)
		default:
			question := database.GetOpenQuestionByAnswerer(user, app.DB)
			if question != nil {
				err = sendCorrespondenceFromAnswerer(question, message, app.Bot)
				if err != nil {
					return l.Err(err)
				}
				err = database.ChangeQuestionHaveAnswer(true, question, app.DB)
				if err != nil {
					return l.Err(err)
				}
				_, err = database.AddCorrespondence(user, message.MessageID, app.DB)
				return l.Err(err)
			}
			return nil
		}
	case SSearchQuestion:
		switch message.Text {
		case "↩️Back":
			err = database.ChangeUserState(SMain, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SSearchQuestion, user, app.DB)
			}
			return l.Err(err)
		default:
			loadFullQuestionById(message.Text, user, app)
			return nil
		}
	default:
		return nil
	}
}

// parseCommand parse commands
func parseCommand(message *tg.Message, app *App) (bool, error) {
	switch message.Text {
	case "/start":
		user, err := database.AddUser(message.From.ID, message.From.UserName, SNew, app.DB)
		if err != nil {
			return true, l.Err(err)
		}
		question := database.GetOpenQuestionByUser(user, app.DB)
		if question != nil {
			err = database.ChangeQuestionIsClosed(true, question, app.DB)
			if err != nil {
				return true, l.Err(err)
			}
		}
		err = responserCommand(message.Text, user, app)
		return true, l.Err(err)
	default:
		return false, nil
	}
}

// parseCallback parse CallbackQuery
func parseCallback(callback *tg.CallbackQuery, app *App) error {
	user := database.GetUserByChatID(callback.Message.Chat.ID, app.DB)
	if user == nil {
		return l.Err(l.NewError("User " + strconv.Itoa(int(callback.Message.Chat.ID)) + " is not found"))
	}
	if user.IsEmployee {
		return l.Err(parseCallbackEmployee(user, callback, app))
	}
	return l.Err(parseCallbackUser(user, callback, app))
}

// parseCallbackUser parse CallbackQuery from user
//
// (not used)
func parseCallbackUser(user *database.User, callback *tg.CallbackQuery, app *App) (err error) {
	return nil
}

// parseCallbackUser parse CallbackQuery from employee
func parseCallbackEmployee(user *database.User, callback *tg.CallbackQuery, app *App) (err error) {
	key, data := splitCallbackData(callback)
	switch user.State {
	case SMain:
		switch key {
		case CBQuestion:
			id, err := strconv.Atoi(data)
			if err != nil {
				return l.Err(l.NewError("no id"))
			}
			err = loadCorrespondence(id, user, app)
			if err != nil {
				return l.Err(err)
			}
			err = database.ChangeUserState(SQuestionDiscussion, user, app.DB)
			if err != nil {
				return l.Err(err)
			}
			err = responser(user, app)
			if err != nil {
				database.ChangeUserState(SMain, user, app.DB)
			}
			return l.Err(err)
		default:
			return nil
		}
	default:
		return nil
	}
}

// parseReview parse rating Review
func parseReview(rating string, user *database.User, app *App) error {
	var r int
	switch rating {
	case "⭐", "1":
		r = 1
	case "⭐⭐", "2":
		r = 2
	case "⭐⭐⭐", "3":
		r = 3
	case "⭐⭐⭐⭐", "4":
		r = 4
	case "⭐⭐⭐⭐⭐", "5":
		r = 5
	}
	review := database.Review{User: *user, Rating: r}
	return l.Err(app.DB.Save(&review).Error)
}

// splitCallbackData split data from CallbackQuery
func splitCallbackData(callback *tg.CallbackQuery) (int, string) {
	parts := strings.Split(callback.Data, "-")
	key, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, ""
	}
	if len(parts) > 1 {
		return key, parts[1]
	}
	return key, ""
}
