package model

type PermType int

const (
	Visitor PermType = iota
	Student
	Faculty
	Admin
)
