package discovery

import (
	"errors"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/serf/serf"
)

type Config struct {
	NodeName       string
	BindAddr       string
	Tags           map[string]string
	StartJoinAddrs []string
}

type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
}

func New(handler Handler, config Config) (*Membership, error) {
	c := &Membership{
		Config:  config,
		handler: handler,
	}
	if err := c.setupSerf(); err != nil {
		return nil, err
	}
	return c, nil
}

type Handler interface {
	Join(name, addr string) error
	Leave(name string) error
}

func (m *Membership) setupSerf() (err error) {
	log.Info().Msg("setting up serf...")
	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)
	if err != nil {
		log.Error().Err(err).Msg("unable to resolve address")
		return errors.New("unable to setup serf")
	}
	log.Info().Str("bind-addr", m.BindAddr).Str("addr", addr.String()).Msg("resolved address for serf...")
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	m.events = make(chan serf.Event)
	config.EventCh = m.events
	config.Tags = m.Tags
	config.NodeName = m.Config.NodeName
	m.serf, err = serf.Create(config)
	if err != nil {
		log.Error().Err(err).Msg("unable to create serf config")
		return errors.New("unable to setup serf")
	}

	go m.eventHandler()
	if m.StartJoinAddrs != nil {
		_, err = m.serf.Join(m.StartJoinAddrs, true)
		if err != nil {
			log.Error().Err(err).Strs("addrs", m.StartJoinAddrs).Msg("unable join")
			return errors.New("unable to setup serf")
		}
	}
	log.Info().Msg("serf setup.")
	return nil

}

func (m *Membership) eventHandler() {
	for e := range m.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleJoin(member)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					return
				}
				m.handleLeave(member)
			}
		}
	}
}

func (m *Membership) handleJoin(member serf.Member) {
	if err := m.handler.Join(
		member.Name,
		member.Tags["rpc_addr"],
	); err != nil {
		m.logError(err, "failed to join", member)
	}
	log.Info().Str("name", member.Name).Str("addr", member.Tags["rpc_addr"]).Msg("member joined")
}
func (m *Membership) handleLeave(member serf.Member) {
	if err := m.handler.Leave(
		member.Name,
	); err != nil {
		m.logError(err, "failed to leave", member)
	}
	log.Info().Str("name", member.Name).Str("addr", member.Tags["rpc_addr"]).Msg("member left")
}

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}
func (m *Membership) Members() []serf.Member {
	log.Info().Msg("finding members...")
	return m.serf.Members()
}
func (m *Membership) Leave() error {
	log.Info().Msg("leaving...")
	return m.serf.Leave()
}
func (m *Membership) logError(err error, msg string, member serf.Member) {
	log.Error().Err(err).Str("name", member.Name).Str("addr", member.Tags["rpc_addr"]).Msg(msg)
}
