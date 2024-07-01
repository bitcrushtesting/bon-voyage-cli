// Copyright 2024 Bitcrush Testing

package models

const ClientsApiPath = "/api/clients"

type Information struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Version    string `json:"version"`
	ApiVersion string `json:"api_version"`
}

type DeviceBase struct {
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

//-----------------------------------
//
//  User API
//
//-----------------------------------

type UserRegisterCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//-----------------------------------
//
//   Session API
//
//-----------------------------------

type Priority int

const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

type SessionListReply struct {
	SessionId string `json:"session_id"`
	DeviceId  string `json:"device_id"`
	Created   string `json:"created"`
}

type SessionCreatePayload struct {
	Command  string `json:"command"`
	DeviceId string `json:"device_id"`
}

type SessionCreateReply struct {
	SessionId string `json:"session_id"`
}
