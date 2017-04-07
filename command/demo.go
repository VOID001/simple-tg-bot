package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	log "github.com/Sirupsen/logrus"
	session "github.com/VOID001/simple-tg-bot/module/session"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Demo Command use to show dialog switch and command execute

type Demo struct {
	adj   string
	noun  string
	verb  string
	oadj  string
	onoun string
	msg   *tg.Message
	sess  *session.Session
}

// Use to store in user session
type DemoArgs struct {
	Adj   string
	Noun  string
	Verb  string
	OAdj  string
	ONoun string
}

type btn struct {
	Text string
	Data string
}

var stateMsg = []string{
	"",
	"欢迎使用非常简易的句子生成器\n\n请输入/选择一个名词",
	"请输入/选择一个形容词, 点击/back可以返回上一级重新填写",
	"请输入/选择一个动词",
	"请输入/选择另一个名词",
	"请输入/选择另一个形容词",
}

var btnMatrix = [][][]btn{
	[][]btn{},
	[][]btn{
		[]btn{btn{Text: "夏娜", Data: "夏娜"}, btn{Text: "彩虹猫", Data: "彩虹猫"}},
		[]btn{btn{Text: "罗技鼠标", Data: "罗技鼠标"}, btn{Text: "女装少年", Data: "女装少年"}},
	},
	[][]btn{
		[]btn{btn{Text: "可爱", Data: "可爱"}, btn{Text: "傲娇", Data: "傲娇"}},
		[]btn{btn{Text: "色情", Data: "色情"}},
		[]btn{btn{Text: "返回", Data: "/back"}},
	},
	[][]btn{
		[]btn{btn{Text: "肛", Data: "肛"}, btn{Text: "调戏", Data: "调戏"}},
		[]btn{btn{Text: "返回", Data: "/back"}},
	},
	[][]btn{
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "洋葱", Data: "洋葱"}, btn{Text: "洋葱", Data: "洋葱"}},
		[]btn{btn{Text: "返回", Data: "/back"}},
	},
	[][]btn{
		[]btn{btn{Text: "可爱", Data: "可爱"}, btn{Text: "傲娇", Data: "傲娇"}},
		[]btn{btn{Text: "色情", Data: "色情"}},
		[]btn{btn{Text: "返回", Data: "/back"}},
	},
}

var stateBtn []tg.InlineKeyboardMarkup

func (d *Demo) Init(sess *session.Session, args ...interface{}) {
	stateBtn = make([]tg.InlineKeyboardMarkup, 0)
	markup := tg.NewInlineKeyboardMarkup()
	btnRow := tg.NewInlineKeyboardRow()
	btn := tg.NewInlineKeyboardButtonData("", "")
	for _, v := range btnMatrix {
		markup = tg.NewInlineKeyboardMarkup()
		for _, vv := range v {
			btnRow = tg.NewInlineKeyboardRow()
			for _, vvv := range vv {
				btn = tg.NewInlineKeyboardButtonData(vvv.Text, vvv.Data)
				btnRow = append(btnRow, btn)
			}
			markup.InlineKeyboard = append(markup.InlineKeyboard, btnRow)
		}
		stateBtn = append(stateBtn, markup)
	}
	d.sess = sess
}

func (d *Demo) Dialog(state int, data string) (reply tg.Chattable, nextState int, err error) {
	log.Infof("State = %d Data = %s\n", state, data)
	if data == "/back" {
		log.Infof("Go back: State = %d Data = %s\n", nextState, data)
		nextState = state - 1
		repl := tg.NewEditMessageText(d.msg.Chat.ID, d.sess.CurrentDialogID, stateMsg[state-1])
		repl.ReplyMarkup = &stateBtn[state-1]
		reply = repl
		return
	}
	args := DemoArgs{}
	if state == 0 {
		repl := tg.NewMessage(d.msg.Chat.ID, stateMsg[state+1])
		repl.ReplyMarkup = stateBtn[state+1]
		repl.ReplyToMessageID = d.msg.MessageID
		reply = repl
		nextState = state + 1
		return
	}
	if !d.IsFinalState(state) {
		if d.sess.Data != "" {
			err = json.Unmarshal([]byte(d.sess.Data), &args)
			if err != nil {
				err = errors.Wrap(err, "Demo.Dialog")
				return
			}
		}
		// Diffrent state set different val
		if state == 1 {
			args.Noun = data
		}
		if state == 2 {
			args.Adj = data
		}
		if state == 3 {
			args.Verb = data
		}
		if state == 4 {
			args.ONoun = data
		}

		dat, er := json.Marshal(args)
		if er != nil {
			err = errors.Wrap(er, "Demo.Dialog")
			return
		}
		d.sess.Data = fmt.Sprintf("%s", dat)
		repl := tg.NewEditMessageText(d.msg.Chat.ID, d.sess.CurrentDialogID, stateMsg[state+1])
		repl.ReplyMarkup = &stateBtn[state+1]
		reply = repl
		nextState = state + 1
	}
	// Final state does not return dialog data
	if d.IsFinalState(state) {
		err = json.Unmarshal([]byte(d.sess.Data), &args)
		if err != nil {
			err = errors.Wrap(err, "Demo.Dialog")
			return
		}
		args.OAdj = data
		dat, er := json.Marshal(args)
		if er != nil {
			err = errors.Wrap(er, "Demo.Dialog")
			return
		}
		d.sess.Data = fmt.Sprintf("%s", dat)
		return
	}
	return
}

func (d *Demo) Info() (reply tg.Chattable, err error) {
	repl := tg.NewMessage(d.msg.Chat.ID, "用法: /demo [名词] [形容词]")
	repl.ReplyToMessageID = d.msg.MessageID
	reply = repl
	return
}

func (d *Demo) Run() (reply tg.Chattable, err error) {
	errRepl := tg.NewMessage(d.msg.Chat.ID, "哎呀，出错了呢, 请再试一次哦")
	errRepl.ReplyToMessageID = d.msg.MessageID
	args := DemoArgs{}
	if d.sess.Data != "" {
		err = json.Unmarshal([]byte(d.sess.Data), &args)
		if err != nil {
			err = errors.Wrap(err, "Demo.Run")
			reply = errRepl
			return
		}
		d.noun = args.Noun
		d.adj = args.Adj
		d.verb = args.Verb
		d.onoun = args.ONoun
		d.oadj = args.OAdj
	}
	reply = tg.NewMessage(d.msg.Chat.ID, fmt.Sprintf("智障造句的结果为: %s的%s%s了%s的%s", d.adj, d.noun, d.verb, d.oadj, d.onoun))
	return
}

func (d *Demo) ParseMessage(msg *tg.Message) (errMsg tg.Chattable, err error) {
	d.msg = msg
	errRepl := tg.NewMessage(d.msg.Chat.ID, "请给五个参数哦，不要多也不要少 OAO")
	errRepl.ReplyToMessageID = d.msg.MessageID
	args := msg.CommandArguments()
	if args == "" {
		return
	}
	argList := strings.Split(args, " ")
	if len(argList) != 5 {
		err = errors.New("Wrong argument given")
		errMsg = errRepl
		return
	}
	d.adj = argList[1]
	d.noun = argList[0]
	d.verb = argList[2]
	d.onoun = argList[3]
	d.oadj = argList[4]
	return
}

func (d *Demo) IsFinalState(state int) (ok bool) {
	if state == 5 {
		return true
	}
	return false
}
