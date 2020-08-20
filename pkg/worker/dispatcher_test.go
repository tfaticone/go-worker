package worker

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewDispatcher(t *testing.T) {
	t.Run("Retrieve object", func(t *testing.T) {
		newDispatcher := NewDispatcher(5)
		assert.Equal(t, 5, newDispatcher.MaxWorkers,
			"[TestNewDispatcher] Max workers not 5 in new dispatcher")
		assert.IsType(t, *new(chan chan JobInterface), newDispatcher.WorkerPool,
			"[TestNewDispatcher] Workerpool is not of type chan chan JobInterface")
	})
}