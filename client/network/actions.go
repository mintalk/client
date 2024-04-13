package network

func (connector *Connector) LoadUser(uid uint) {
	connector.sender <- NetworkData{
		"action": "user",
		"uid":    uid,
	}
}

func (connector *Connector) LoadMessages(limit int, channel uint) {
	connector.sender <- NetworkData{
		"action": "messages",
		"limit":  limit,
		"cid":    channel,
	}
}

func (connector *Connector) LoadGroups() {
	connector.sender <- NetworkData{
		"action": "groups",
	}
}

func (connector *Connector) LoadChannels() {
	connector.sender <- NetworkData{
		"action": "channels",
	}
}

func (connector *Connector) LoadUsers() {
	connector.sender <- NetworkData{
		"action": "users",
	}
}

func (connector *Connector) SendMessage(contents string, channel uint) {
	connector.sender <- NetworkData{
		"action":   "new-message",
		"contents": contents,
		"cid":      channel,
	}
}
