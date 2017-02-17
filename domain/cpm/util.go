package cpm

import (
	"time"
)

func allKeysInMap(ks []int, m map[int]struct{}) bool {
	for _, k := range ks {
		_, found := m[k]
		if !found {
			return false
		}
	}
	return true
}

func maxTime(t1 time.Time, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func minTime(t1 time.Time, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}