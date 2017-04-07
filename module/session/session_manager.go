package module

import (
	"database/sql"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/VOID001/simple-tg-bot/model"
	"github.com/pkg/errors"
)

type Session struct {
	model           *model.Session
	userID          int
	groupID         int64
	Cmd             string      `json:"cmd"`
	State           int         `json:"state"`
	LastData        string      `json:"last_data"`
	CurrentDialogID int         `json:"current_dialog_id"`
	Data            interface{} `json"data"`
}

func (s *Session) Get(userID int, groupID int64) (err error) {
	log.Infof("Call Session.Get userID=%d groupID=%d", userID, groupID)
	ss, er := model.SessionByUserIDAndGroupID(model.DB, userID, groupID)
	if er == sql.ErrNoRows {
		// The session is not created yet, create it
		ss = new(model.Session)
		ss.GroupID = groupID
		ss.UserID = userID
		s.model = ss
		s.Cmd = ""
		s.State = 0
		s.LastData = ""
		dat, er := json.Marshal(s)
		ss.Session = fmt.Sprintf("%s", dat)
		if er != nil {
			err = errors.Wrap(er, "Session.Get")
			return
		}
		err = ss.Insert(model.DB)
		if err != nil {
			err = errors.Wrap(err, "Session.Get")
			return
		}
		log.Infof("Sucessfully create session for userID %d", userID)
		return
	}
	if er != nil {
		err = errors.Wrap(er, "Session.Get")
		return
	}
	err = json.Unmarshal([]byte(ss.Session), &s)
	if err != nil {
		err = errors.Wrap(err, "Session.Get")
		return
	}
	s.model = ss
	s.userID = ss.UserID
	s.groupID = ss.GroupID
	return
}

func (s *Session) Put() (err error) {
	dat, er := json.Marshal(s)
	if er != nil {
		err = errors.Wrap(er, "Session.Put error")
		return
	}
	s.model.Session = fmt.Sprintf("%s", dat)
	err = s.model.Update(model.DB)
	if err != nil {
		err = errors.Wrap(err, "Session.Put error")
		return
	}
	return
}

func (s *Session) End() (err error) {
	s.Cmd = ""
	s.State = 0
	err = s.Put()
	if err != nil {
		err = errors.Wrap(err, "Session.Put error")
		return
	}
	return
}
