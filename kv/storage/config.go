package storage

type Config struct {
	//存储位置
	DBPath string

	//
	HeartbeatInterval uint64
	ElectionTimeout   electionTimeoutRange
	BaseTimeInternal  uint64
}

type electionTimeoutRange struct {
	start uint64
	end   uint64
}
