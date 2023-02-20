package snowflake1

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

// Init 初始化雪花算法
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000 //需要毫秒级的时间
	node, err = sf.NewNode(machineID)
	return
}

// GenID 生成有时间顺序的id
func GenID() int64 {
	return node.Generate().Int64()
}
