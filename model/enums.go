package model

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleUser
)
