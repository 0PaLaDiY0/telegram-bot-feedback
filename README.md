﻿# Telegram Bot for feedback

This bot implements a simple functionality for receiving feedback.

### How to start:

1. Create a telegram bot (You need a token) [Guide](https://core.telegram.org/bots/tutorial#getting-ready)
2. Download the repository and compile OR download [here](https://github.com/0PaLaDiY0/telegram-bot-feedback/releases)
3. Run the file
4. Specify host for local server or "-" for standard
5. Enter token

*In the folder of the executable file, the database and error folders, as well as the configuration file, will be automatically created.*

### Console

Here are the available commands:
```
abi <id> - adds employee by user ID
abn <nickname> - adds an employee by user Nickname
rbi <id> - removes an employee by user ID
rbn <nickname> - removes an employee by user Nickname
ge - displays a list of employees
close - closes the program
```

### User functionality
The user can leave reviews with or without comments:

![](https://i.ibb.co/rmm6TCY/UserRev.gif)

---
The user can ask a question:

![](https://i.ibb.co/23rrBtN/UserQue.gif)

*An employee of the company answers the question, and the answer comes to the user from the bot*

### Employee functionality

An employee can toggle receiving questions:

![](https://i.ibb.co/JjM0KZ6/ERec.gif)

*If receiving is enabled, the bot will send new questions in real time.*

---
An employee сan get a list of questions and take a question. 

When you take a question, the message history is loaded:

![](https://i.ibb.co/bNXQsSr/ETake.gif)

*In the example, when sending messages, the bot responds with the same message. This happens because user and employee are one person. In a real case, messages will come to the user who asked the question.*

*Only one employee can take a question. Also, if the question has been answered, it will disappear from the list.*

---
An employee can find a question by number. Message history will be loaded:

![](https://i.ibb.co/F0Z96hH/EQue.gif)

---
An employee can view reviews for a period or for all time.:

![](https://i.ibb.co/zPPTJHB/ERev.gif)
