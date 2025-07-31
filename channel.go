package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Channel struct {
	name       string
	topic      string
	topicBy    string
	topicTime  time.Time
	clients    map[string]*Client
	operators  map[string]*Client
	halfops    map[string]*Client
	voices     map[string]*Client
	owners     map[string]*Client
	modes      map[rune]bool
	key        string
	limit      int
	banList    []string
	quietList  []string
	exceptList []string
	inviteList []string
	created    time.Time
	mu         sync.RWMutex
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:       name,
		clients:    make(map[string]*Client),
		operators:  make(map[string]*Client),
		halfops:    make(map[string]*Client),
		voices:     make(map[string]*Client),
		owners:     make(map[string]*Client),
		modes:      make(map[rune]bool),
		banList:    make([]string, 0),
		quietList:  make([]string, 0),
		exceptList: make([]string, 0),
		inviteList: make([]string, 0),
		created:    time.Now(),
	}
}

func (ch *Channel) Name() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.name
}

func (ch *Channel) Topic() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.topic
}

func (ch *Channel) TopicBy() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.topicBy
}

func (ch *Channel) TopicTime() time.Time {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.topicTime
}

func (ch *Channel) SetTopic(topic, by string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.topic = topic
	ch.topicBy = by
	ch.topicTime = time.Now()
}

func (ch *Channel) AddClient(client *Client) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.clients[strings.ToLower(client.Nick())] = client
	client.AddChannel(ch)

	// First user becomes operator (not owner - owner is for special designation)
	if len(ch.clients) == 1 {
		ch.operators[strings.ToLower(client.Nick())] = client
	}
}

func (ch *Channel) RemoveClient(client *Client) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	nick := strings.ToLower(client.Nick())
	delete(ch.clients, nick)
	delete(ch.operators, nick)
	delete(ch.halfops, nick)
	delete(ch.voices, nick)
	delete(ch.owners, nick)
	client.RemoveChannel(ch.name)
}

func (ch *Channel) HasClient(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	_, exists := ch.clients[strings.ToLower(client.Nick())]
	return exists
}

func (ch *Channel) IsOperator(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	_, exists := ch.operators[strings.ToLower(client.Nick())]
	return exists
}

func (ch *Channel) IsVoice(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	_, exists := ch.voices[strings.ToLower(client.Nick())]
	return exists
}

func (ch *Channel) IsHalfop(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	_, exists := ch.halfops[strings.ToLower(client.Nick())]
	return exists
}

func (ch *Channel) IsOwner(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	_, exists := ch.owners[strings.ToLower(client.Nick())]
	return exists
}

func (ch *Channel) IsQuieted(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.isQuietedUnsafe(client)
}

func (ch *Channel) isQuietedUnsafe(client *Client) bool {
	nick := strings.ToLower(client.Nick())
	hostmask := fmt.Sprintf("%s!%s@%s", client.Nick(), client.user, client.host)

	for _, quiet := range ch.quietList {
		if matched, _ := filepath.Match(quiet, nick); matched {
			return true
		}
		if matched, _ := filepath.Match(quiet, hostmask); matched {
			return true
		}
	}
	return false
}

func (ch *Channel) SetOperator(client *Client, isOp bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	nick := strings.ToLower(client.Nick())
	if isOp {
		ch.operators[nick] = client
	} else {
		delete(ch.operators, nick)
	}
}

func (ch *Channel) SetVoice(client *Client, hasVoice bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	nick := strings.ToLower(client.Nick())
	if hasVoice {
		ch.voices[nick] = client
	} else {
		delete(ch.voices, nick)
	}
}

func (ch *Channel) SetHalfop(client *Client, isHalfop bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	nick := strings.ToLower(client.Nick())
	if isHalfop {
		ch.halfops[nick] = client
	} else {
		delete(ch.halfops, nick)
	}
}

func (ch *Channel) SetOwner(client *Client, isOwner bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	nick := strings.ToLower(client.Nick())
	if isOwner {
		ch.owners[nick] = client
	} else {
		delete(ch.owners, nick)
	}
}

func (ch *Channel) GetClients() []*Client {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	clients := make([]*Client, 0, len(ch.clients))
	for _, client := range ch.clients {
		clients = append(clients, client)
	}
	return clients
}

func (ch *Channel) GetClientCount() int {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.clients)
}

func (ch *Channel) UserCount() int {
	return ch.GetClientCount()
}

func (ch *Channel) Broadcast(message string, exclude *Client) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	for _, client := range ch.clients {
		if exclude != nil && client.Nick() == exclude.Nick() {
			continue
		}
		client.SendMessage(message)
	}
}

func (ch *Channel) BroadcastFrom(source, message string, exclude *Client) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	for _, client := range ch.clients {
		if exclude != nil && client.Nick() == exclude.Nick() {
			continue
		}
		client.SendFrom(source, message)
	}
}

func (ch *Channel) HasMode(mode rune) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.modes[mode]
}

func (ch *Channel) SetMode(mode rune, set bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if set {
		ch.modes[mode] = true
	} else {
		delete(ch.modes, mode)
	}
}

func (ch *Channel) GetModes() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	var modes []rune
	for mode := range ch.modes {
		modes = append(modes, mode)
	}

	if len(modes) == 0 {
		return ""
	}

	return "+" + string(modes)
}

func (ch *Channel) CanSendMessage(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	// Check if user is quieted first
	if ch.isQuietedUnsafe(client) {
		// Only owners, operators, and halfops can speak when quieted
		nick := strings.ToLower(client.Nick())
		_, isOwner := ch.owners[nick]
		_, isOp := ch.operators[nick]
		_, isHalfop := ch.halfops[nick]

		if !isOwner && !isOp && !isHalfop {
			return false
		}
	}

	// If channel is not moderated, anyone in the channel can send
	if !ch.modes['m'] {
		return true
	}

	// In moderated channels, only owners, operators, halfops and voiced users can send messages
	nick := strings.ToLower(client.Nick())
	_, isOwner := ch.owners[nick]
	_, isOp := ch.operators[nick]
	_, isHalfop := ch.halfops[nick]
	_, hasVoice := ch.voices[nick]

	return isOwner || isOp || isHalfop || hasVoice
}

func (ch *Channel) Key() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.key
}

func (ch *Channel) SetKey(key string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.key = key
}

func (ch *Channel) Limit() int {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.limit
}

func (ch *Channel) SetLimit(limit int) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.limit = limit
}

func (ch *Channel) GetNamesReply() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	var names []string
	for _, client := range ch.clients {
		prefix := ""
		if ch.IsOwner(client) {
			prefix = "~"
		} else if ch.IsOperator(client) {
			prefix = "@"
		} else if ch.IsHalfop(client) {
			prefix = "%"
		} else if ch.IsVoice(client) {
			prefix = "+"
		}
		names = append(names, prefix+client.Nick())
	}

	return strings.Join(names, " ")
}

func (ch *Channel) CanSpeak(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	// If channel is not moderated, anyone can speak
	if !ch.modes['m'] {
		return true
	}

	// Operators and voiced users can always speak
	return ch.IsOperator(client) || ch.IsVoice(client)
}

func (ch *Channel) CanJoin(client *Client, key string) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	// Check if invite-only
	if ch.modes['i'] {
		// Check invite list
		for _, mask := range ch.inviteList {
			if ch.matchesMask(client.Prefix(), mask) {
				return true
			}
		}
		return false
	}

	// Check key
	if ch.modes['k'] && ch.key != key {
		return false
	}

	// Check limit
	if ch.modes['l'] && len(ch.clients) >= ch.limit {
		return false
	}

	// Check ban list
	for _, mask := range ch.banList {
		if ch.matchesMask(client.Prefix(), mask) {
			// Check exception list
			for _, exceptMask := range ch.exceptList {
				if ch.matchesMask(client.Prefix(), exceptMask) {
					return true
				}
			}
			return false
		}
	}

	return true
}

func (ch *Channel) matchesMask(target, mask string) bool {
	// Simple mask matching - should be enhanced for production
	return strings.Contains(strings.ToLower(target), strings.ToLower(mask))
}

func (ch *Channel) AddBan(mask string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.banList = append(ch.banList, mask)
}

func (ch *Channel) RemoveBan(mask string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for i, ban := range ch.banList {
		if ban == mask {
			ch.banList = append(ch.banList[:i], ch.banList[i+1:]...)
			break
		}
	}
}

func (ch *Channel) GetBans() []string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	bans := make([]string, len(ch.banList))
	copy(bans, ch.banList)
	return bans
}

func (ch *Channel) Created() time.Time {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.created
}

// IsBanned checks if a client matches any ban mask in the channel
func (ch *Channel) IsBanned(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	hostmask := fmt.Sprintf("%s!%s@%s", client.Nick(), client.User(), client.Host())
	
	for _, ban := range ch.banList {
		if matchWildcard(ban, hostmask) {
			return true
		}
	}
	return false
}

// IsInvited checks if a client is on the invite list for the channel
func (ch *Channel) IsInvited(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	
	hostmask := fmt.Sprintf("%s!%s@%s", client.Nick(), client.User(), client.Host())
	
	for _, invite := range ch.inviteList {
		if matchWildcard(invite, hostmask) {
			return true
		}
	}
	return false
}

// matchWildcard checks if a pattern with wildcards (* and ?) matches a string
func matchWildcard(pattern, str string) bool {
	matched, _ := filepath.Match(strings.ToLower(pattern), strings.ToLower(str))
	return matched
}
