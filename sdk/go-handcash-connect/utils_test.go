package handcash

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatISOTimestamp(t *testing.T) {
	t.Parallel()

	timestamp := currentISOTimestamp()
	extractedTime, err := time.Parse(isoFormat, timestamp)
	assert.NoError(t, err)
	assert.NotNil(t, extractedTime)
	assert.WithinDuration(t, time.Now().UTC(), extractedTime, 1*time.Second)
}
