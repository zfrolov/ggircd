package irc

type Channel struct {
  Name string

  //Mode Mode

  Topic string

  Limit int
  Key   string

  BanNick string
  BanUser string
  BanHost string

  Clients map[*User]bool
  Ops     map[*User]bool
  Voice   map[*User]bool

  Sink Sink
}

type ChannelSink struct {
  Channel *Channel
}

func (s *ChannelSink) Send(msg Message) {
  for user := range s.Channel.Clients {
    user.Sink.Send(msg)
  }
}

// ForUsers iterates over all of the users in the channel and passes a pointer
// to each to the supplied callback.
func (ch Channel) ForChannels(callback func(*User)) {
  for u := range ch.Clients {
    callback(u)
  }
}
