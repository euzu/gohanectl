package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepare(t *testing.T) {

	deviceTopic := DeviceTopic{
		Topic: "abcd",
	}
	deviceTopic.Prepare()
	m := deviceTopic.Matches("abcd")
	assert.True(t, m)

	deviceTopic.Topic = "abcd/#"
	deviceTopic.Prepare()
	m = deviceTopic.Matches("abcd")
	assert.True(t, m)

	deviceTopic.Topic = "abcd/"
	deviceTopic.Prepare()
	m = deviceTopic.Matches("abcd")
	assert.True(t, m)

	deviceTopic.Topic = "abcd/-"
	deviceTopic.Prepare()
	m = deviceTopic.Matches("abcd")
	assert.False(t, m)

	deviceTopic.Topic = "abcde"
	deviceTopic.Prepare()
	m = deviceTopic.Matches("abcd")
	assert.False(t, m)

}
