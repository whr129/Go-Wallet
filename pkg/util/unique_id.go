package util

import (
	"fmt"

	"github.com/sony/sonyflake/v2"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings

	var err error
	sf, err = sonyflake.New(st)
	if err != nil {
		panic(err)
	}
}

func GenerateID() (int64, error) {
	var err error
	id, err := sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("failed to generate ID: %w", err)
	}

	return id, nil
}
