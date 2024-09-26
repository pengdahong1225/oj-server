package snowflake

import (
	"sync"
	"time"
)

// Worker 工作节点
type Worker struct {
	mtx             sync.Mutex
	timestamp       int64 // 生成id的时间戳，从startTimestamp到现在的差值
	timestampMax    int64
	timestampOffset uint8
	machineId       int64
	machineIdOffset uint8
	seq             int64
	seqMax          int64
	seqOffset       uint8
}

func (w *Worker) next() (int64, error) {
	// 先校验参数
	if w.timestamp < 0 || w.timestamp > w.timestampMax {
		return 0, TimestampIllegal
	}
	if w.seq < 0 || w.seq > w.seqMax {
		return 0, SeqIllegal
	}

	// 生成id
	var id int64
	id |= w.timestamp << w.timestampOffset
	id |= w.machineId << w.machineIdOffset
	id |= w.seq << w.seqOffset

	return id, nil
}

// GenerateId 生成雪花id
func (w *Worker) GenerateId() (int64, error) {
	// 加锁
	w.mtx.Lock()
	defer w.mtx.Unlock()

	// 生成id
	id, err := w.next()
	if err != nil {
		return 0, err
	}

	// 更新参数：序列号用完了，就增加时间戳并重置序列号
	if w.seq == w.seqMax {
		w.seq = 0
		w.timestamp++
	} else {
		w.seq++
	}

	return id, err
}

// NewWorker 创建一个雪花算法的工作节点
func NewWorker(c SnowFlakeConfig, machineId int64) (*Worker, error) {
	// 检查配置
	sumBits := c.TimestampBits + c.MachineIdBits + c.SeqBits
	if sumBits != 63 {
		return nil, BitsSumError
	}

	// 检查机器码
	var machineMax int64 = (1 << c.MachineIdBits) - 1
	if machineId < 0 || machineId > machineMax {
		return nil, MachineIdIllegal
	}

	return &Worker{
		mtx:             sync.Mutex{},
		timestamp:       time.Now().UnixMilli() - c.StartTimestamp,
		timestampMax:    (1 << c.TimestampBits) - 1,
		timestampOffset: c.SeqBits + c.MachineIdBits,
		machineId:       machineId,
		machineIdOffset: c.SeqBits,
		seq:             0,
		seqMax:          (1 << c.SeqBits) - 1,
		seqOffset:       0,
	}, nil
}
