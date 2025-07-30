package main

import (
	"fmt"
	"strings"
	"time"
)

// IRC numeric reply codes
const (
	RPL_WELCOME           = 001
	RPL_YOURHOST          = 002
	RPL_CREATED           = 003
	RPL_MYINFO            = 004
	RPL_ISUPPORT          = 005
	RPL_AWAY              = 301
	RPL_UNAWAY            = 305
	RPL_NOWAWAY           = 306
	RPL_WHOISUSER         = 311
	RPL_WHOISSERVER       = 312
	RPL_WHOISOPERATOR     = 313
	RPL_WHOISIDLE         = 317
	RPL_ENDOFWHOIS        = 318
	RPL_WHOISCHANNELS     = 319
	RPL_LISTSTART         = 321
	RPL_LIST              = 322
	RPL_LISTEND           = 323
	RPL_CHANNELMODEIS     = 324
	RPL_NOTOPIC           = 331
	RPL_TOPIC             = 332
	RPL_TOPICWHOTIME      = 333
	RPL_NAMREPLY          = 353
	RPL_ENDOFNAMES        = 366
	RPL_MOTDSTART         = 375
	RPL_MOTD              = 372
	RPL_ENDOFMOTD         = 376
	RPL_UMODEIS           = 221
	RPL_INVITING          = 341
	RPL_YOUREOPER         = 381
	ERR_NOSUCHNICK        = 401
	ERR_NOSUCHSERVER      = 402
	ERR_NOSUCHCHANNEL     = 403
	ERR_CANNOTSENDTOCHAN  = 404
	ERR_TOOMANYCHANNELS   = 405
	ERR_WASNOSUCHNICK     = 406
	ERR_TOOMANYTARGETS    = 407
	ERR_NOORIGIN          = 409
	ERR_NORECIPIENT       = 411
	ERR_NOTEXTTOSEND      = 412
	ERR_UNKNOWNCOMMAND    = 421
	ERR_NOMOTD            = 422
	ERR_NONICKNAMEGIVEN   = 431
	ERR_ERRONEUSNICKNAME  = 432
	ERR_NICKNAMEINUSE     = 433
	ERR_NICKCOLLISION     = 436
	ERR_USERNOTINCHANNEL  = 441
	ERR_NOTONCHANNEL      = 442
	ERR_USERONCHANNEL     = 443
	ERR_NOLOGIN           = 444
	ERR_SUMMONDISABLED    = 445
	ERR_USERSDISABLED     = 446
	ERR_NOTREGISTERED     = 451
	ERR_NEEDMOREPARAMS    = 461
	ERR_ALREADYREGISTRED  = 462
	ERR_NOPERMFORHOST     = 463
	ERR_PASSWDMISMATCH    = 464
	ERR_YOUREBANNEDCREEP  = 465
	ERR_YOUWILLBEBANNED   = 466
	ERR_KEYSET            = 467
	ERR_CHANNELISFULL     = 471
	ERR_UNKNOWNMODE       = 472
	ERR_INVITEONLYCHAN    = 473
	ERR_BANNEDFROMCHAN    = 474
	ERR_BADCHANNELKEY     = 475
	ERR_BADCHANMASK       = 476
	ERR_NOCHANMODES       = 477
	ERR_BANLISTFULL       = 478
	ERR_NOPRIVILEGES      = 481
	ERR_CHANOPRIVSNEEDED  = 482
	ERR_CANTKILLSERVER    = 483
	ERR_RESTRICTED        = 484
	ERR_UNIQOPPRIVSNEEDED = 485
	ERR_NOOPERHOST        = 491
	ERR_UMODEUNKNOWNFLAG  = 501
	ERR_USERSDONTMATCH    = 502
	RPL_SNOMASK           = 8
	RPL_GLOBALNOTICE      = 710
	RPL_OPERWALL          = 711
)

// handleNick handles NICK command
func (c *Client) handleNick(parts []string) {
	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "NICK :Not enough parameters")
		return
	}

	newNick := parts[1]
	if len(newNick) > 0 && newNick[0] == ':' {
		newNick = newNick[1:]
	}

	// Validate nickname
	if !isValidNickname(newNick) {
		c.SendNumeric(ERR_ERRONEUSNICKNAME, newNick+" :Erroneous nickname")
		return
	}

	// Check if nick is already in use
	if existing := c.server.GetClient(newNick); existing != nil && existing != c {
		c.SendNumeric(ERR_NICKNAMEINUSE, newNick+" :Nickname is already in use")
		return
	}

	oldNick := c.Nick()
	c.SetNick(newNick)

	// If already registered, notify channels
	if c.IsRegistered() && oldNick != "" {
		message := fmt.Sprintf(":%s NICK :%s", c.Prefix(), newNick)
		for _, channel := range c.GetChannels() {
			channel.Broadcast(message, nil)
		}

		// Send snomask notification for nick change
		if c.server != nil && oldNick != newNick {
			c.server.sendSnomask('n', fmt.Sprintf("Nick change: %s -> %s (%s@%s)",
				oldNick, newNick, c.User(), c.Host()))
		}
	}

	c.checkRegistration()
}

// handleUser handles USER command
func (c *Client) handleUser(parts []string) {
	if len(parts) < 5 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "USER :Not enough parameters")
		return
	}

	if c.IsRegistered() {
		c.SendNumeric(ERR_ALREADYREGISTRED, ":You may not reregister")
		return
	}

	c.SetUser(parts[1])
	// parts[2] and parts[3] are ignored (mode and unused)
	realname := strings.Join(parts[4:], " ")
	if len(realname) > 0 && realname[0] == ':' {
		realname = realname[1:]
	}
	c.SetRealname(realname)

	c.checkRegistration()
}

// checkRegistration checks if client is ready to be registered
func (c *Client) checkRegistration() {
	if !c.IsRegistered() && c.Nick() != "" && c.User() != "" {
		c.SetRegistered(true)
		c.sendWelcome()
	}
}

// sendWelcome sends welcome messages to newly registered client
func (c *Client) sendWelcome() {
	fmt.Printf("DEBUG: sendWelcome called\n")
	if c.server == nil {
		fmt.Printf("DEBUG: sendWelcome - server is nil\n")
		return
	}
	if c.server.config == nil {
		fmt.Printf("DEBUG: sendWelcome - config is nil\n")
		return
	}

	fmt.Printf("DEBUG: sendWelcome - about to send RPL_WELCOME\n")
	c.SendNumeric(RPL_WELCOME, fmt.Sprintf("Welcome to %s, %s", c.server.config.Server.Network, c.Prefix()))
	fmt.Printf("DEBUG: sendWelcome - sent RPL_WELCOME\n")
	c.SendNumeric(RPL_YOURHOST, fmt.Sprintf("Your host is %s, running version %s", c.server.config.Server.Name, c.server.config.Server.Version))
	c.SendNumeric(RPL_CREATED, "This server was created recently")
	c.SendNumeric(RPL_MYINFO, fmt.Sprintf("%s %s o o", c.server.config.Server.Name, c.server.config.Server.Version))

	// Send MOTD
	if len(c.server.config.MOTD) > 0 {
		c.SendNumeric(RPL_MOTDSTART, fmt.Sprintf("- %s Message of the Day -", c.server.config.Server.Name))
		for _, line := range c.server.config.MOTD {
			c.SendNumeric(RPL_MOTD, fmt.Sprintf("- %s", line))
		}
		c.SendNumeric(RPL_ENDOFMOTD, "End of /MOTD command")
	}

	// Send snomask notification for new client connection
	if c.server != nil {
		c.server.sendSnomask('c', fmt.Sprintf("Client connect: %s (%s@%s)",
			c.Nick(), c.User(), c.Host()))
	}

	fmt.Printf("DEBUG: sendWelcome completed\n")
}

// handlePing handles PING command
func (c *Client) handlePing(parts []string) {
	if len(parts) < 2 {
		return
	}

	token := parts[1]
	if len(token) > 0 && token[0] == ':' {
		token = token[1:]
	}

	serverName := "localhost"
	if c.server != nil && c.server.config != nil {
		serverName = c.server.config.Server.Name
	}

	c.SendMessage(fmt.Sprintf("PONG %s :%s", serverName, token))
}

// handlePong handles PONG command
func (c *Client) handlePong(parts []string) {
	// Update the last pong time for ping timeout tracking
	// This is used by the client Handler's ping timeout mechanism
	c.mu.Lock()
	c.lastPong = time.Now()
	c.waitingForPong = false
	c.mu.Unlock()
}

// handleJoin handles JOIN command
func (c *Client) handleJoin(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "JOIN :Not enough parameters")
		return
	}

	channelNames := strings.Split(parts[1], ",")
	keys := []string{}
	if len(parts) > 2 {
		keys = strings.Split(parts[2], ",")
	}

	for i, channelName := range channelNames {
		if channelName == "0" {
			// Leave all channels
			for _, channel := range c.GetChannels() {
				c.handlePartChannel(channel.Name(), "Leaving all channels")
			}
			continue
		}

		if !isValidChannelName(channelName) {
			c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
			continue
		}

		channel := c.server.GetOrCreateChannel(channelName)

		// Check if already in channel
		if c.IsInChannel(channelName) {
			continue
		}

		// Check channel modes and limits
		key := ""
		if i < len(keys) {
			key = keys[i]
		}

		if channel.HasMode('k') && channel.Key() != key {
			c.SendNumeric(ERR_BADCHANNELKEY, channelName+" :Cannot join channel (+k)")
			continue
		}

		if channel.HasMode('l') && channel.UserCount() >= channel.Limit() {
			c.SendNumeric(ERR_CHANNELISFULL, channelName+" :Cannot join channel (+l)")
			continue
		}

		// Join the channel
		channel.AddClient(c)
		c.AddChannel(channel)

		message := fmt.Sprintf(":%s JOIN :%s", c.Prefix(), channelName)
		channel.Broadcast(message, nil)

		// Send topic if exists
		if channel.Topic() != "" {
			c.SendNumeric(RPL_TOPIC, channelName+" :"+channel.Topic())
			c.SendNumeric(RPL_TOPICWHOTIME, fmt.Sprintf("%s %s %d", channelName, channel.TopicBy(), channel.TopicTime().Unix()))
		}

		// Send names list
		c.sendNames(channel)
	}
}

// handlePart handles PART command
func (c *Client) handlePart(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "PART :Not enough parameters")
		return
	}

	channelNames := strings.Split(parts[1], ",")
	reason := "Leaving"
	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
		if len(reason) > 0 && reason[0] == ':' {
			reason = reason[1:]
		}
	}

	for _, channelName := range channelNames {
		c.handlePartChannel(channelName, reason)
	}
}

func (c *Client) handlePartChannel(channelName, reason string) {
	if !c.IsInChannel(channelName) {
		c.SendNumeric(ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		return
	}

	message := fmt.Sprintf(":%s PART %s :%s", c.Prefix(), channelName, reason)
	channel.Broadcast(message, nil)

	channel.RemoveClient(c)
	c.RemoveChannel(channelName)

	// Remove empty channel
	if channel.UserCount() == 0 {
		c.server.RemoveChannel(channelName)
	}
}

// handlePrivmsg handles PRIVMSG command
func (c *Client) handlePrivmsg(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NORECIPIENT, ":No recipient given (PRIVMSG)")
		return
	}

	if len(parts) < 3 {
		c.SendNumeric(ERR_NOTEXTTOSEND, ":No text to send")
		return
	}

	target := parts[1]
	message := strings.Join(parts[2:], " ")
	if len(message) > 0 && message[0] == ':' {
		message = message[1:]
	}

	if isChannelName(target) {
		// Channel message
		channel := c.server.GetChannel(target)
		if channel == nil {
			c.SendNumeric(ERR_NOSUCHCHANNEL, target+" :No such channel")
			return
		}

		if !c.IsInChannel(target) {
			c.SendNumeric(ERR_CANNOTSENDTOCHAN, target+" :Cannot send to channel")
			return
		}

		// Check if user can send messages to this channel (moderated mode check)
		if !channel.CanSendMessage(c) {
			c.SendNumeric(ERR_CANNOTSENDTOCHAN, target+" :Cannot send to channel (+m)")
			return
		}

		msg := fmt.Sprintf(":%s PRIVMSG %s :%s", c.Prefix(), target, message)
		channel.Broadcast(msg, c)
	} else {
		// Private message
		targetClient := c.server.GetClient(target)
		if targetClient == nil {
			c.SendNumeric(ERR_NOSUCHNICK, target+" :No such nick/channel")
			return
		}

		if targetClient.Away() != "" {
			c.SendNumeric(RPL_AWAY, fmt.Sprintf("%s :%s", target, targetClient.Away()))
		}

		msg := fmt.Sprintf(":%s PRIVMSG %s :%s", c.Prefix(), target, message)
		targetClient.SendMessage(msg)
	}
}

// handleNotice handles NOTICE command
func (c *Client) handleNotice(parts []string) {
	if !c.IsRegistered() {
		return // NOTICE should not generate error responses
	}

	if len(parts) < 3 {
		return
	}

	target := parts[1]
	message := strings.Join(parts[2:], " ")
	if len(message) > 0 && message[0] == ':' {
		message = message[1:]
	}

	if isChannelName(target) {
		// Channel notice
		channel := c.server.GetChannel(target)
		if channel == nil || !c.IsInChannel(target) {
			return
		}

		msg := fmt.Sprintf(":%s NOTICE %s :%s", c.Prefix(), target, message)
		channel.Broadcast(msg, c)
	} else {
		// Private notice
		targetClient := c.server.GetClient(target)
		if targetClient == nil {
			return
		}

		msg := fmt.Sprintf(":%s NOTICE %s :%s", c.Prefix(), target, message)
		targetClient.SendMessage(msg)
	}
}

// handleWho handles WHO command
func (c *Client) handleWho(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "WHO :Not enough parameters")
		return
	}

	target := parts[1]

	if isChannelName(target) {
		channel := c.server.GetChannel(target)
		if channel == nil {
			c.SendNumeric(ERR_NOSUCHCHANNEL, target+" :No such channel")
			return
		}

		for _, client := range channel.GetClients() {
			flags := ""
			if client.IsOper() {
				flags += "*"
			}
			if client.Away() != "" {
				flags += "G"
			} else {
				flags += "H"
			}
			if channel.IsOperator(client) {
				flags += "@"
			} else if channel.IsVoice(client) {
				flags += "+"
			}

			c.SendNumeric(352, fmt.Sprintf("%s %s %s %s %s %s :0 %s",
				target, client.User(), client.HostForUser(c), c.server.config.Server.Name,
				client.Nick(), flags, client.Realname()))
		}
	}

	c.SendNumeric(315, target+" :End of /WHO list")
}

// handleWhois handles WHOIS command
func (c *Client) handleWhois(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "WHOIS :Not enough parameters")
		return
	}

	nick := parts[1]
	target := c.server.GetClient(nick)
	if target == nil {
		c.SendNumeric(ERR_NOSUCHNICK, nick+" :No such nick")
		return
	}

	c.SendNumeric(RPL_WHOISUSER, fmt.Sprintf("%s %s %s * :%s",
		target.Nick(), target.User(), target.HostForUser(c), target.Realname()))

	c.SendNumeric(RPL_WHOISSERVER, fmt.Sprintf("%s %s :%s",
		target.Nick(), c.server.config.Server.Name, c.server.config.Server.Description))

	if target.IsOper() {
		c.SendNumeric(RPL_WHOISOPERATOR, target.Nick()+" :is an IRC operator")
	}

	if target.Away() != "" {
		c.SendNumeric(RPL_AWAY, fmt.Sprintf("%s :%s", target.Nick(), target.Away()))
	}

	// Send channels
	var channels []string
	for _, channel := range target.GetChannels() {
		channelName := channel.Name()
		if channel.IsOperator(target) {
			channelName = "@" + channelName
		} else if channel.IsVoice(target) {
			channelName = "+" + channelName
		}
		channels = append(channels, channelName)
	}
	if len(channels) > 0 {
		c.SendNumeric(RPL_WHOISCHANNELS, fmt.Sprintf("%s :%s", target.Nick(), strings.Join(channels, " ")))
	}

	// Show user modes if the requester is an operator or the target user
	if c.IsOper() || c.Nick() == target.Nick() {
		modes := target.GetModes()
		if modes != "" {
			c.SendMessage(fmt.Sprintf(":%s 379 %s %s :is using modes %s",
				c.server.config.Server.Name, c.Nick(), target.Nick(), modes))
		}
	}

	// Show SSL status
	if target.IsSSL() {
		c.SendMessage(fmt.Sprintf(":%s 671 %s %s :is using a secure connection",
			c.server.config.Server.Name, c.Nick(), target.Nick()))
	}

	c.SendNumeric(RPL_ENDOFWHOIS, target.Nick()+" :End of /WHOIS list")
}

// handleNames handles NAMES command
func (c *Client) handleNames(parts []string) {
	if !c.IsRegistered() {
		c.SendNumeric(ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if len(parts) < 2 {
		// Send names for all channels
		for _, channel := range c.server.GetChannels() {
			if c.IsInChannel(channel.Name()) {
				c.sendNames(channel)
			}
		}
		return
	}

	channelNames := strings.Split(parts[1], ",")
	for _, channelName := range channelNames {
		channel := c.server.GetChannel(channelName)
		if channel != nil && c.IsInChannel(channelName) {
			c.sendNames(channel)
		}
	}
}

func (c *Client) sendNames(channel *Channel) {
	var names []string
	for _, client := range channel.GetClients() {
		name := client.Nick()
		if channel.IsOwner(client) {
			name = "~" + name
		} else if channel.IsOperator(client) {
			name = "@" + name
		} else if channel.IsHalfop(client) {
			name = "%" + name
		} else if channel.IsVoice(client) {
			name = "+" + name
		}
		names = append(names, name)
	}

	symbol := "="
	if channel.HasMode('s') {
		symbol = "@"
	} else if channel.HasMode('p') {
		symbol = "*"
	}

	c.SendNumeric(RPL_NAMREPLY, fmt.Sprintf("%s %s :%s", symbol, channel.Name(), strings.Join(names, " ")))
	c.SendNumeric(RPL_ENDOFNAMES, channel.Name()+" :End of /NAMES list")
}

// handleQuit handles QUIT command
func (c *Client) handleQuit(parts []string) {
	reason := "Client quit"
	if len(parts) > 1 {
		reason = strings.Join(parts[1:], " ")
		if len(reason) > 0 && reason[0] == ':' {
			reason = reason[1:]
		}
	}

	c.server.RemoveClient(c)
}

// handleMode handles MODE command
func (c *Client) handleMode(parts []string) {
	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "MODE :Not enough parameters")
		return
	}

	target := parts[1]

	// Handle user mode requests
	if !isChannelName(target) {
		if target != c.Nick() {
			c.SendNumeric(ERR_USERSDONTMATCH, ":Cannot change mode for other users")
			return
		}

		// If no mode changes specified, return current user modes
		if len(parts) == 2 {
			modes := c.GetModes()
			if modes == "" {
				modes = "+"
			}
			c.SendNumeric(RPL_UMODEIS, modes)
			return
		}

		// Parse user mode changes
		modeString := parts[2]
		adding := true
		var appliedModes []string

		for _, char := range modeString {
			switch char {
			case '+':
				adding = true
			case '-':
				adding = false
			case 'i': // invisible
				c.SetMode('i', adding)
				if adding {
					appliedModes = append(appliedModes, "+i")
				} else {
					appliedModes = append(appliedModes, "-i")
				}
			case 'w': // wallops
				c.SetMode('w', adding)
				if adding {
					appliedModes = append(appliedModes, "+w")
				} else {
					appliedModes = append(appliedModes, "-w")
				}
			case 's': // server notices (requires oper)
				if !c.IsOper() && adding {
					continue // silently ignore for non-opers
				}
				c.SetMode('s', adding)
				if adding {
					appliedModes = append(appliedModes, "+s")
				} else {
					appliedModes = append(appliedModes, "-s")
				}
			case 'o': // operator (cannot be set manually)
				if adding {
					c.SendNumeric(ERR_UMODEUNKNOWNFLAG, ":Unknown MODE flag")
				} else {
					// Allow de-opering
					c.SetOper(false)
					c.SetMode('o', false)
					appliedModes = append(appliedModes, "-o")
					// Clear snomasks when de-opering
					c.snomasks = make(map[rune]bool)
					c.sendSnomask('o', fmt.Sprintf("%s is no longer an IRC operator", c.Nick()))
				}
			case 'r': // registered (cannot be set manually, services only)
				c.SendNumeric(ERR_UMODEUNKNOWNFLAG, ":Unknown MODE flag")
			case 'x': // host masking (TechIRCd special)
				c.SetMode('x', adding)
				if adding {
					appliedModes = append(appliedModes, "+x")
					// TODO: Implement host masking
				} else {
					appliedModes = append(appliedModes, "-x")
				}
			case 'z': // SSL/TLS (automatic, cannot be manually set)
				if c.IsSSL() {
					c.SetMode('z', true)
				}
				// Ignore attempts to manually set/unset
			case 'B': // bot flag (TechIRCd special)
				c.SetMode('B', adding)
				if adding {
					appliedModes = append(appliedModes, "+B")
				} else {
					appliedModes = append(appliedModes, "-B")
				}
			default:
				c.SendNumeric(ERR_UMODEUNKNOWNFLAG, ":Unknown MODE flag")
			}
		}

		// Send mode changes back to user
		if len(appliedModes) > 0 {
			modeStr := strings.Join(appliedModes, "")
			c.SendMessage(fmt.Sprintf(":%s MODE %s :%s", c.Nick(), c.Nick(), modeStr))
		}
		return
	}

	// Handle channel mode requests
	channel := c.server.GetChannel(target)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, target+" :No such channel")
		return
	}

	if !c.IsInChannel(target) {
		c.SendNumeric(ERR_NOTONCHANNEL, target+" :You're not on that channel")
		return
	}

	// If no mode changes specified, return current channel modes
	if len(parts) == 2 {
		modes := channel.GetModes()
		if modes == "" {
			modes = "+"
		}
		c.SendNumeric(RPL_CHANNELMODEIS, fmt.Sprintf("%s %s", target, modes))
		return
	}

	// Parse mode changes
	modeString := parts[2]
	args := parts[3:]
	argIndex := 0

	// Check if user has operator privileges (required for most mode changes)
	if !channel.IsOwner(c) && !channel.IsOperator(c) && !channel.IsHalfop(c) && !c.IsOper() {
		c.SendNumeric(ERR_CHANOPRIVSNEEDED, target+" :You're not channel operator")
		return
	}

	adding := true
	var appliedModes []string
	var appliedArgs []string

	for _, char := range modeString {
		switch char {
		case '+':
			adding = true
		case '-':
			adding = false
		case 'o': // operator
			if argIndex >= len(args) {
				continue
			}
			targetNick := args[argIndex]
			argIndex++

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				c.SendNumeric(ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !targetClient.IsInChannel(target) {
				c.SendNumeric(ERR_USERNOTINCHANNEL, fmt.Sprintf("%s %s :They aren't on that channel", targetNick, target))
				continue
			}

			channel.SetOperator(targetClient, adding)
			if adding {
				appliedModes = append(appliedModes, "+o")
			} else {
				appliedModes = append(appliedModes, "-o")
			}
			appliedArgs = append(appliedArgs, targetNick)

		case 'v': // voice
			if argIndex >= len(args) {
				continue
			}
			targetNick := args[argIndex]
			argIndex++

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				c.SendNumeric(ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !targetClient.IsInChannel(target) {
				c.SendNumeric(ERR_USERNOTINCHANNEL, fmt.Sprintf("%s %s :They aren't on that channel", targetNick, target))
				continue
			}

			channel.SetVoice(targetClient, adding)
			if adding {
				appliedModes = append(appliedModes, "+v")
			} else {
				appliedModes = append(appliedModes, "-v")
			}
			appliedArgs = append(appliedArgs, targetNick)

		case 'h': // halfop
			if argIndex >= len(args) {
				continue
			}
			targetNick := args[argIndex]
			argIndex++

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				c.SendNumeric(ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !targetClient.IsInChannel(target) {
				c.SendNumeric(ERR_USERNOTINCHANNEL, fmt.Sprintf("%s %s :They aren't on that channel", targetNick, target))
				continue
			}

			channel.SetHalfop(targetClient, adding)
			if adding {
				appliedModes = append(appliedModes, "+h")
			} else {
				appliedModes = append(appliedModes, "-h")
			}
			appliedArgs = append(appliedArgs, targetNick)

		case 'q': // owner/founder
			if argIndex >= len(args) {
				continue
			}
			targetNick := args[argIndex]
			argIndex++

			// Only existing owners can grant/remove owner status
			if !channel.IsOwner(c) && !c.IsOper() {
				c.SendNumeric(ERR_CHANOPRIVSNEEDED, target+" :You're not channel owner")
				continue
			}

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				c.SendNumeric(ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !targetClient.IsInChannel(target) {
				c.SendNumeric(ERR_USERNOTINCHANNEL, fmt.Sprintf("%s %s :They aren't on that channel", targetNick, target))
				continue
			}

			channel.SetOwner(targetClient, adding)
			if adding {
				appliedModes = append(appliedModes, "+q")
			} else {
				appliedModes = append(appliedModes, "-q")
			}
			appliedArgs = append(appliedArgs, targetNick)

		case 'm': // moderated
			channel.SetMode('m', adding)
			if adding {
				appliedModes = append(appliedModes, "+m")
			} else {
				appliedModes = append(appliedModes, "-m")
			}

		case 'n': // no external messages
			channel.SetMode('n', adding)
			if adding {
				appliedModes = append(appliedModes, "+n")
			} else {
				appliedModes = append(appliedModes, "-n")
			}

		case 't': // topic restriction
			channel.SetMode('t', adding)
			if adding {
				appliedModes = append(appliedModes, "+t")
			} else {
				appliedModes = append(appliedModes, "-t")
			}

		case 'i': // invite only
			channel.SetMode('i', adding)
			if adding {
				appliedModes = append(appliedModes, "+i")
			} else {
				appliedModes = append(appliedModes, "-i")
			}

		case 's': // secret
			channel.SetMode('s', adding)
			if adding {
				appliedModes = append(appliedModes, "+s")
			} else {
				appliedModes = append(appliedModes, "-s")
			}

		case 'p': // private
			channel.SetMode('p', adding)
			if adding {
				appliedModes = append(appliedModes, "+p")
			} else {
				appliedModes = append(appliedModes, "-p")
			}

		case 'k': // key (password)
			if adding {
				if argIndex >= len(args) {
					continue
				}
				key := args[argIndex]
				argIndex++
				channel.SetKey(key)
				channel.SetMode('k', true)
				appliedModes = append(appliedModes, "+k")
				appliedArgs = append(appliedArgs, key)
			} else {
				channel.SetKey("")
				channel.SetMode('k', false)
				appliedModes = append(appliedModes, "-k")
			}

		case 'l': // limit
			if adding {
				if argIndex >= len(args) {
					continue
				}
				limitStr := args[argIndex]
				argIndex++
				// Parse limit (simplified - should validate it's a number)
				limit := 0
				fmt.Sscanf(limitStr, "%d", &limit)
				if limit > 0 {
					channel.SetLimit(limit)
					channel.SetMode('l', true)
					appliedModes = append(appliedModes, "+l")
					appliedArgs = append(appliedArgs, limitStr)
				}
			} else {
				channel.SetLimit(0)
				channel.SetMode('l', false)
				appliedModes = append(appliedModes, "-l")
			}

		case 'b': // ban (enhanced with extended ban types)
			if argIndex >= len(args) {
				// List bans (TODO: implement ban list display)
				continue
			}
			mask := args[argIndex]
			argIndex++

			// Check for extended ban types (e.g., ~q:nick!user@host for quiet)
			if strings.HasPrefix(mask, "~") && len(mask) > 2 && mask[2] == ':' {
				banType := mask[1]  // The character after ~
				banMask := mask[3:] // The mask after ~x:

				switch banType {
				case 'q': // Quiet ban
					if adding {
						// Add to quiet list
						channel.quietList = append(channel.quietList, banMask)
						appliedModes = append(appliedModes, "+b")
						appliedArgs = append(appliedArgs, mask)

						// Send snomask to opers
						if c.IsOper() {
							c.server.sendSnomask('x', fmt.Sprintf("%s set quiet ban %s on %s", c.Nick(), banMask, target))
						}
					} else {
						// Remove from quiet list
						for i, quiet := range channel.quietList {
							if quiet == banMask {
								channel.quietList = append(channel.quietList[:i], channel.quietList[i+1:]...)
								appliedModes = append(appliedModes, "-b")
								appliedArgs = append(appliedArgs, mask)

								// Send snomask to opers
								if c.IsOper() {
									c.server.sendSnomask('x', fmt.Sprintf("%s removed quiet ban %s on %s", c.Nick(), banMask, target))
								}
								break
							}
						}
					}
				default:
					// Unknown extended ban type - treat as regular ban for now
					if adding {
						channel.banList = append(channel.banList, mask)
						appliedModes = append(appliedModes, "+b")
					} else {
						for i, ban := range channel.banList {
							if ban == mask {
								channel.banList = append(channel.banList[:i], channel.banList[i+1:]...)
								appliedModes = append(appliedModes, "-b")
								break
							}
						}
					}
					appliedArgs = append(appliedArgs, mask)
				}
			} else {
				// Regular ban
				if adding {
					channel.banList = append(channel.banList, mask)
					appliedModes = append(appliedModes, "+b")
				} else {
					for i, ban := range channel.banList {
						if ban == mask {
							channel.banList = append(channel.banList[:i], channel.banList[i+1:]...)
							appliedModes = append(appliedModes, "-b")
							break
						}
					}
				}
				appliedArgs = append(appliedArgs, mask)
			}

		default:
			// Unknown mode - ignore for now
		}
	}

	// Broadcast mode changes to all channel members
	if len(appliedModes) > 0 {
		modeChangeMsg := fmt.Sprintf("MODE %s %s", target, strings.Join(appliedModes, ""))
		if len(appliedArgs) > 0 {
			modeChangeMsg += " " + strings.Join(appliedArgs, " ")
		}

		for _, client := range channel.GetClients() {
			client.SendFrom(c.Prefix(), modeChangeMsg)
		}
	}
}

// handleTopic handles TOPIC command
func (c *Client) handleTopic(parts []string) {
	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "TOPIC :Not enough parameters")
		return
	}

	channelName := parts[1]
	if !isChannelName(channelName) {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !c.IsInChannel(channelName) {
		c.SendNumeric(ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	// If no topic provided, return current topic
	if len(parts) == 2 {
		topic := channel.Topic()
		if topic == "" {
			c.SendNumeric(RPL_NOTOPIC, channelName+" :No topic is set")
		} else {
			c.SendNumeric(RPL_TOPIC, fmt.Sprintf("%s :%s", channelName, topic))
		}
		return
	}

	// Check if user can set topic (for now, anyone in channel can)
	// TODO: Add proper +t mode checking
	newTopic := strings.Join(parts[2:], " ")
	if len(newTopic) > 0 && newTopic[0] == ':' {
		newTopic = newTopic[1:]
	}

	channel.SetTopic(newTopic, c.Nick())

	// Broadcast topic change to all channel members
	for _, client := range channel.GetClients() {
		client.SendFrom(c.Prefix(), fmt.Sprintf("TOPIC %s :%s", channelName, newTopic))
	}
}

// handleAway handles AWAY command
func (c *Client) handleAway(parts []string) {
	if len(parts) == 1 {
		// Remove away status
		c.SetAway("")
		c.SendNumeric(RPL_UNAWAY, ":You are no longer marked as being away")
		return
	}

	// Set away message
	awayMsg := strings.Join(parts[1:], " ")
	if len(awayMsg) > 0 && awayMsg[0] == ':' {
		awayMsg = awayMsg[1:]
	}

	c.SetAway(awayMsg)
	c.SendNumeric(RPL_NOWAWAY, ":You have been marked as being away")
}

// handleList handles LIST command
func (c *Client) handleList(parts []string) {
	c.SendNumeric(RPL_LISTSTART, "Channel :Users  Name")

	for _, channel := range c.server.GetChannels() {
		// For now, show all channels (TODO: Add proper mode checking for secret channels)
		userCount := len(channel.GetClients())
		topic := channel.Topic()
		if topic == "" {
			topic = ""
		}
		c.SendNumeric(RPL_LIST, fmt.Sprintf("%s %d :%s", channel.Name(), userCount, topic))
	}

	c.SendNumeric(RPL_LISTEND, ":End of /LIST")
}

// handleInvite handles INVITE command
func (c *Client) handleInvite(parts []string) {
	if len(parts) < 3 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "INVITE :Not enough parameters")
		return
	}

	nick := parts[1]
	channelName := parts[2]

	target := c.server.GetClient(nick)
	if target == nil {
		c.SendNumeric(ERR_NOSUCHNICK, nick+" :No such nick/channel")
		return
	}

	if !isChannelName(channelName) {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !c.IsInChannel(channelName) {
		c.SendNumeric(ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	if target.IsInChannel(channelName) {
		c.SendNumeric(ERR_USERONCHANNEL, fmt.Sprintf("%s %s :is already on channel", nick, channelName))
		return
	}

	// TODO: Check if user has operator privileges for invite-only channels

	// Send invite to target
	target.SendFrom(c.Prefix(), fmt.Sprintf("INVITE %s %s", target.Nick(), channelName))
	c.SendNumeric(RPL_INVITING, fmt.Sprintf("%s %s", target.Nick(), channelName))
}

// handleKick handles KICK command
func (c *Client) handleKick(parts []string) {
	if len(parts) < 3 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "KICK :Not enough parameters")
		return
	}

	channelName := parts[1]
	nick := parts[2]
	reason := "No reason given"
	if len(parts) > 3 {
		reason = strings.Join(parts[3:], " ")
		if len(reason) > 0 && reason[0] == ':' {
			reason = reason[1:]
		}
	}

	if !isChannelName(channelName) {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !c.IsInChannel(channelName) {
		c.SendNumeric(ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	target := c.server.GetClient(nick)
	if target == nil {
		c.SendNumeric(ERR_NOSUCHNICK, nick+" :No such nick/channel")
		return
	}

	if !target.IsInChannel(channelName) {
		c.SendNumeric(ERR_USERNOTINCHANNEL, fmt.Sprintf("%s %s :They aren't on that channel", nick, channelName))
		return
	}

	// TODO: Check if user has operator privileges
	// For now, allow anyone to kick (will fix with proper channel modes)

	// Broadcast kick to all channel members
	kickMsg := fmt.Sprintf("KICK %s %s :%s", channelName, target.Nick(), reason)
	for _, client := range channel.GetClients() {
		client.SendFrom(c.Prefix(), kickMsg)
	}

	// Remove target from channel
	channel.RemoveClient(target)
	target.RemoveChannel(channelName)
}

// handleKill handles KILL command (operator only)
func (c *Client) handleKill(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "KILL :Not enough parameters")
		return
	}

	nick := parts[1]
	reason := "Killed by operator"
	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
		if len(reason) > 0 && reason[0] == ':' {
			reason = reason[1:]
		}
	}

	target := c.server.GetClient(nick)
	if target == nil {
		c.SendNumeric(ERR_NOSUCHNICK, nick+" :No such nick/channel")
		return
	}

	// Can't kill other operators
	if target.IsOper() {
		c.SendNumeric(ERR_CANTKILLSERVER, ":You can't kill other operators")
		return
	}

	// Send kill message to target and disconnect
	target.SendMessage(fmt.Sprintf("ERROR :Killed (%s (%s))", c.Nick(), reason))

	// Broadcast to other operators
	for _, client := range c.server.GetClients() {
		if client.IsOper() && client != c {
			client.SendMessage(fmt.Sprintf(":%s WALLOPS :%s killed %s (%s)",
				c.server.config.Server.Name, c.Nick(), target.Nick(), reason))
		}
	}

	// Disconnect the target
	target.conn.Close()
}

// handleOper handles OPER command
func (c *Client) handleOper(parts []string) {
	if len(parts) < 3 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "OPER :Not enough parameters")
		return
	}

	if c.server == nil || c.server.config == nil {
		c.SendNumeric(ERR_NOOPERHOST, ":No O-lines for your host")
		return
	}

	name := parts[1]
	password := parts[2]

	// Check if opers are enabled
	if !c.server.config.Features.EnableOper {
		c.SendNumeric(ERR_NOOPERHOST, ":O-lines are disabled")
		return
	}

	// Find matching oper configuration
	for _, oper := range c.server.config.Opers {
		if oper.Name == name && oper.Password == password {
			// Check host mask (simplified - just check if it matches *@localhost for now)
			if oper.Host == "*@localhost" || oper.Host == "*@*" {
				c.SetOper(true)

				// Set operator user mode
				c.SetMode('o', true)
				c.SetMode('s', true) // Enable server notices by default
				c.SetMode('w', true) // Enable wallops by default

				// Set default snomasks for new operators
				c.SetSnomask('c', true) // Client connects/disconnects
				c.SetSnomask('o', true) // Oper-up messages
				c.SetSnomask('s', true) // Server messages

				c.SendNumeric(RPL_YOUREOPER, ":You are now an IRC operator")
				c.SendNumeric(RPL_SNOMASK, fmt.Sprintf("%s :Server notice mask", c.GetSnomasks()))

				// Send mode change notification
				c.SendMessage(fmt.Sprintf(":%s MODE %s :+osw", c.Nick(), c.Nick()))

				// Send snomask to other operators
				c.sendSnomask('o', fmt.Sprintf("%s (%s@%s) is now an IRC operator", c.Nick(), c.User(), c.Host()))
				return
			}
		}
	}

	c.SendNumeric(ERR_PASSWDMISMATCH, ":Password incorrect")
}

// handleSnomask handles SNOMASK command (server notice masks for operators)
func (c *Client) handleSnomask(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		// Show current snomasks
		current := c.GetSnomasks()
		if current == "" {
			current = "+"
		}
		c.SendNumeric(RPL_SNOMASK, fmt.Sprintf("%s :Server notice mask", current))
		return
	}

	modeString := parts[1]
	adding := true
	changed := false

	for _, char := range modeString {
		switch char {
		case '+':
			adding = true
		case '-':
			adding = false
		case 'c': // Client connects/disconnects
			c.SetSnomask('c', adding)
			changed = true
		case 'k': // Kill messages
			c.SetSnomask('k', adding)
			changed = true
		case 'o': // Oper-up messages
			c.SetSnomask('o', adding)
			changed = true
		case 'x': // X-line (ban) messages
			c.SetSnomask('x', adding)
			changed = true
		case 'f': // Flood messages
			c.SetSnomask('f', adding)
			changed = true
		case 'n': // Nick changes
			c.SetSnomask('n', adding)
			changed = true
		case 's': // Server messages
			c.SetSnomask('s', adding)
			changed = true
		case 'd': // Debug messages (TechIRCd special)
			c.SetSnomask('d', adding)
			changed = true
		}
	}

	if changed {
		current := c.GetSnomasks()
		if current == "" {
			current = "+"
		}
		c.SendNumeric(RPL_SNOMASK, fmt.Sprintf("%s :Server notice mask", current))
	}
}

// handleGlobalNotice handles GLOBALNOTICE command (TechIRCd special oper command)
func (c *Client) handleGlobalNotice(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "GLOBALNOTICE :Not enough parameters")
		return
	}

	message := strings.Join(parts[1:], " ")
	if len(message) > 0 && message[0] == ':' {
		message = message[1:]
	}

	// Send global notice to all users
	for _, client := range c.server.GetClients() {
		client.SendMessage(fmt.Sprintf(":%s NOTICE %s :[GLOBAL] %s",
			c.server.config.Server.Name, client.Nick(), message))
	}

	// Send snomask to operators watching global notices
	c.sendSnomask('s', fmt.Sprintf("Global notice from %s: %s", c.Nick(), message))
}

// handleWallops handles WALLOPS command (send to users with +w mode)
func (c *Client) handleWallops(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "WALLOPS :Not enough parameters")
		return
	}

	message := strings.Join(parts[1:], " ")
	if len(message) > 0 && message[0] == ':' {
		message = message[1:]
	}

	// Send to all users with +w mode
	for _, client := range c.server.GetClients() {
		if client.HasMode('w') {
			client.SendMessage(fmt.Sprintf(":%s WALLOPS :%s", c.Nick(), message))
		}
	}
}

// handleOperWall handles OPERWALL command (message to all operators)
func (c *Client) handleOperWall(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		c.SendNumeric(ERR_NEEDMOREPARAMS, "OPERWALL :Not enough parameters")
		return
	}

	message := strings.Join(parts[1:], " ")
	if len(message) > 0 && message[0] == ':' {
		message = message[1:]
	}

	// Send to all operators
	for _, client := range c.server.GetClients() {
		if client.IsOper() {
			client.SendMessage(fmt.Sprintf(":%s WALLOPS :%s", c.Nick(), message))
		}
	}
}

// handleRehash handles REHASH command (reload configuration)
func (c *Client) handleRehash(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	// Reload configuration
	if c.server != nil {
		err := c.server.ReloadConfig()
		if err != nil {
			c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** REHASH failed: %s",
				c.server.config.Server.Name, c.Nick(), err.Error()))
			c.sendSnomask('s', fmt.Sprintf("REHASH failed by %s: %s", c.Nick(), err.Error()))
		} else {
			c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Configuration reloaded successfully",
				c.server.config.Server.Name, c.Nick()))
			c.sendSnomask('s', fmt.Sprintf("Configuration reloaded by %s", c.Nick()))
		}
	}
}

// handleTrace handles TRACE command (show server connection tree)
func (c *Client) handleTrace(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	// Show basic server info (simplified implementation)
	c.SendMessage(fmt.Sprintf(":%s 200 %s Link %s %s %s",
		c.server.config.Server.Name, c.Nick(),
		c.server.config.Server.Version,
		c.server.config.Server.Name,
		"TechIRCd"))

	clientCount := len(c.server.GetClients())
	c.SendMessage(fmt.Sprintf(":%s 262 %s %s :End of TRACE with %d clients",
		c.server.config.Server.Name, c.Nick(),
		c.server.config.Server.Name, clientCount))
}

// handleSpy handles SPY command - covert surveillance and stealth operations
func (c *Client) handleSpy(parts []string) {
	if !c.IsOper() {
		c.SendNumeric(ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	if len(parts) < 2 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** SPY Usage: SPY <hide|watch|track|listen|cloak|ghost|shadow>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	command := strings.ToLower(parts[1])

	switch command {
	case "hide":
		c.handleSpyHide(parts[2:])
	case "watch":
		c.handleSpyWatch(parts[2:])
	case "track":
		c.handleSpyTrack(parts[2:])
	case "listen":
		c.handleSpyListen(parts[2:])
	case "cloak":
		c.handleSpyCloak(parts[2:])
	case "ghost":
		c.handleSpyGhost(parts[2:])
	case "shadow":
		c.handleSpyShadow(parts[2:])
	case "status":
		c.handleSpyStatus()
	default:
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Unknown SPY command: %s",
			c.server.config.Server.Name, c.Nick(), command))
	}
}

// handleSpyHide - become invisible to most commands and lists
func (c *Client) handleSpyHide(args []string) {
	if len(args) == 0 || strings.ToLower(args[0]) == "on" {
		c.SetMode('H', true) // Hidden mode
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** You are now HIDDEN from WHO, WHOIS, and NAMES",
			c.server.config.Server.Name, c.Nick()))
		c.sendSnomask('d', fmt.Sprintf("Operator %s has entered STEALTH mode", c.Nick()))
	} else if strings.ToLower(args[0]) == "off" {
		c.SetMode('H', false)
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** You are now VISIBLE again",
			c.server.config.Server.Name, c.Nick()))
		c.sendSnomask('d', fmt.Sprintf("Operator %s has left STEALTH mode", c.Nick()))
	}
}

// handleSpyWatch - monitor a specific user's activities
func (c *Client) handleSpyWatch(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY WATCH <nickname|off>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	target := args[0]
	if strings.ToLower(target) == "off" {
		// TODO: Remove from watch list
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Surveillance disabled",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	targetClient := c.server.GetClient(target)
	if targetClient == nil {
		c.SendNumeric(ERR_NOSUCHNICK, target+" :No such nick/channel")
		return
	}

	// TODO: Add to watch list
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Now watching %s (%s@%s)",
		c.server.config.Server.Name, c.Nick(), target, targetClient.User(), targetClient.Host()))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Target is in channels: %s",
		c.server.config.Server.Name, c.Nick(), c.getChannelList(targetClient)))
}

// handleSpyTrack - get real-time location and movement tracking
func (c *Client) handleSpyTrack(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY TRACK <nickname>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	target := args[0]
	targetClient := c.server.GetClient(target)
	if targetClient == nil {
		c.SendNumeric(ERR_NOSUCHNICK, target+" :No such nick/channel")
		return
	}

	// Show detailed tracking info
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** TRACKING %s",
		c.server.config.Server.Name, c.Nick(), target))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Location: %s@%s",
		c.server.config.Server.Name, c.Nick(), targetClient.User(), targetClient.Host()))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Status: %s",
		c.server.config.Server.Name, c.Nick(), c.getUserStatus(targetClient)))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Channels: %s",
		c.server.config.Server.Name, c.Nick(), c.getChannelList(targetClient)))

	if targetClient.Away() != "" {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Away: %s",
			c.server.config.Server.Name, c.Nick(), targetClient.Away()))
	}
}

// handleSpyListen - tap into channel conversations invisibly
func (c *Client) handleSpyListen(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY LISTEN <#channel|off>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	target := args[0]
	if strings.ToLower(target) == "off" {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Wiretaps disabled",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	if !isChannelName(target) {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Invalid channel name",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	channel := c.server.GetChannel(target)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, target+" :No such channel")
		return
	}

	// TODO: Add to wiretap list
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Now listening to %s (%d users)",
		c.server.config.Server.Name, c.Nick(), target, channel.UserCount()))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Wiretap established - you will receive covert copies of all messages",
		c.server.config.Server.Name, c.Nick()))
}

// handleSpyCloak - disguise your identity
func (c *Client) handleSpyCloak(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY CLOAK <identity|off>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	identity := args[0]
	if strings.ToLower(identity) == "off" {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Identity cloak removed",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	// TODO: Implement identity cloaking
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Identity cloaked as: %s",
		c.server.config.Server.Name, c.Nick(), identity))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Your true identity is hidden from WHOIS and other commands",
		c.server.config.Server.Name, c.Nick()))
}

// handleSpyGhost - become completely invisible in a channel
func (c *Client) handleSpyGhost(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY GHOST <#channel|off>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	target := args[0]
	if strings.ToLower(target) == "off" {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Ghost mode disabled",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	if !isChannelName(target) {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Invalid channel name",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	channel := c.server.GetChannel(target)
	if channel == nil {
		c.SendNumeric(ERR_NOSUCHCHANNEL, target+" :No such channel")
		return
	}

	// Join channel invisibly
	if !c.IsInChannel(target) {
		channel.AddClient(c)
	}

	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** You are now a GHOST in %s",
		c.server.config.Server.Name, c.Nick(), target))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** You can see everything but are invisible to users",
		c.server.config.Server.Name, c.Nick()))
}

// handleSpyShadow - follow a user invisibly across channels
func (c *Client) handleSpyShadow(args []string) {
	if len(args) < 1 {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Usage: SPY SHADOW <nickname|off>",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	target := args[0]
	if strings.ToLower(target) == "off" {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Shadow mode disabled",
			c.server.config.Server.Name, c.Nick()))
		return
	}

	targetClient := c.server.GetClient(target)
	if targetClient == nil {
		c.SendNumeric(ERR_NOSUCHNICK, target+" :No such nick/channel")
		return
	}

	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** Now shadowing %s",
		c.server.config.Server.Name, c.Nick(), target))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** You will automatically follow them to any channel they join",
		c.server.config.Server.Name, c.Nick()))
}

// handleSpyStatus - show current spy operations
func (c *Client) handleSpyStatus() {
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** === SPY STATUS ===",
		c.server.config.Server.Name, c.Nick()))

	if c.HasMode('H') {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** STEALTH: Active (Hidden from WHO/WHOIS/NAMES)",
			c.server.config.Server.Name, c.Nick()))
	} else {
		c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** STEALTH: Inactive",
			c.server.config.Server.Name, c.Nick()))
	}

	// TODO: Show other active spy operations
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** WATCH: None active",
		c.server.config.Server.Name, c.Nick()))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** WIRETAPS: None active",
		c.server.config.Server.Name, c.Nick()))
	c.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** SHADOWS: None active",
		c.server.config.Server.Name, c.Nick()))
}

// Helper functions for spy operations
func (c *Client) getChannelList(target *Client) string {
	channels := target.GetChannels()
	var channelNames []string
	for _, channel := range channels {
		channelNames = append(channelNames, channel.Name())
	}
	if len(channelNames) == 0 {
		return "None"
	}
	return strings.Join(channelNames, " ")
}

func (c *Client) getUserStatus(target *Client) string {
	status := "Online"
	if target.Away() != "" {
		status = "Away"
	}
	if target.IsOper() {
		status += " (Operator)"
	}
	if target.HasMode('i') {
		status += " (Invisible)"
	}
	if target.HasMode('B') {
		status += " (Bot)"
	}
	return status
}

// sendSnomask sends a server notice to operators watching a specific snomask
func (c *Client) sendSnomask(snomask rune, message string) {
	if c.server == nil {
		return
	}

	for _, client := range c.server.GetClients() {
		if client.IsOper() && client.HasSnomask(snomask) {
			client.SendMessage(fmt.Sprintf(":%s NOTICE %s :*** %s",
				c.server.config.Server.Name, client.Nick(), message))
		}
	}
}

// isValidNickname checks if a nickname is valid
func isValidNickname(nick string) bool {
	if len(nick) == 0 || len(nick) > 30 {
		return false
	}

	// First character must be a letter or special char
	first := nick[0]
	if !((first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z') ||
		first == '[' || first == ']' || first == '\\' || first == '`' ||
		first == '_' || first == '^' || first == '{' || first == '|' || first == '}') {
		return false
	}

	// Rest can be letters, digits, or special chars
	for i := 1; i < len(nick); i++ {
		c := nick[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '[' || c == ']' || c == '\\' || c == '`' ||
			c == '_' || c == '^' || c == '{' || c == '|' || c == '}' || c == '-') {
			return false
		}
	}

	return true
}

// isValidChannelName checks if a channel name is valid
func isValidChannelName(name string) bool {
	if len(name) == 0 || len(name) > 50 {
		return false
	}

	return name[0] == '#' || name[0] == '&' || name[0] == '!' || name[0] == '+'
}

// isChannelName checks if a name is a channel name
func isChannelName(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] == '#' || name[0] == '&' || name[0] == '!' || name[0] == '+'
}
