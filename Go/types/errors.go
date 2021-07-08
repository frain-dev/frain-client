package types

import "fmt"

type FrainException struct {
	Message string
}

func (e *FrainException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *FrainException) ErrorMessage() string {
	return e.Message
}

func (e *FrainException) ErrorCode() string { return "FrainException" }

type ClientException struct {
	Message string
}

func (e *ClientException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ClientException) ErrorMessage() string {
	return e.Message
}
func (e *ClientException) ErrorCode() string {
	return "ClientException"
}

type ServerException struct {
	Message string
}

func (e *ServerException) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}
func (e *ServerException) ErrorMessage() string {
	return e.Message
}
func (e *ServerException) ErrorCode() string {
	return "ServerException"
}
