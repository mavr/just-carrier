package staff_msg

import (
	"golang.org/x/text/language"
)

var msgWelcome = translate{
	m: map[language.Tag]string{
		language.English: `*Welcome to anonymous bot mailer!*
		
To send anonymous message for user use
` + "```" + `
/send @username text of the message you want to send
` + "```",

		language.Russian: `*Добро пожаловать в бот для отправки анонимных сообщений!*
		
Для того что б послать анонимное сообщение человеку используйте следующий формат:
` + "```" + `
/send @username text of the message you want to send
` + "```",
	},
}

var msgInvalidSendCommand = translate{
	m: map[language.Tag]string{
		language.English: `Invalid format of the /send command`,
		language.Russian: `Неверный формат команды /send`,
	},
}
