package network

func (connector *Connector) LoadUser(uid uint) {
	connector.sender <- NetworkData{
		"action": "user",
		"uid":    uid,
	}
}

func (connector *Connector) LoadMessages(limit int, channel uint) {
	connector.sender <- NetworkData{
		"action": "fetchmsg",
		"limit":  limit,
		"cid":    channel,
	}
}

func (connector *Connector) LoadGroups() {
	connector.sender <- NetworkData{
		"action": "fetchgroup",
	}
}

func (connector *Connector) LoadChannels() {
	connector.sender <- NetworkData{
		"action": "fetchchannel",
	}
}

func (connector *Connector) SendMessage(text string, channel uint) {
	connector.sender <- NetworkData{
		"action": "message",
		"text":   text,
		"cid":    channel,
	}
}
