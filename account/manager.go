package account

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
)


type Manager struct {
	component.Base
	channel *starx.Channel
}

type CreateNew struct {
	Name string
}

func NewManager() *Manager {
	return &Manager{
		channel: starx.ChannelService.NewChannel("Account"),
	}
}

func (m *Manager) Register(s *session.Session, msg *CreateNew) error {

	// try to create account here.

	_, credentials, _ := NewAccount(msg.Name)

	return s.Response(response.OKWithData(credentials))
}

func (m *Manager) Authenticate(s *session.Session, msg *Credentials) error {

	acc, err := RetrieveAccount(msg.ID, msg.Key)

	if err != nil {
		return s.Response(response.InternalError())
	}

	s.Bind(acc.ID)     // binding session uid
	m.channel.Add(s) // add session to channel


	s.Set("account", &acc)

	return s.Response(response.OK())
}