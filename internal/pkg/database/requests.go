package database

import (
	l "telegram-bot-feedback/internal/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// AddEmployeeByID creates/updates User by Telegram ID with field IsEmployee = true
func AddEmployeeByID(db *gorm.DB, id int) error {
	user := User{}
	db.Where("chat_id = ?", id).First(&user)
	user.ChatID = id
	user.IsEmployee = true
	return l.Err(db.Save(&user).Error)
}

// AddEmployeeByNickname creates/updates User by Telegram Nickname with field IsEmployee = true
func AddEmployeeByNickname(db *gorm.DB, nick string) error {
	user := User{}
	db.Where("nickname = ?", nick).First(&user)
	user.Nickname = nick
	user.IsEmployee = true
	return l.Err(db.Save(&user).Error)
}

// RemoveEmployeeByID creates/updates User by Telegram ID with field IsEmployee = false
func RemoveEmployeeByID(db *gorm.DB, id int) error {
	user := User{}
	db.Where("chat_id = ?", id).First(&user)
	user.ChatID = id
	user.IsEmployee = false
	return l.Err(db.Save(&user).Error)
}

// RemoveEmployeeByNickname creates/updates User by Telegram Nickname with field IsEmployee = false
func RemoveEmployeeByNickname(db *gorm.DB, nick string) error {
	user := User{}
	db.Where("nickname = ?", nick).First(&user)
	user.Nickname = nick
	user.IsEmployee = false
	return l.Err(db.Save(&user).Error)
}

// AddUser creates/updates User
func AddUser(chatId int, nick string, state int, db *gorm.DB) (*User, error) {
	user := User{}
	db.Where("chat_id = ? OR nickname = ?", chatId, nick).First(&user)
	user.Nickname = nick
	user.ChatID = chatId
	user.State = state
	user.IsReceiver = false
	err := db.Save(&user).Error
	return &user, l.Err(err)
}

// AddQuestion creates Question from User
func AddQuestion(header string, user *User, db *gorm.DB) (*Question, error) {
	question := Question{}
	question.UserID = int(user.ID)
	question.Header = header
	err := db.Save(&question).Error
	return &question, l.Err(err)
}

// AddCorrespondence creates Correspondence from User
func AddCorrespondence(user *User, messageId int, db *gorm.DB) (*QuestionCorrespondence, error) {
	question := &Question{}
	if user.IsEmployee {
		question = GetOpenQuestionByAnswerer(user, db)
	} else {
		question = GetOpenQuestionByUser(user, db)
	}
	if question == nil {
		return nil, nil
	}
	corr := QuestionCorrespondence{
		QuestionID: int(question.ID),
		MessageID:  messageId,
		User:       *user,
		IsEmployee: false,
	}
	err := db.Save(&corr).Error
	return &corr, l.Err(err)
}

// GetEmployees returns the Users with field IsEmployee = true
func GetEmployees(db *gorm.DB) []User {
	users := []User{}
	err := db.Where("is_employee = ?", true).Find(&users).Error
	if err != nil || len(users) == 0 {
		return nil
	}
	return users
}

// GetReceivers returns the Users with fields IsEmployee = true and IsReceiver = true
func GetReceivers(db *gorm.DB) []User {
	users := []User{}

	err := db.Where("is_employee = ? AND is_receiver = ?", true, true).Where("NOT EXISTS (?)", db.Table("questions").Select("id").Where("answerer_id = users.id")).Find(&users).Error
	if err != nil || len(users) == 0 {
		return nil
	}
	return users
}

// GetUserByChatID returns User by Telegram ID (or private Chat ID)
func GetUserByChatID(chatId int, db *gorm.DB) *User {
	user := User{}
	err := db.Where("chat_id = ?", chatId).First(&user).Error
	if err != nil || user.ID == 0 {
		return nil
	}
	return &user
}

// GetEmptyReview returns Review from User with empty Text
func GetEmptyReview(user *User, db *gorm.DB) *Review {
	review := Review{}
	err := db.Preload("User").Where("user_id = ? AND text = ?", user.ID, "").First(&review).Error
	if err != nil || review.ID == 0 {
		return nil
	}
	return &review
}

// GetReviewsInRange returns Reviews between two dates
func GetReviewsInRange(fDate time.Time, sDate time.Time, db *gorm.DB) []Review {
	reviews := []Review{}
	err := db.Order("id asc").Where("created_at BETWEEN ? AND ?", sDate, fDate).Find(&reviews).Error
	if err != nil || len(reviews) == 0 {
		return nil
	}
	return reviews
}

// GetCountReviewsByRating returns the number of Reviews with each rating
func GetCountReviewsByRating(db *gorm.DB) [5]int64 {
	number := [5]int64{}
	for i := 0; i < 5; i++ {
		var c int64
		db.Model(&Review{}).Where("rating = ?", i+1).Count(&c)
		number[i] = c
	}
	return number
}

// GetQuestionById returns Question by ID with preloading User and Answerer
func GetQuestionById(id int, db *gorm.DB) *Question {
	question := Question{}
	err := db.Preload("User").Preload("Answerer").First(&question, id).Error
	if err != nil || question.ID == 0 {
		return nil
	}
	return &question
}

// GetOpenQuestionByUser returns Question by User with preloading User and Answerer
func GetOpenQuestionByUser(user *User, db *gorm.DB) *Question {
	question := Question{}
	err := db.Preload("User").Preload("Answerer").Where("user_id = ? AND is_closed = ?", user.ID, false).First(&question).Error
	if err != nil || question.ID == 0 {
		return nil
	}
	return &question
}

// GetOpenQuestionByAnswerer returns Question by Answerer with preloading User and Answerer
func GetOpenQuestionByAnswerer(user *User, db *gorm.DB) *Question {
	question := Question{}
	err := db.Preload("User").Preload("Answerer").Where("answerer_id = ? AND is_closed = ?", user.ID, false).First(&question).Error
	if err != nil || question.ID == 0 {
		return nil
	}
	return &question
}

// GetNewQuestionById returns open Question without answer and Answerer by ID
func GetNewQuestionById(id int, db *gorm.DB) *Question {
	question := Question{}
	err := db.Where("id = ? AND (answerer_id IS NULL OR answerer_id = 0) AND have_answer = ? AND is_closed = ?", id, false, false).First(&question).Error
	if err != nil || question.ID == 0 {
		return nil
	}
	return &question
}

// GetNewQuestions returns open Questions without answer and Answerer
func GetNewQuestions(db *gorm.DB) []Question {
	questions := []Question{}
	err := db.Order("id asc").Find(&questions, "(answerer_id IS NULL OR answerer_id = 0) AND have_answer = ? AND is_closed = ?", false, false).Error
	if err != nil || len(questions) == 0 {
		return nil
	}
	return questions
}

// GetCorrespondenceByQuestion returns Correspondence by Question with preloading User
func GetCorrespondenceByQuestion(questions *Question, db *gorm.DB) []QuestionCorrespondence {
	corr := []QuestionCorrespondence{}
	err := db.Preload("User").Where("question_id = ?", questions.ID).Order("id asc").Find(&corr).Error
	if err != nil || len(corr) == 0 {
		return nil
	}
	return corr
}

// ChangeUserState change User "State"
func ChangeUserState(state int, user *User, db *gorm.DB) error {
	user.State = state
	err := db.Save(user).Error
	return l.Err(err)
}

// ChangeUserIsReceiver change User "IsReceiver"
func ChangeUserIsReceiver(isReceiver bool, user *User, db *gorm.DB) error {
	user.IsReceiver = isReceiver
	err := db.Save(user).Error
	return l.Err(err)
}

// ChangeTextReviewByUser change Review "Text" (by User)
func ChangeTextReviewByUser(text string, user *User, db *gorm.DB) error {
	review := GetEmptyReview(user, db)
	if review == nil {
		return nil
	}
	review.Text = text
	return l.Err(db.Save(review).Error)
}

// ChangeQuestionHaveAnswer change Question "HaveAnswer"
func ChangeQuestionHaveAnswer(state bool, question *Question, db *gorm.DB) error {
	question.HaveAnswer = state
	err := db.Save(question).Error
	return l.Err(err)
}

// ChangeQuestionAnswerer change Question "Answerer"
func ChangeQuestionAnswerer(answererID int, question *Question, db *gorm.DB) error {
	question.AnswererID = answererID
	err := db.Save(question).Error
	return l.Err(err)
}

// ChangeQuestionIsClosed change Question "IsClosed"
func ChangeQuestionIsClosed(closed bool, question *Question, db *gorm.DB) error {
	question.IsClosed = closed
	err := db.Save(question).Error
	return l.Err(err)
}
