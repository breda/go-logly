package discover

import (
	"fmt"
	"net"

	"github.com/hashicorp/serf/serf"
)

func (m *Membership) setupSerf() error {
	addr, err := net.ResolveTCPAddr("tcp4", m.BindAddr)
	if err != nil {
		return err
	}

	config := serf.DefaultConfig()
	config.Init()

	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port

	m.events = make(chan serf.Event)
	config.EventCh = m.events

	config.Tags = m.Tags
	config.NodeName = m.NodeName

	m.serf, err = serf.Create(config)
	if err != nil {
		return err
	}

	go m.eventHandler()
	if m.StartJoinAddrs != nil {
		_, err = m.serf.Join(m.StartJoinAddrs, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Membership) eventHandler() {
	for event := range m.events {

		m.Logger.Info().Msg(fmt.Sprintf("new event: %s", event.String()))

		switch event.EventType() {
		case serf.EventMemberJoin:
			for _, member := range event.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}

				m.handleJoin(member)
			}

		case serf.EventMemberFailed, serf.EventMemberLeave:
			for _, member := range event.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					return
				}

				m.handleLeave(member)
			}
		}
	}
}

func (m *Membership) handleJoin(member serf.Member) {
	if err := m.handler.Join(member.Name, member.Tags["rpc_addr"]); err != nil {
		m.logErr(err, "node join handle failed", member)
	}
}

func (m *Membership) handleLeave(member serf.Member) {
	if err := m.handler.Leave(member.Name); err != nil {
		m.logErr(err, "node leave handle failed", member)
	}
}
