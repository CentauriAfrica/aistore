// Command-line mounting utility for aisfs.
/*
 * Copyright (c) 2019, NVIDIA CORPORATION. All rights reserved.
 */
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/NVIDIA/aistore/cmn"
	"github.com/NVIDIA/aistore/fuse/fs"
)

const configDirName = fs.Name

var defaultConfig = Config{
	Cluster: ClusterConfig{
		URL: "http://127.0.0.1:8080",
	},
	Timeout: TimeoutConfig{
		TCPTimeoutStr:  "60s",
		TCPTimeout:     60 * time.Second,
		HTTPTimeoutStr: "300s",
		HTTPTimeout:    300 * time.Second,
	},
	Periodic: PeriodicConfig{
		SyncIntervalStr: "20m",
		SyncInterval:    20 * time.Minute,
	},
	Log: LogConfig{
		ErrorLogFile: "",
		DebugLogFile: "",
	},
	IO: IOConfig{
		// Determines the size of chunks that we write with append. The only exception
		// when we write less is Flush (end-of-file).
		WriteBufSize: cmn.MiB,
	},
}

type (
	Config struct {
		Cluster     ClusterConfig  `json:"cluster"`
		Timeout     TimeoutConfig  `json:"timeout"`
		Periodic    PeriodicConfig `json:"periodic"`
		Log         LogConfig      `json:"log"`
		IO          IOConfig       `json:"io"`
		MemoryLimit string         `json:"memory_limit"`
	}

	ClusterConfig struct {
		URL string `json:"url"`
	}

	TimeoutConfig struct {
		TCPTimeoutStr  string        `json:"tcp_timeout"`
		TCPTimeout     time.Duration `json:"-"`
		HTTPTimeoutStr string        `json:"http_timeout"`
		HTTPTimeout    time.Duration `json:"-"`
	}

	PeriodicConfig struct {
		SyncIntervalStr string        `json:"sync_interval"`
		SyncInterval    time.Duration `json:"-"`
	}

	LogConfig struct {
		ErrorLogFile string `json:"error_log_file"`
		DebugLogFile string `json:"debug_log_file"`
	}

	IOConfig struct {
		WriteBufSize int64 `json:"write_buf_size"`
	}
)

func (c *Config) validate() (err error) {
	if c.Timeout.TCPTimeout, err = time.ParseDuration(c.Timeout.TCPTimeoutStr); err != nil {
		return fmt.Errorf("invalid timeout.tcp_timeout format %q: %v", c.Timeout.TCPTimeoutStr, err)
	}
	if c.Timeout.HTTPTimeout, err = time.ParseDuration(c.Timeout.HTTPTimeoutStr); err != nil {
		return fmt.Errorf("invalid timeout.http_timeout format %q: %v", c.Timeout.HTTPTimeoutStr, err)
	}
	if c.Periodic.SyncInterval, err = time.ParseDuration(c.Periodic.SyncIntervalStr); err != nil {
		return fmt.Errorf("invalid periodic.sync_interval format %q: %v", c.Periodic.SyncInterval, err)
	}
	if c.Log.ErrorLogFile != "" && !filepath.IsAbs(c.Log.ErrorLogFile) {
		return fmt.Errorf("invalid log.error_log_file format %q: path needs to be absolute", c.Log.ErrorLogFile)
	}
	if c.Log.DebugLogFile != "" && !filepath.IsAbs(c.Log.DebugLogFile) {
		return fmt.Errorf("invalid log.debug_log_file format %q: path needs to be absolute", c.Log.DebugLogFile)
	}
	if c.IO.WriteBufSize < 0 {
		return fmt.Errorf("invalid io.write_buf_size value: %d: expected non-negative value", c.IO.WriteBufSize)
	}
	if v, err := cmn.S2B(c.MemoryLimit); err != nil {
		return fmt.Errorf("invalid memory_limit value: %q: %v", c.MemoryLimit, err)
	} else if v < 0 {
		return fmt.Errorf("invalid memory_limit value: %q: expected non-negative value", c.MemoryLimit)
	}
	return nil
}

func loadConfig(bucket string) (cfg *Config, err error) {
	cfg = &Config{}
	configFileName := bucket + "_mount.json"
	if err = cmn.LoadAppConfig(configDirName, configFileName, &cfg); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}

		cfg = &defaultConfig
		err = cmn.SaveAppConfig(configDirName, configFileName, cfg)
		if err != nil {
			err = fmt.Errorf("failed to generate config file: %v", err)
		}
		return
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return
}
