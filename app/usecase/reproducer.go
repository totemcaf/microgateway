package usecase

import (
	"net/url"
	"time"

	"github.com/totemcaf/microgateway.git/app/entity"
)

// TimeOutErrorCode indicates the Message Sender could not send the response due to time out
const TimeOutErrorCode = 601

// MessageRepository a repository for messages
type MessageRepository interface {
	Store(mh *entity.MessageHistory) error
}

// MessageSender helper to send a message to target
type MessageSender interface {
	Send(target *url.URL, m *entity.Message) (int, error) // Error code, error
}

// Reproducer allows to reproduce a message received and sent. That means, a message can be sent twice, or more times
type Reproducer struct {
	target *entity.Target

	messageRepository MessageRepository
	sender            MessageSender
}

// NewReproducer builds a reproducer
func NewReproducer(
	target *entity.Target,
	messageRepository MessageRepository,
	sender MessageSender,
) *Reproducer {
	return &Reproducer{
		target:            target,
		messageRepository: messageRepository,
		sender:            sender,
	}
}

// ReceiveAndSend recevies a message, stores it, and send to target
func (r *Reproducer) ReceiveAndSend(msg *entity.Message) error {
	_, found := r.findTarget(msg)

	if !found {
		return NewNoTargetError(msg.Path)
	}

	mh := entity.NewMessageHistory(r.target.GetIDFor(msg), msg)

	if err := r.messageRepository.Store(mh); err != nil {
		return err
	}

	url, err := r.target.MakeURL(mh.Message)

	if err != nil {
		return err
	}

	returnCode, err := r.sender.Send(url, mh.Message)

	if err == nil {
		mh.Sent(time.Now(), returnCode, "")
	} else {
		mh.Sent(time.Now(), returnCode, err.Error())
	}

	err = r.messageRepository.Store(mh)

	return err
}

func (r *Reproducer) findTarget(m *entity.Message) (*entity.Target, bool) {
	return r.target, r.target.Match(m.Path)
}
