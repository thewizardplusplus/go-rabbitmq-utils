package rabbitmqutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_startModeHolder_getStartMode(test *testing.T) {
	holder := &startModeHolder{mode: 23}
	receivedStartMode := holder.getStartMode()

	assert.Equal(test, startMode(23), receivedStartMode)
}
