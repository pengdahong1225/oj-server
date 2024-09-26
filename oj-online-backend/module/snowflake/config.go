package snowflake

type SnowFlakeConfig struct {
	// 起始时间戳，单位为毫秒，默认为0，即1970-01-01 00:00:00.000
	StartTimestamp int64
	// 时间戳位数，默认为41
	TimestampBits uint8
	// 机器码位数，默认为10
	MachineIdBits uint8
	// 序列号位数，默认为12
	SeqBits uint8
}

// defaultConfig 雪花算法默认配置
const (
	startTimestamp int64 = 1727336390 // 2024-09-26 15:39:50
	timestampBits  uint8 = 41
	machineIdBits  uint8 = 10
	seqBits        uint8 = 12
)

// NewDefaultConfig 创建一个雪花算法默认配置
func NewDefaultConfig() SnowFlakeConfig {
	return SnowFlakeConfig{
		StartTimestamp: startTimestamp,
		TimestampBits:  timestampBits,
		MachineIdBits:  machineIdBits,
		SeqBits:        seqBits,
	}
}
