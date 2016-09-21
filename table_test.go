package main

import "testing"

// TestNewSQLTable tests SQLTable creation
func TestNewSQLTable(t *testing.T) {
	finfo := DBTblInfo{"id", "bigint", "bigint(64) unsigned", "NO", "PRI", "auto_increment"}
	fieldInfo := []*DBTblInfo{&finfo}
	_, err := NewSQLTable("TEST", "TEST", fieldInfo)
	if err != nil {
		t.Error(
			"For", fieldInfo,
			"Expected", "",
			"Got", err,
		)
	}
}
