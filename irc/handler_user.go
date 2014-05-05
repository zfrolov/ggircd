package irc

type commandMap map[string]func(State, *User, Connection, Message) Handler

// UserHandler is a handler that handles messages coming from a user connection
// that has successfully associated with the client.
type UserHandler struct {
	state    chan State
	nick     string
	commands commandMap
}

func NewUserHandler(state chan State, nick string) Handler {
	handler := &UserHandler{
		state: state,
		nick:  nick,
	}
	handler.commands = commandMap{
		CmdAway.Command:    handler.handleCmdAway,
		CmdInvite.Command:  handler.handleCmdInvite,
		CmdJoin.Command:    handler.handleCmdJoin,
		CmdKick.Command:    handler.handleCmdKick,
		CmdList.Command:    handler.handleCmdList,
		CmdMode.Command:    handler.handleCmdMode,
		CmdMotd.Command:    handler.handleCmdMotd,
		CmdNames.Command:   handler.handleCmdNames,
		CmdNick.Command:    handler.handleCmdNick,
		CmdNotice.Command:  handler.handleCmdNotice,
		CmdPart.Command:    handler.handleCmdPart,
		CmdPing.Command:    handler.handleCmdPing,
		CmdPrivMsg.Command: handler.handleCmdPrivMsg,
		CmdQuit.Command:    handler.handleCmdQuit,
		CmdTopic.Command:   handler.handleCmdTopic,
		CmdWho.Command:     handler.handleCmdWho,
	}
	return handler
}

func (h *UserHandler) Closed(conn Connection) {
	state := <-h.state
	defer func() { h.state <- state }()

	state.RemoveUser(state.GetUser(h.nick))
	conn.Kill()
}

func (h *UserHandler) Handle(conn Connection, msg Message) Handler {
	state := <-h.state
	defer func() { h.state <- state }()

	command := h.commands[msg.Command]
	if command == nil {
		return h
	}

	user := state.GetUser(h.nick)
	newHandler := command(state, user, conn, msg)
	h.nick = user.Nick
	return newHandler
}
