package slack

import (
	"errors"

	"github.com/botopolis/bot"
	"github.com/nlopes/slack"
)

func (a *Adapter) parseRoom(m *bot.Message) error {
	if len(m.Room) > 0 {
		if m.Room[0] == 'C' || m.Room[0] == 'D' {
			return nil
		}
	}

	if m.Room == "" {
		if msg, ok := m.Envelope.(slack.Message); ok {
			m.Room = msg.Channel
		}
		return nil
	}

	if ch, ok := a.Store.ChannelByName(m.Room); ok {
		m.Room = ch.ID
	}

	return errors.New("Room not found")
}

func (a *Adapter) parseUser(m *bot.Message) error {
	if len(m.User) > 0 {
		if m.User[0] == 'U' {
			return nil
		}
	}

	if m.User == "" {
		if msg, ok := m.Envelope.(slack.Message); ok {
			m.User = msg.User
		}
		return nil
	}

	if u, ok := a.Store.UserByName(m.User); ok {
		m.User = u.ID
	}

	return errors.New("User not found")
}

func (a *Adapter) parseDM(m *bot.Message) error {
	if len(m.Room) > 0 {
		if m.Room[0] == 'D' {
			return nil
		}
	}

	if im, ok := a.Store.IMByUserID(m.User); ok {
		m.Room = im.ID
		return nil
	}

	_, _, imID, err := a.Client.OpenIMChannel(m.User)
	if err != nil {
		return errors.New("Couldn't open IM to User: " + m.User)
	}
	m.Room = imID

	return nil
}