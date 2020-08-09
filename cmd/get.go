// Copyright © 2016 Dropbox, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/ioprogress"
	"github.com/spf13/cobra"
)

type DownloadEntry struct {
	path        string
	destination string
}

func get(cmd *cobra.Command, args []string) (err error) {

	if len(args) == 0 || len(args) > 2 {
		return errors.New("`get` requires `src` and/or `dst` arguments")
	}

	src, err := validatePath(args[0])
	if err != nil {
		return
	}

	dst := path.Base(src)
	if len(args) == 2 {
		dst = args[1]
	}

	recurse, _ := cmd.Flags().GetBool("recurse")
	var entries []DownloadEntry

	if !recurse {
		entries = append(entries, DownloadEntry{
			path: src,
			dest: dest,
		})
	} else {
		// go through each file

	}

	for entry := range entries {
		// Default `dst` to the base segment of the source path; use the second argument if provided.
		// If `dst` is a directory, append the source filename.
		if f, err := os.Stat(dst); err == nil && f.IsDir() {
			dst = path.Join(dst, path.Base(entry.Path))
		}

		arg := files.NewDownloadArg(entry.path)

		dbx := files.New(config)
		res, contents, err := dbx.Download(arg)
		if err != nil {
			return
		}
		defer contents.Close()

		f, err := os.Create(dst)
		if err != nil {
			return
		}
		defer f.Close()

		progressbar := &ioprogress.Reader{
			Reader: contents,
			DrawFunc: ioprogress.DrawTerminalf(os.Stderr, func(progress, total int64) string {
				return fmt.Sprintf("Downloading %s/%s",
					humanize.IBytes(uint64(progress)), humanize.IBytes(uint64(total)))
			}),
			Size: int64(res.Size),
		}

		if _, err = io.Copy(f, progressbar); err != nil {
			return
		}

	}
	return
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "hi", // "get  [flags] <source> [<target>]",
	Short: "Download a file",
	RunE:  get,
}

func init() {
	getCmd.Flags().BoolP("recurse", "R", false, "Download everything under directory")
	RootCmd.AddCommand(getCmd)
}
