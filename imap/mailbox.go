package imap

import (
	"errors"
	"time"

	"github.com/emersion/go-imap"
	"github.com/hectane/hectane/db"
)

var ErrUnimplemented = errors.New("not yet implemented")

// mailbox maintains information about a specific folder.
type mailbox struct {
	imap   *IMAP
	folder *db.Folder
}

func (m *mailbox) count(unseen bool) (uint32, error) {
	var (
		count uint32
		c     = db.C.
			Model(&db.Message{}).
			Where("folder_id = ?", m.folder.ID).
			Count(&count)
	)
	if unseen {
		c = c.Where("is_unread = ?", true)
	}
	if err := c.Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Name retrieves the name of the folder.
func (m *mailbox) Name() string {
	return m.folder.Name
}

// Info retrieves information about the folder.
func (m *mailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{
		Attributes: []string{},
		Name:       m.folder.Name,
	}, nil
}

// TODO

// Status retrieves information about the messages in the folder.
func (m *mailbox) Status(items []string) (*imap.MailboxStatus, error) {
	s := imap.NewMailboxStatus(m.folder.Name, items)
	s.Flags = []string{}
	s.PermanentFlags = []string{}
	for _, item := range items {
		switch item {
		case imap.MailboxMessages:
			fallthrough
		case imap.MailboxUnseen:
			c, err := m.count(item == imap.MailboxUnseen)
			if err != nil {
				return nil, err
			}
			s.Messages = c
		case imap.MailboxRecent:
			s.Recent = 0
		case imap.MailboxUidNext:
		case imap.MailboxUidValidity:
		}
	}
	return s, nil
}

// TODO

func (m *mailbox) SetSubscribed(subscribed bool) error {
	return ErrUnimplemented
}

// Check doesn't do anything.
func (m *mailbox) Check() error {
	return nil
}

// List messages retrieves all of the requested messages in the folder.
func (m *mailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []string, ch chan<- *imap.Message) error {
	defer close(ch)
	return m.walk(uid, seqset, func(seqNum uint32, msg *db.Message) error {
		n, err := m.message(msg, seqNum, items)
		if err != nil {
			return err
		}
		ch <- n
		return nil
	})
}

// SearchMessages is unimplemented.
func (m *mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	return nil, ErrUnimplemented
}

// CreateMessage is unimplemented.
func (m *mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return ErrUnimplemented
}

// UpdateMessagesFlags is unimplemented.
func (m *mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	return ErrUnimplemented
}

// CopyMessages is unimplemented.
func (m *mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	//...
	return ErrUnimplemented
}

// Expunge is unimplemented.
func (m *mailbox) Expunge() error {
	return ErrUnimplemented
}
