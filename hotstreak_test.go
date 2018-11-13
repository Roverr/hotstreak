package hotstreak

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	hotstreak := New(Config{})
	assert.False(t, hotstreak.IsActive())
	assert.False(t, hotstreak.IsHot())
	assert.Equal(t, hotstreak.counter, 0)
	assert.Equal(t, hotstreak.Limit, 20)
	assert.Equal(t, hotstreak.ActiveWait, time.Minute*5)
	assert.Equal(t, hotstreak.HotWait, time.Minute*5)
}

func TestActivation(t *testing.T) {
	t.Run("Should be able to activate", func(t *testing.T) {
		hotstreak := New(Config{})
		hotstreak.Activate()
		assert.True(t, hotstreak.IsActive())
	})

	t.Run("Should be able to end activation after wait time", func(t *testing.T) {
		hotstreak := New(Config{ActiveWait: time.Millisecond * 500})
		hotstreak.Activate()
		assert.True(t, hotstreak.IsActive())
		<-time.After(time.Second)
		assert.False(t, hotstreak.IsActive())
	})

	t.Run("Should not deactivate if any hit has been made in the activation time", func(t *testing.T) {
		hotstreak := New(Config{ActiveWait: time.Millisecond * 500})
		hotstreak.Activate()
		assert.True(t, hotstreak.Hit().IsActive())
		assert.Equal(t, hotstreak.counter, 1)
		<-time.After(time.Millisecond * 600)
		assert.True(t, hotstreak.IsActive())
		assert.Equal(t, hotstreak.counter, 0)
		<-time.After(time.Millisecond * 600)
		assert.False(t, hotstreak.IsActive())
	})
}

func TestHitting(t *testing.T) {
	t.Run("Should be able to reach being hot", func(t *testing.T) {
		hotstreak := New(Config{Limit: 2})
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 1)
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 2)
		assert.True(t, hotstreak.IsHot())
	})

	t.Run("Should not remain hot when deactivated", func(t *testing.T) {
		hotstreak := New(Config{Limit: 2})
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 1)
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 2)
		assert.True(t, hotstreak.IsHot())
		assert.False(t, hotstreak.Deactivate().IsHot())
		assert.Equal(t, hotstreak.counter, 0)
	})

	t.Run("Should not be able to increase counter after getting hot", func(t *testing.T) {
		hotstreak := New(Config{Limit: 2})
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 1)
		hotstreak.Hit()
		assert.Equal(t, hotstreak.counter, 2)
		assert.True(t, hotstreak.IsHot())
		assert.Equal(t, hotstreak.Hit().counter, 2)
	})
}

func TestHot(t *testing.T) {
	t.Run("Should be able to cool down after hot wait time", func(t *testing.T) {
		hotstreak := New(Config{Limit: 2, HotWait: time.Millisecond * 500})
		assert.True(t, hotstreak.Activate().Hit().Hit().IsHot())
		assert.True(t, hotstreak.IsActive())
		<-time.After(time.Second)
		assert.False(t, hotstreak.IsHot())
		assert.True(t, hotstreak.IsActive())
	})
}
