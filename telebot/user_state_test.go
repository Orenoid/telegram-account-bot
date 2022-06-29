package telebot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInMemoryUserStateManager(t *testing.T) {
	m := NewInMemoryUserStateManager()
	var _ UserStateManager = m
	assert.IsType(t, &InMemoryUserStateManager{}, m)
}

func TestInMemoryUserStateManager_ClearUserState(t *testing.T) {
	m := NewInMemoryUserStateManager()

	m.cache.Store(int64(7), nil)
	m.cache.Store(int64(8), nil)

	err := m.ClearUserState(7)
	assert.NoError(t, err)
	_, found := m.cache.Load(int64(7))
	assert.False(t, found)
	_, found = m.cache.Load(int64(8))
	assert.True(t, found)
}

func TestInMemoryUserStateManager_GetUserState(t *testing.T) {
	m := NewInMemoryUserStateManager()
	preState := &UserState{}
	m.cache.Store(int64(42), preState)
	m.cache.Store(int64(44), "not *UserState")

	state, exists, err := m.GetUserState(42)
	assert.True(t, exists)
	assert.NoError(t, err)
	assert.True(t, preState == state)

	state, exists, err = m.GetUserState(43)
	assert.NoError(t, err)
	assert.False(t, exists)

	state, exists, err = m.GetUserState(44)
	assert.ErrorContains(t, err, "invalid type of state value: string")
}

func TestInMemoryUserStateManager_SetUserState(t *testing.T) {
	testCases := map[int64]*UserState{
		83: {Type: CreatingBill},
		84: {Type: ""},
	}
	m := NewInMemoryUserStateManager()
	for userID, state := range testCases {
		err := m.SetUserState(userID, state)
		assert.NoError(t, err)
	}
	cacheLen := 0
	m.cache.Range(func(key, value interface{}) bool {
		userID, ok := key.(int64)
		assert.True(t, ok)
		expectedState, found := testCases[userID]
		assert.True(t, found)
		state, ok := value.(*UserState)
		assert.True(t, expectedState == state)
		cacheLen++
		return true
	})
	assert.Equal(t, len(testCases), cacheLen)
}
