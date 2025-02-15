package token

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Akvicor/util"
	"time"
)

type Model struct {
	ID     int64  `json:"id"`
	Type   Type   `json:"type"`
	Time   int64  `json:"time"`
	Random string `json:"random"`
}

func (t *Model) String() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func NewLoginToken(id int64) string {
	token := &Model{
		ID:     id,
		Time:   time.Now().Unix(),
		Type:   TypeLoginToken,
		Random: util.RandomString(32),
	}
	return token.String()
}

func NewAccessToken(id int64) string {
	token := &Model{
		ID:     id,
		Time:   time.Now().Unix(),
		Type:   TypeAccessToken,
		Random: util.RandomString(32),
	}
	return token.String()
}

func Parse(token string) *Model {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil
	}
	tokenV := new(Model)
	err = json.Unmarshal(data, tokenV)
	if err != nil {
		return nil
	}
	return tokenV
}
