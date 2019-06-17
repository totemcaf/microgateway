package usecase

import (
	"errors"
	"net/url"
	"testing"

	"github.com/totemcaf/microgateway.git/app/entity"

	"github.com/stretchr/testify/assert"
)

func TestReproducer_ReceiveAndSend_NoTarget(t *testing.T) {
	// GIVEN a reproducer with a target
	target, _ := entity.NewTarget("1", "/data1/", "")
	repo := &messageRepositoryMock{}
	sender := &messageSenderMock{}

	reproducer := NewReproducer(target, repo, sender)

	// WHEN it receives an unhandled message
	m := entity.Message{Path: "/other"}
	err := reproducer.ReceiveAndSend(&m)

	// THEN it reports the error and does not store anything
	assert.Error(t, err)
	assert.IsType(t, &NoTargetError{}, err)
	assert.Nil(t, repo.mh, "it was not called")
	assert.Nil(t, sender.target, "it was not called")
}

func TestReproducer_ReceiveAndSend_StoreAndSend(t *testing.T) {
	// GIVEN A reproducer
	target, _ := entity.NewTarget("1", "/data1/", "")
	repo := &messageRepositoryMock{}
	sender := &messageSenderMock{}

	reproducer := NewReproducer(target, repo, sender)

	// WHEN received a new message
	m := entity.Message{Path: "/data1/"}
	err := reproducer.ReceiveAndSend(&m)

	// THEN it is stored and send
	assert.NoError(t, err)
	assert.NotNil(t, repo.mh, "it was called")
	assert.NotNil(t, sender.target, "it was called")
}

func TestReproducer_ReceiveAndSend_GeneratesID(t *testing.T) {
	// GIVEN A reproducer
	target, _ := entity.NewTarget("1", "/data1/", "")
	repo := &messageRepositoryMock{}
	sender := &messageSenderMock{returnCode: 200}

	reproducer := NewReproducer(target, repo, sender)

	// WHEN received a new message and it is send ok
	m := entity.Message{Path: "/data1/"}
	err := reproducer.ReceiveAndSend(&m)

	// THEN it is stored with the success result
	assert.NoError(t, err)
	assert.NotEmpty(t, repo.mh.ID, "an ID was provided")
}

func TestReproducer_ReceiveAndSend_PersistSendSuccess(t *testing.T) {
	// GIVEN A reproducer
	target, _ := entity.NewTarget("1", "/data1/", "")
	repo := &messageRepositoryMock{}
	sender := &messageSenderMock{returnCode: 200}

	reproducer := NewReproducer(target, repo, sender)

	// WHEN received a new message and it is send ok
	m := entity.Message{Path: "/data1/"}
	err := reproducer.ReceiveAndSend(&m)

	// THEN it is stored with the success result
	assert.NoError(t, err)
	assert.Len(t, repo.mh.Deliveries, 1, "a history was added")
	assert.Equal(t, 200, repo.mh.Deliveries[0].TargetResponseCode, "the send result was persisted")
	assert.Equal(t, "", repo.mh.Deliveries[0].TargetResponse, "the send result was persisted")
}

// fails to send
func TestReproducer_ReceiveAndSend_FailToSend(t *testing.T) {
	// GIVEN A reproducer
	target, _ := entity.NewTarget("1", "/data1/", "")
	repo := &messageRepositoryMock{}
	sender := &messageSenderMock{returnCode: TimeOutErrorCode, errToReturn: errors.New("cannot send")}

	reproducer := NewReproducer(target, repo, sender)

	// WHEN received a new message
	m := entity.Message{Path: "/data1/"}
	err := reproducer.ReceiveAndSend(&m)

	// THEN it is stored and send
	assert.NoError(t, err)
	assert.NotNil(t, repo.mh, "it was called")
	assert.NotNil(t, sender.target, "it was called")
	assert.Len(t, repo.mh.Deliveries, 1, "a history was added")
	assert.Equal(t, TimeOutErrorCode, repo.mh.Deliveries[0].TargetResponseCode, "the send result was persisted")
	assert.Equal(t, "cannot send", repo.mh.Deliveries[0].TargetResponse, "the send result was persisted")
}

// resend

type messageRepositoryMock struct {
	errToReturn error
	mh          *entity.MessageHistory
}

func (r *messageRepositoryMock) Store(mh *entity.MessageHistory) error {
	r.mh = mh
	return r.errToReturn
}

// MessageSender helper to send a message to target
type messageSenderMock struct {
	returnCode  int
	errToReturn error
	target      *url.URL
	m           *entity.Message
}

func (s *messageSenderMock) Send(target *url.URL, m *entity.Message) (int, error) {
	s.target = target
	s.m = m

	return s.returnCode, s.errToReturn
}
