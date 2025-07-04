package defs

type NoteSendMsg = *Note
type TopicSendMsg = *Topic
type CurrentViewMsg = string

type InitMsg struct{}

/*
Daily Task
*/
type TaskSucMsg struct{}

type DeleteTaskMsg struct{}

type SwitchContextMsg struct{}
