package parser

import (
	"content_autogen/dto"
	"encoding/json"
	"fmt"
)

type Parser interface {
	Parse(msg []byte) (dto.ContentAutogenEvent, error)
}

type parserImpl struct{}

func NewParser() Parser {
	return &parserImpl{}
}

func (p *parserImpl) Parse(msg []byte) (dto.ContentAutogenEvent, error) {
	var event dto.ContentAutogenEvent

	err := json.Unmarshal(msg, &event)

	if err != nil {
		return dto.ContentAutogenEvent{}, err
	}
	fmt.Println(event)

	return event, err

}
