package telebot

import (
	"github.com/pkg/errors"
	"sync"
)

type UserStateType string

const (
	CreatingBill UserStateType = "creatingBill"
)

type UserState struct {
	Type UserStateType `json:"type"`

	// Type 为 CreatingBill 时不为空
	BillCategory *string `json:"bill_category"`
	// Type 为 CreatingBill 时可能有值，若为空则代表未设置名称
	BillName *string `json:"bill_name"`
}

type UserStateManager interface {
	GetUserState(userID int64) (state *UserState, exists bool, err error)
	SetUserState(userID int64, state *UserState) error
	ClearUserState(userID int64) error
}

type InMemoryUserStateManager struct {
	cache *sync.Map
}

func (manager *InMemoryUserStateManager) GetUserState(userID int64) (state *UserState, exists bool, err error) {
	value, found := manager.cache.Load(userID)
	if !found {
		return nil, false, nil
	}
	state, ok := value.(*UserState)
	if !ok {
		return nil, false, errors.Errorf("invalid type of state value: %T", value)
	}
	return state, true, nil
}

func (manager *InMemoryUserStateManager) SetUserState(userID int64, state *UserState) error {
	manager.cache.Store(userID, state)
	return nil
}

func (manager *InMemoryUserStateManager) ClearUserState(userID int64) error {
	manager.cache.Delete(userID)
	return nil
}

func NewInMemoryUserStateManager() *InMemoryUserStateManager {
	manager := &InMemoryUserStateManager{cache: &sync.Map{}}
	return manager
}
