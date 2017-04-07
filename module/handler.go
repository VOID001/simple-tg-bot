package module

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/VOID001/simple-tg-bot/command"
	session "github.com/VOID001/simple-tg-bot/module/session"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type CommandHandler struct {
	cmdMap map[string]command.Command
	Bot    *tg.BotAPI
}

func (h *CommandHandler) Register(cmdStr string, cmdInstance command.Command) (err error) {
	if h.cmdMap == nil {
		log.Debugf("Initialized Command Table")
		h.cmdMap = make(map[string]command.Command)
	}
	if _, ok := h.cmdMap[cmdStr]; ok == true {
		err = errors.New(fmt.Sprintf("CommandHandler.Register: Command '%s' already registered as another instance", cmdStr))
		return
	}
	h.cmdMap[cmdStr] = cmdInstance
	log.Infof("Register Command %s complete", cmdStr)
	return
}

func (h *CommandHandler) Unregister(cmdStr string) (err error) {
	return
}

func (h *CommandHandler) Command(cmdStr string) (cmdObj command.Command) {
	log.Debugf("CommandHandler.Command cmdStr = %s", cmdStr)
	cmdObj = nil
	if _, ok := h.cmdMap[cmdStr]; ok == true {
		cmdObj = h.cmdMap[cmdStr]
	}
	return
}

func (h *CommandHandler) handle(update tg.Update) (replyMsg tg.Chattable, sess *session.Session, doSend bool, err error) {
	log.Debugf("CommandHandler.handle")
	sess = new(session.Session)
	m := update.Message
	userID := 0
	groupID := int64(0)
	if m == nil {
		m = update.CallbackQuery.Message
	}
	cquery := update.CallbackQuery
	doSend = true // Only turn off when no data should send out

	// If no other errors are set, it will return this default error when needed
	internalErrorReply := tg.NewMessage(m.Chat.ID, "")
	internalErrorReply.Text = "Oops! Something wrong happend, please try again"
	internalErrorReply.ReplyToMessageID = m.MessageID
	replyMsg = internalErrorReply

	// We need to get the realUserID & groupID here
	// If forwarded from A place to B, set to B
	if update.Message != nil {
		userID = m.From.ID
		groupID = m.Chat.ID
	}
	if update.CallbackQuery != nil {
		userID = cquery.From.ID
		groupID = m.Chat.ID
	}

	// Excatly only ONE of the above will be triggered
	if update.Message != nil {
		log.Debugf("Message Recv, Data = %+v", m)
		// userID := m.From.ID
		// groupID := m.Chat.ID
		err = sess.Get(userID, groupID)
		if err != nil {
			err = errors.Wrap(err, "CommandHandler.handle error")
			return
		}
		//process
		replyMsg = tg.NewMessage(m.Chat.ID, "")
		// If message is a command then process it
		if m.IsCommand() {
			cmd := h.Command(m.Command())

			// Command do not exist, don't reply anything
			if cmd == nil {
				doSend = false
				err = errors.New(fmt.Sprintf("CommandHandler.handle error: unknown command %s", m.Command()))
				return
			}
			// Command not null, then we will start here =w=

			// First, no arguments given, show Dialog
			if m.CommandArguments() == "" || m.CommandArguments == nil {
				log.Debugf("Show dialog")
				sess.Cmd = m.Command()
				sess.State = 0
				cmd.Init(sess)
				replyMsg, err = cmd.ParseMessage(m)
				if err != nil {
					err = errors.Wrap(err, "CommandHandler.handle error")
					return
				}
				replyMsg, sess.State, err = cmd.Dialog(0, sess.Cmd) // 0 means the start state
				if err != nil {
					err = errors.Wrap(err, "CommandHandler.handle error")
					return
				}
				return
			}
			// If has arguments, but arguments are wrong, Show usage
			cmd.Init(sess)
			replyMsg, err = cmd.ParseMessage(m)
			if err != nil {
				err = errors.Wrap(err, "CommandHandler.handle error")
				return
			}
			// If arguments are correct, then Run it
			replyMsg, err = cmd.Run()
			if err != nil {
				err = errors.Wrap(err, "CommandHandler.handle error")
				return
			}
		}
		// user input something, but not a command
		if !m.IsCommand() {
			// If m is not command, then user must be in a state, with command, not zero
			if sess.Cmd == "" {
				doSend = false
				return
			}
			cmd := h.Command(sess.Cmd)
			if cmd == nil {
				err = errors.New("CommandHandler.handler error: session data contains unexpected cmd, possible corrupt")
				doSend = false
				return
			}
			cmd.Init(sess)
			replyMsg, err = cmd.ParseMessage(m)
			if err != nil {
				err = errors.Wrap(err, "CommandHandler.handle error")
				doSend = false
				return
			}
			// Command Not nil, Check if is final
			if cmd.IsFinalState(sess.State) {
				_, _, err = cmd.Dialog(sess.State, m.Text)
				if err != nil {
					err = errors.Wrap(err, "commandHandler.handle error")
					return
				}

				replyMsg, err = cmd.Run()
				if err != nil {
					err = errors.Wrap(err, "commandHandler.handle error")
					return
				}
				sess.End()
				return
			}
			// Not Final
			replyMsg, sess.State, err = cmd.Dialog(sess.State, m.Text)
			if err != nil {
				err = errors.Wrap(err, "commandHandler.handle error")
				return
			}
		}
		return // We can skip code below safely
	}
	if cquery != nil {
		defer func() {
			acq := tg.CallbackConfig{}
			acq.CallbackQueryID = cquery.ID
			h.Bot.AnswerCallbackQuery(acq)
		}()
		log.Debugf("Call back Query Recv, Data = %s", cquery.Data)
		m := cquery.Message
		log.Debugf("Message from CallbackData get = %+v", m)
		// userID := cquery.From.ID
		// groupID := m.Chat.ID
		sess.Get(userID, groupID)

		// If not the user who initiate the dialog, abort
		// If the dialog expired, abort
		if sess.CurrentDialogID != m.MessageID {
			doSend = false
			acq := tg.CallbackConfig{}
			acq.CallbackQueryID = cquery.ID
			acq.Text = "You cannot operate on dialogs opened by other people, or dialog expired"
			h.Bot.AnswerCallbackQuery(acq)
			return
		}

		// Get data in callback
		cmd := h.Command(sess.Cmd)
		if cmd == nil {
			err = errors.New("CommandHandler.handler error: session data contains unexpected cmd, possible corrupt")
			doSend = false
			return
		}
		cmd.Init(sess)
		replyMsg, err = cmd.ParseMessage(m)
		if err != nil {
			err = errors.Wrap(err, "CommandHandler.handle error")
			return
		}
		// Is the last step =w=
		if cmd.IsFinalState(sess.State) {
			_, _, err = cmd.Dialog(sess.State, cquery.Data)
			if err != nil {
				err = errors.Wrap(err, "commandHandler.handle error")
				return
			}
			replyMsg, err = cmd.Run()
			if err != nil {
				err = errors.Wrap(err, "commandHandler.handle error")
				return
			}
			sess.End()
			return
		}

		replyMsg, sess.State, err = cmd.Dialog(sess.State, cquery.Data)
		if err != nil {
			err = errors.Wrap(err, "CommandHandler.Run error")
			log.Error(err)
			return
		}
		// Callback Query should be answered
		//acq := tg.CallbackConfig{}
		//acq.CallbackQueryID = cquery.ID
		//h.Bot.AnswerCallbackQuery(acq)
		return
	}
	return
}

func (h *CommandHandler) Run() {
	log.Infof("CommandHandler Running...")

	//get the message in
	u := tg.UpdateConfig{}
	u.Timeout = 60
	updates, err := h.Bot.GetUpdatesChan(u)
	if err != nil {
		err = errors.Wrap(err, "CommandHandler.Run error")
		log.Fatal(err)
	}

	for update := range updates {
		//u.Offset++
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		reply, sess, ok, err := h.handle(update)
		if err != nil {
			err = errors.Wrap(err, "CommandHandler.Run")
			log.Error(err)
		}
		// Only send commands when need to
		if ok && reply != nil {
			log.Debugf("Reply to Send %+v", reply)
			m, err := h.Bot.Send(reply)
			if err != nil {
				err = errors.Wrap(err, "CommandHandler.Run")
				log.Error(err)
				continue
			}
			sess.CurrentDialogID = m.MessageID
			err = sess.Put()
			if err != nil {
				err = errors.Wrap(err, "CommandHandler.Run error")
				continue
			}
		}
	}
	return
}
