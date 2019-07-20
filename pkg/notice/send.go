package notice

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuebing1110/notify-inspect/pkg/notice/models"
	"sync"
)

var (
	multiSender []Sender
	mutex       sync.Mutex
)

func init() {
	multiSender = make([]Sender, 0)
}

type Sender interface {
	Name() string
	Send(ctx context.Context, n *models.Notice) error
}

func SenderRegister(s Sender) error {
	mutex.Lock()
	defer mutex.Unlock()

	multiSender = append(multiSender, s)
	return nil
}

func Send(ctx context.Context, n *models.Notice) error {
	errs := make(map[string]error)
	for _, s := range multiSender {
		err := s.Send(ctx, n)
		if err != nil {
			errs[s.Name()] = err
		}
	}

	if len(errs) > 0 {
		return errors.Wrap(fmt.Errorf("%+v", errs), "send the message to some channel failed")
	}

	return nil
}
