package wss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type subMessage struct {
	Sub []string `json:"sub"`
}

type MsgBody struct {
	To       string
	Response Response
}

type Response struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func (m *MsgBody) BodyHash() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%v", m.Response)))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MsgBody) GetBody() []byte {
	re := m.Response
	data, _ := json.Marshal(re)
	return data
}
