package command

// Demo to show Dialog Basic usage

import (
	"fmt"

	session "github.com/VOID001/simple-tg-bot/module/session"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// You do not need to contain any export member here
type Ping struct {
	times int
	msg   *tg.Message
	sess  session.Session // Read Only Data, you should not modify session
}

func (p *Ping) SetSession(sess *session.Session) {
	p.sess = *sess
}

func (p *Ping) Dialog(state int, data string) (reply tg.Chattable, nextState int, err error) {
	if data != "ping" && data != "pong" {
		return
	}

	// Let's get the user session
	sess := new(session.Session)
	sess.Get(p.msg.From.ID, p.msg.Chat.ID)

	nextState = state + 1 // No next state, it will pause here
	btnRow := tg.NewInlineKeyboardRow()
	btn := tg.InlineKeyboardButton{}
	dat := "ping"
	btn.CallbackData = &dat
	btn.Text = "Click To Ping!"
	btnRow = append(btnRow, btn)
	btnObj := tg.NewInlineKeyboardMarkup(btnRow)

	if state == 0 {
		repl := tg.NewMessage(p.msg.Chat.ID, "")
		repl.Text = "You have Pinged 1 Time!"
		repl.Text += fmt.Sprintf("Session Data\n%+v\n", p.sess)
		repl.ReplyToMessageID = p.msg.MessageID
		// Add InlineButton
		repl.ReplyMarkup = &btnObj
		reply = repl
		return
	}
	repl := tg.NewEditMessageText(p.msg.Chat.ID, p.msg.MessageID, "")
	repl.Text = fmt.Sprintf("You have Pinged %d Times!", state+1)
	repl.Text += fmt.Sprintf("Session Data\n%+v\n", p.sess)
	repl.ReplyMarkup = &btnObj
	reply = repl
	return
}

func (p *Ping) Info() (reply tg.Chattable, err error) {
	repl := tg.NewMessage(p.msg.Chat.ID, "usage: Just Ping as you wish")
	reply = repl
	return
}

func (p *Ping) Run() (reply tg.Chattable, err error) {
	return
}

func (p *Ping) ParseMessage(msg *tg.Message) (errMsg tg.Chattable, err error) {
	p.msg = msg
	if msg.CommandArguments() != "" {
		errMsg, err = p.Info()
		err = errors.New("Ping.ParseMessage: argument should be null")
		return
	}
	return
}

func (p *Ping) IsFinalState(state int) (ok bool) {
	// This will never end
	ok = false
	return
}
