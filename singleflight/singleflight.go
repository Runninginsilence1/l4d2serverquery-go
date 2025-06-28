package singleflight

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

var _group = &singleflight.Group{}
var once sync.Once

func Sf() *singleflight.Group {
	once.Do(func() {
		_group = &singleflight.Group{}
	})
	return _group
}
