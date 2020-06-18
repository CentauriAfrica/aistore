// Package fs provides mountpath and FQN abstractions and methods to resolve/map stored content
/*
 * Copyright (c) 2019, NVIDIA CORPORATION. All rights reserved.
 */
package fs_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/NVIDIA/aistore/cmn"
	"github.com/NVIDIA/aistore/fs"
	"github.com/NVIDIA/aistore/ios"
	"github.com/NVIDIA/aistore/tutils/tassert"
)

func TestWalkBck(t *testing.T) {
	var (
		bck   = cmn.Bck{Name: "name", Provider: cmn.ProviderAIS}
		tests = []struct {
			name     string
			mpathCnt int
			sorted   bool
		}{
			{name: "simple_sorted", mpathCnt: 1, sorted: true},
			{name: "10mpaths_sorted", mpathCnt: 10, sorted: true},
		}
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs.Mountpaths = fs.NewMountedFS(ios.NewIOStaterMock())
			fs.Mountpaths.DisableFsIDCheck()
			_ = fs.CSM.RegisterContentType(fs.ObjectType, &fs.ObjectContentResolver{})

			mpaths := make([]string, 0, test.mpathCnt)
			defer func() {
				for _, mpath := range mpaths {
					os.RemoveAll(mpath)
				}
			}()

			for i := 0; i < test.mpathCnt; i++ {
				mpath, err := ioutil.TempDir("", "testwalk")
				tassert.CheckFatal(t, err)

				err = cmn.CreateDir(mpath)
				tassert.CheckFatal(t, err)

				err = fs.Mountpaths.Add(mpath)
				tassert.CheckFatal(t, err)

				mpaths = append(mpaths, mpath)
			}

			avail, _ := fs.Mountpaths.Get()
			var fileNames []string
			for _, mpath := range avail {
				dir := mpath.MakePathCT(bck, fs.ObjectType)
				err := cmn.CreateDir(dir)
				tassert.CheckFatal(t, err)

				_, names := prepareDirTree(t, dirTreeDesc{
					initDir: dir,
					dirs:    rand.Int()%100 + 1,
					files:   rand.Int()%100 + 1,
					depth:   rand.Int()%4 + 1,
					empty:   false,
				})
				fileNames = append(fileNames, names...)
			}

			var (
				objs = make([]string, 0, 100)
				fqns = make([]string, 0, 100)
			)
			err := fs.WalkBck(&fs.Options{
				Bck:         bck,
				CTs:         []string{fs.ObjectType},
				ErrCallback: nil,
				Callback: func(fqn string, de fs.DirEntry) error {
					parsedFQN, err := fs.Mountpaths.ParseFQN(fqn)
					tassert.CheckError(t, err)
					objs = append(objs, parsedFQN.ObjName)
					fqns = append(fqns, fqn)
					return nil
				},
				Sorted: test.sorted,
			})
			tassert.CheckFatal(t, err)

			sorted := sort.IsSorted(sort.StringSlice(objs))
			tassert.Fatalf(t, sorted == test.sorted, "expected the output to be sorted=%t", test.sorted)

			sort.Strings(fqns)
			sort.Strings(fileNames)
			tassert.Fatalf(t, reflect.DeepEqual(fqns, fileNames), "found objects don't match expected objects")
		})
	}
}
