package console

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"telegram-bot-feedback/internal/pkg/database"
	l "telegram-bot-feedback/internal/pkg/logger"

	"gorm.io/gorm"
)

// Run starts reading commands from the console
func Run(cancel context.CancelFunc, db *gorm.DB) {
	for {
		in := bufio.NewScanner(os.Stdin)
		in.Scan()
		command := strings.Split(in.Text(), " ")
		switch command[0] {
		case "":
		case "help":
			fmt.Println("Here are the available commands:")
			fmt.Println("abi <id> - adds employee by user ID")
			fmt.Println("abn <nickname> - adds an employee by user Nickname")
			fmt.Println("rbi <id> - removes an employee by user ID")
			fmt.Println("rbn <nickname> - removes an employee by user Nickname")
			fmt.Println("ge - displays a list of employees")
			fmt.Println("close - closes the program")
		case "abi":
			if len(command) > 1 {
				id, err := strconv.Atoi(command[1])
				if err != nil {
					fmt.Println("Wrong format")
					break
				}
				err = database.AddEmployeeByID(db, id)
				if err != nil {
					l.Error(err)
					break
				}
				fmt.Println("Employee added")
				break
			}
			fmt.Println("Enter value")
		case "abn":
			if len(command) > 1 {
				nick := command[1]
				err := database.AddEmployeeByNickname(db, nick)
				if err != nil {
					l.Error(err)
					break
				}
				fmt.Println("Employee added")
				break
			}
			fmt.Println("Enter value")
		case "rbi":
			if len(command) > 1 {
				id, err := strconv.Atoi(command[1])
				if err != nil {
					fmt.Println("Wrong format")
					break
				}
				err = database.RemoveEmployeeByID(db, id)
				if err != nil {
					l.Error(err)
					break
				}
				fmt.Println("Employee removed")
				break
			}
			fmt.Println("Enter value")
		case "rbn":
			if len(command) > 1 {
				nick := command[1]
				err := database.RemoveEmployeeByNickname(db, nick)
				if err != nil {
					l.Error(err)
					break
				}
				fmt.Println("Employee removed")
				break
			}
			fmt.Println("Enter value")
		case "ge":
			users := database.GetEmployees(db)
			for _, user := range users {
				fmt.Printf("UserID: %d Nickname: %s\n", user.ChatID, user.Nickname)
				fmt.Println("(empty fields are filled when the employee uses the bot)")
			}
		case "close":
			cancel()
			return
		default:
			fmt.Println("Unknown command, use \"help\"")
		}
	}
}
