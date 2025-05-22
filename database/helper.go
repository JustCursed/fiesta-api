package database

import (
	"fmt"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"strconv"
)

func GetLastChatLogs(size uint64, private bool) []string {
	executeQuery(ch.Query{Body: "SELECT * FROM chat ORDER BY id DESC LIMIT " + strconv.FormatUint(size, 10), Result: proto.Results{}})

	return []string{}
}

func executeQuery(query ch.Query) bool {
	err := pool.Do(ctx, query)

	if err != nil {
		_ = fmt.Errorf("failed to execute query: %v", err)
		return false
	}

	return true
}
