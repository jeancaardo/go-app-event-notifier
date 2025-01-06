package internalevents

import "errors"

var (
	ErrEventNotFound          = errors.New("event not found")
	ErrEventNameAlreadyExists = errors.New("event name already exists")
	ErrOnStoreEvent           = errors.New("error on store event")
	ErrOnUpdateEvent          = errors.New("error on update event")
	ErrOnDeleteEvent          = errors.New("error on delete event")
)
