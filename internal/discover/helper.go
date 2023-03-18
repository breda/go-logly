package discover

import "github.com/hashicorp/serf/serf"

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

func (m *Membership) Leave() error {
	return m.serf.Leave()
}

func (m *Membership) logErr(err error, msg string, member serf.Member) {
	m.Logger.Fatal().
		Err(err).
		Str("message", msg).
		Str("name", member.Name).
		Str("rpc_addr", member.Tags["rpc_addr"]).
		Send()
}
