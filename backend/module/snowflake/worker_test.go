package snowflake

import (
	"testing"
	"time"
)

func TestWorker_GenerateId(t *testing.T) {
	idMap := make(map[int64]int)

	config := NewDefaultConfig()
	worker, _ := NewWorker(config, 1)
	for i := 0; i < 100; i++ {
		id, _ := worker.GenerateId()
		if _, ok := idMap[id]; ok {
			t.Error("id重复")
		} else {
			idMap[id] = 1
		}
		t.Log(id)
		time.Sleep(1)
	}
}
