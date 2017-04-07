package command

import (
	session "github.com/VOID001/simple-tg-bot/module/session"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Command interface {
	Dialog(state int, cdata string) (msg tg.Chattable, nextState int, err error)
	Info() (msg tg.Chattable, err error)
	Run() (msg tg.Chattable, err error)
	ParseMessage(msg *tg.Message) (errMsg tg.Chattable, err error)
	IsFinalState(state int) (ok bool)
	SetSession(sess *session.Session)
}
