// Package commands provides the set of CLI commands used to communicate with the AIS cluster.
// This specific file handles the CLI commands that interact with buckets in the cluster
/*
 * Copyright (c) 2019, NVIDIA CORPORATION. All rights reserved.
 */
package commands

import (
	"time"

	"github.com/NVIDIA/aistore/cmn"
	"github.com/urfave/cli"
)

// Top level Command names
const (
	commandBucket = "bucket"
	commandConfig = "config"
	// commandDaeclu - all subcommands are top-level commands
	commandDownload = "download"
	commandDsort    = cmn.DSortNameLowercase
	commandLRU      = "lru"
	commandObject   = "object"
	commandXaction  = "xaction"
)

// Subcommand names
const (
	// Common
	subcommandRename = "rename"
	subcommandEvict  = "evict"
	subcommandStart  = "start"
	subcommandStatus = "status"
	subcommandAbort  = "abort"
	subcommandRemove = "remove"
	subcommandList   = "ls"
	subcommandSet    = "set"
	subcommandGet    = "get"

	// Bucket
	bucketCreate     = "create"
	bucketDestroy    = "destroy"
	bucketNWayMirror = cmn.ActMakeNCopies
	bucketEvict      = subcommandEvict
	bucketSummary    = "summary"
	bucketObjects    = "objects"

	// Bucket props
	commandBucketProps = "props"
	propsList          = subcommandList
	propsReset         = "reset"
	propsSet           = subcommandSet

	// Config
	configGet = subcommandGet
	configSet = subcommandSet

	// Daeclu
	daecluSmap      = cmn.GetWhatSmap
	daecluStats     = cmn.GetWhatStats
	daecluDiskStats = cmn.GetWhatDiskStats
	daecluStatus    = cmn.GetWhatDaemonStatus

	// Download
	downloadStart  = subcommandStart
	downloadStatus = subcommandStatus
	downloadAbort  = subcommandAbort
	downloadRemove = subcommandRemove
	downloadList   = subcommandList

	// dSort
	dsortGen    = "gen"
	dsortStart  = subcommandStart
	dsortStatus = subcommandStatus
	dsortAbort  = subcommandAbort
	dsortRemove = subcommandRemove
	dsortList   = subcommandList

	// Lru
	lruStart  = subcommandStart
	lruStop   = "stop"
	lruStatus = subcommandStatus

	// Object
	objGet      = subcommandGet
	objPut      = "put"
	objDel      = "delete"
	objStat     = "stat"
	objPrefetch = "prefetch"
	objEvict    = subcommandEvict

	// Xaction
	xactStart = cmn.ActXactStart
	xactStop  = cmn.ActXactStop
	xactStats = cmn.ActXactStats
)

// Flag related constants
const (
	aisBucketEnvVar         = "AIS_BUCKET"
	aisBucketProviderEnvVar = "AIS_BUCKET_PROVIDER"

	refreshRateDefault = time.Second
)

// Flags
var (
	// Common
	bucketFlag      = cli.StringFlag{Name: "bucket", Usage: "bucket where the objects are stored, eg. 'imagenet'", EnvVar: aisBucketEnvVar}
	bckProviderFlag = cli.StringFlag{Name: "provider",
		Usage:  "determines which bucket ('local' or 'cloud') should be used. By default, locality is determined automatically",
		EnvVar: aisBucketProviderEnvVar,
	}
	objPropsFlag    = cli.StringFlag{Name: "props", Usage: "properties to return with object names, comma separated", Value: "size,version"}
	prefixFlag      = cli.StringFlag{Name: "prefix", Usage: "prefix for string matching"}
	refreshFlag     = cli.StringFlag{Name: "refresh", Usage: "refresh period", Value: refreshRateDefault.String()}
	regexFlag       = cli.StringFlag{Name: "regex", Usage: "regex pattern for matching"}
	jsonFlag        = cli.BoolFlag{Name: "json,j", Usage: "json input/output"}
	noHeaderFlag    = cli.BoolFlag{Name: "no-headers,H", Usage: "display tables without headers"}
	progressBarFlag = cli.BoolFlag{Name: "progress", Usage: "display progress bar"}

	// Bucket
	jsonspecFlag      = cli.StringFlag{Name: "jsonspec", Usage: "bucket properties in JSON format"}
	markerFlag        = cli.StringFlag{Name: "marker", Usage: "start listing bucket objects starting from the object that follows the marker(alphabetically), ignored in fast mode"}
	objLimitFlag      = cli.StringFlag{Name: "limit", Usage: "limit object count", Value: "0"}
	pageSizeFlag      = cli.StringFlag{Name: "page-size", Usage: "maximum number of entries by list bucket call", Value: "1000"}
	templateFlag      = cli.StringFlag{Name: "template", Usage: "template for matching object names"}
	copiesFlag        = cli.IntFlag{Name: "copies", Usage: "number of object replicas", Value: 1}
	maxPagesFlag      = cli.IntFlag{Name: "max-pages", Usage: "display up to this number pages of bucket objects"}
	allFlag           = cli.BoolTFlag{Name: "all", Usage: "show all items including old, duplicated etc"}
	fastFlag          = cli.BoolTFlag{Name: "fast", Usage: "use fast API to list all object names in a bucket. Flags 'props', 'all', 'limit', and 'page-size' are ignored in this mode"}
	pagedFlag         = cli.BoolFlag{Name: "paged", Usage: "fetch and print the bucket list page by page, ignored in fast mode"}
	propsFlag         = cli.BoolFlag{Name: "props", Usage: "properties of a bucket"}
	showUnmatchedFlag = cli.BoolTFlag{Name: "show-unmatched", Usage: "also list objects that were not matched by regex and template"}

	// Daeclu
	countFlag = cli.IntFlag{Name: "count", Usage: "total number of generated reports", Value: countDefault}

	// Download
	descriptionFlag = cli.StringFlag{Name: "description,desc", Usage: "description of the job - can be useful when listing all downloads"}
	timeoutFlag     = cli.StringFlag{Name: "timeout", Usage: "timeout for request to external resource, eg. '30m'"}
	verboseFlag     = cli.BoolFlag{Name: "verbose,v", Usage: "verbose"}

	// dSort
	dsortBucketFlag   = cli.StringFlag{Name: "bucket", Value: cmn.DSortNameLowercase + "-testing", Usage: "bucket where shards will be put"}
	dsortTemplateFlag = cli.StringFlag{Name: "template", Value: "shard-{0..9}", Usage: "template of input shard name"}
	extFlag           = cli.StringFlag{Name: "ext", Value: ".tar", Usage: "extension for shards (either '.tar' or '.tgz')"}
	fileSizeFlag      = cli.StringFlag{Name: "fsize", Value: "1024", Usage: "single file size inside the shard"}
	logFlag           = cli.StringFlag{Name: "log", Usage: "path to file where the metrics will be saved"}
	cleanupFlag       = cli.BoolFlag{Name: "cleanup", Usage: "when set, the old bucket will be deleted and created again"}
	concurrencyFlag   = cli.IntFlag{Name: "conc", Value: 10, Usage: "limits number of concurrent put requests and number of concurrent shards created"}
	fileCountFlag     = cli.IntFlag{Name: "fcount", Value: 5, Usage: "number of files inside single shard"}

	// Object
	deadlineFlag = cli.StringFlag{Name: "deadline", Usage: "amount of time (Go Duration string) before the request expires", Value: "0s"}
	fileFlag     = cli.StringFlag{Name: "file", Usage: "filepath for content of the object"}
	lengthFlag   = cli.StringFlag{Name: "length", Usage: "object read length"}
	listFlag     = cli.StringFlag{Name: "list", Usage: "comma separated list of object names, eg. 'o1,o2,o3'"}
	nameFlag     = cli.StringFlag{Name: "name", Usage: "name of object"}
	newNameFlag  = cli.StringFlag{Name: "new-name", Usage: "new name of object"}
	outFileFlag  = cli.StringFlag{Name: "out-file", Usage: "name of the file where the contents will be saved"}
	offsetFlag   = cli.StringFlag{Name: "offset", Usage: "object read offset"}
	rangeFlag    = cli.StringFlag{Name: "range", Usage: "colon separated interval of object indices, eg. <START>:<STOP>"}
	cachedFlag   = cli.BoolFlag{Name: "cached", Usage: "check if an object is cached"}
	checksumFlag = cli.BoolFlag{Name: "checksum", Usage: "validate checksum"}
	waitFlag     = cli.BoolTFlag{Name: "wait", Usage: "wait for operation to finish before returning response"}
)

// Command argument texts
const (
	// Common
	idArgumentText            = "ID"
	noArgumentsText           = " "
	keyValuePairArgumentsText = "KEY=VALUE [KEY=VALUE...]"

	// Bucket
	bucketArgumentText       = "BUCKET_NAME"
	bucketRenameArgumentText = bucketArgumentText + " NEW_NAME"
	bucketPropsArgumentText  = bucketArgumentText + " " + keyValuePairArgumentsText

	// Config
	daemonIDArgumentText  = "[DAEMON_ID]"
	configSetArgumentText = daemonIDArgumentText + " " + keyValuePairArgumentsText

	// Daeclu
	daemonTypeArgumentText = "[DAEMON_TYPE]"
	targetIDArgumentText   = "[TARGET_ID]"

	// Download
	downloadStartArgumentText = "SOURCE DESTINATION"

	// dSort
	jsonSpecArgumentText = "JSON_SPECIFICATION"

	// Xaction
	xactionNameArgumentText         = "XACTION_NAME"
	xactionNameOptionalArgumentText = "[XACTION_NAME]"
)
