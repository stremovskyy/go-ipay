package http

import "time"

type Options struct {
	Timeout         time.Duration
	KeepAlive       time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
	IsDebug         bool
}

func DefaultOptions() *Options {
	return &Options{
		Timeout:         5 * time.Second,
		KeepAlive:       30 * time.Second,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		IsDebug:         false,
	}
}
