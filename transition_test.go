package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetStepTime(t *testing.T) {
	assert.Equal(t, 100*time.Millisecond, getStepTime(10, 15, 500*time.Millisecond))
	assert.Equal(t, 100*time.Millisecond, getStepTime(15, 10, 500*time.Millisecond))
}
