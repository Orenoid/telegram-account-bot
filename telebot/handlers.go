package telebot

import (
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type HandlersHub struct {
}

func NewHandlerHub() *HandlersHub {
	return &HandlersHub{}
}

func (hub *HandlersHub) HandleStartCommand(ctx telebot.Context) error {
	// TODO 单元测试
	err := ctx.Send("hello")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
