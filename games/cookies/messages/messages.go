package messages

import (
	"encoding/json"
)

type msgType int8

const (
	ViewPortRequestType      = 1
	ViewPortResponseType     = 2
	UserJoinRequestType      = 3
	UserJoinResponseType     = 4
	CreateCookieRequestType  = 5
	CreateCookieResponseType = 6
	StatsResponseType        = 7
)

type Message interface {
	GetType() msgType
	SetType(msgType)
}

type BaseMessage struct {
	Type msgType         `json:"t"`
	Data json.RawMessage `json:"d,omitempty"`
}

func (m *BaseMessage) GetType() msgType {
	return m.Type
}

func (m *BaseMessage) SetType(t msgType) {
	m.Type = t
}

type ViewPortRequest struct {
	BaseMessage
	X     float32 `json:"X"`
	Y     float32 `json:"Y"`
	XX    float32 `json:"XX"`
	YY    float32 `json:"YY"`
	Angle float32 `json:"R"`
	Turbo bool    `json:"T"`
}

func NewViewPortRequest(X, Y, XX, YY, Angle float32, turbo bool) *ViewPortRequest {
	resp := &ViewPortRequest{X: X, Y: Y, XX: XX, YY: YY, Angle: Angle, Turbo: turbo}
	resp.SetType(ViewPortRequestType)
	return resp
}

type ViewportResponse struct {
	BaseMessage
	Cookies []*CookieInfo `json:"C"`
	Food    []*FoodInfo   `json:"F"`
}

type CookieInfo struct {
	ID    uint64  `json:"ID"`
	Score uint64  `json:"SC"`
	X     float32 `json:"X"`
	Y     float32 `json:"Y"`
}

type FoodInfo struct {
	ID    uint64  `json:"ID"`
	Score uint64  `json:"SC"`
	X     float32 `json:"X"`
	Y     float32 `json:"Y"`
}

type CreateCookieRequest struct {
	BaseMessage
}

func NewCreateCookieRequest() *CreateCookieRequest {
	resp := &CreateCookieRequest{}
	resp.SetType(CreateCookieRequestType)
	return resp
}

type CreateCookieResponse struct {
	BaseMessage
	Data CookieInfo `json:"d"`
}

func NewCreateCookieResponse(ID uint64, sc uint64, X float32, Y float32) *CreateCookieResponse {
	resp := &CreateCookieResponse{
		Data: CookieInfo{
			ID:    ID,
			Score: sc,
			X:     X,
			Y:     Y,
		},
	}
	resp.SetType(CreateCookieResponseType)
	return resp
}

type UserJoinRequest struct {
	BaseMessage
	Username string `json:"UN"`
}

func NewUserJoinRequest(name string) *UserJoinRequest {
	resp := &UserJoinRequest{Username: name}
	resp.SetType(UserJoinRequestType)
	return resp
}

type userJoinResponseData struct {
	Ok       bool     `json:"OK"`
	AltNames []string `json:"AN"`
}

type UserJoinResponse struct {
	BaseMessage
	Data userJoinResponseData `json:"d"`
}

func NewUserJoinResponse(ok bool, altNames []string) *UserJoinResponse {
	resp := &UserJoinResponse{Data: userJoinResponseData{Ok: ok, AltNames: altNames}}
	resp.SetType(UserJoinResponseType)
	return resp
}

type StatsResponseData struct {
	FoodCount    uint64 `json:"FC"`
	CookiesCount uint64 `json:"CC"`
}

type StatsResponse struct {
	BaseMessage
	Data StatsResponseData
}

func NewStatsResponse(food uint64, cookies uint64) *StatsResponse {
	resp := &StatsResponse{Data: StatsResponseData{FoodCount: food, CookiesCount: cookies}}
	resp.SetType(StatsResponseType)
	return resp
}
