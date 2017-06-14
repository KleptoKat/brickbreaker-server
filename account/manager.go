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

type IDRequest struct {
	AccountID int64 `json:"account_id"`
}

type IsAuthenticatedRequest struct {}
type IsAuthenticatedResponse struct {
	Authenticated bool `json:"authenticated"`
}

func NewManager() *Manager {
	return &Manager{
		channel: starx.ChannelService.NewChannel("Account"),
	}
}

func (m *Manager) Register(s *session.Session, msg *CreateNew) error {
	// try to create account here.

	if s.Uid > 0 {
		return s.Response(response.BadRequest())
	}

	_, credentials, err := NewAccount(msg.Name)

	if err != nil {
		return s.Response(response.BadRequestWithError(err))
	}

	return s.Response(response.OKWithData(credentials))
}

func (m *Manager) Authenticate(s *session.Session, msg *Credentials) error {

	if s.Uid > 0 {
		return s.Response(response.BadRequest())
	}

	acc := AuthenticateAccount(msg.ID, msg.Key)
	if acc == nil {
		return s.Response(response.BadRequest())
	}

	s.Bind(acc.ID)

	s.Set("Name", acc.Name)
	m.channel.Add(s)

	return s.Response(response.OK())
}


func (m *Manager) IsAuthenticated(s *session.Session, msg *IsAuthenticatedRequest) error {

	return s.Response(response.OKWithData(IsAuthenticatedResponse{
		Authenticated:s.Uid > 0,
	}))
}


func (m *Manager) RetrieveAccountInfo(s *session.Session, msg *IDRequest) error {

	acc := RetrieveAccount(msg.AccountID)

	if acc == nil {
		return s.Response(response.BadRequest())
	}

	return s.Response(response.OKWithData(acc))
}