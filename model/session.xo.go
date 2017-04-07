// Package model contains the types for schema 'todolist'.
package model

// GENERATED BY XO. DO NOT EDIT.

import (
	"errors"
)

// Session represents a row from 'todolist.sessions'.
type Session struct {
	ID      int    `json:"id"`       // id
	UserID  int    `json:"user_id"`  // user_id
	GroupID int64  `json:"group_id"` // group_id
	Session string `json:"session"`  // session

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Session exists in the database.
func (s *Session) Exists() bool {
	return s._exists
}

// Deleted provides information if the Session has been deleted from the database.
func (s *Session) Deleted() bool {
	return s._deleted
}

// Insert inserts the Session to the database.
func (s *Session) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if s._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO todolist.sessions (` +
		`id, user_id, group_id, session` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, s.ID, s.UserID, s.GroupID, s.Session)
	_, err = db.Exec(sqlstr, s.ID, s.UserID, s.GroupID, s.Session)
	if err != nil {
		return err
	}

	// set existence
	s._exists = true

	return nil
}

// Update updates the Session in the database.
func (s *Session) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if s._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE todolist.sessions SET ` +
		`user_id = ?, group_id = ?, session = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, s.UserID, s.GroupID, s.Session, s.ID)
	_, err = db.Exec(sqlstr, s.UserID, s.GroupID, s.Session, s.ID)
	return err
}

// Save saves the Session to the database.
func (s *Session) Save(db XODB) error {
	if s.Exists() {
		return s.Update(db)
	}

	return s.Insert(db)
}

// Delete deletes the Session from the database.
func (s *Session) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return nil
	}

	// if deleted, bail
	if s._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM todolist.sessions WHERE id = ?`

	// run query
	XOLog(sqlstr, s.ID)
	_, err = db.Exec(sqlstr, s.ID)
	if err != nil {
		return err
	}

	// set deleted
	s._deleted = true

	return nil
}

// SessionByID retrieves a row from 'todolist.sessions' as a Session.
//
// Generated from index 'sessions_id_pkey'.
func SessionByID(db XODB, id int) (*Session, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, group_id, session ` +
		`FROM todolist.sessions ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	s := Session{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&s.ID, &s.UserID, &s.GroupID, &s.Session)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Below are customized code added by VOID001

func SessionByUserIDAndGroupID(db XODB, userID int, groupID int64) (*Session, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, group_id, session ` +
		`FROM todolist.sessions ` +
		`WHERE user_id = ? AND group_id = ?`

	// run query
	XOLog(sqlstr, userID, groupID)
	s := Session{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, userID, groupID).Scan(&s.ID, &s.UserID, &s.GroupID, &s.Session)
	if err != nil {
		return nil, err
	}

	return &s, nil
}