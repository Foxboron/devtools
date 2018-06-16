package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var errSourceNotDirectory = errors.New("source is not a directory")
var errDestinationExists = errors.New("destination already exists")

/*
 *
 * From: https://gist.github.com/kuznero/41acd2afdbe7cfd8e135c4573041e1da
 *
 *
 * MIT License
 *
 * Copyright (c) 2017 Roland Singer [roland.singer@desertbit.com]
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// CopyFile copies the contents from one file to another.
// If the file does not exist, it will be created.
// If the file exists, the contents will be replaced.
// The file mode of the source file will also be set on the destination file.
// The copied data is synced after being copied.
func CopyFile(source, destination string) (err error) {

	// Open the source file
	in, err := os.Open(source)
	if err != nil {
		return
	}
	defer in.Close()

	// Create a new file, if it doesn't already exist
	out, err := os.Create(destination)
	if err != nil {
		return
	}

	// If the output stream can not be closed at the end,
	// return that as the error.
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	// Copy over the data
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	// Sync/flush
	err = out.Sync()
	if err != nil {
		return
	}

	// Retrieve the file mode
	si, err := os.Stat(source)
	if err != nil {
		return
	}

	// Apply the file mode to the destination file
	err = os.Chmod(destination, si.Mode())
	if err != nil {
		return
	}

	// Success
	return nil
}

// CopyDir recursively copies a directory tree while attempting to
// preserve permissions.
// * The source directory must exist.
// * The destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(source, destination string) error {

	// Clean and shorten the paths to the respective directories
	srcDir := filepath.Clean(source)
	dstDir := filepath.Clean(destination)

	// Check if the source directory is there and valid
	sourceInfo, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !sourceInfo.IsDir() {
		return errSourceNotDirectory
	}

	// Check the destination directory is not there
	if _, err := os.Stat(dstDir); err == nil {
		return errDestinationExists
	} else if os.IsExist(err) {
		return err
	}

	// Create the destination directory
	if err := os.MkdirAll(dstDir, sourceInfo.Mode()); err != nil {
		return err
	}

	// Retrieve a list of entries from the source directory
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// Copy over files, recursively
	for _, entry := range entries {
		fileName := entry.Name()
		srcPath := filepath.Join(srcDir, fileName)
		dstPath := filepath.Join(dstDir, fileName)

		// File or directory?
		if !entry.IsDir() {

			// Skip symlinks
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			// Skip sockets. Probably don't need dirmngr
			if entry.Mode()&os.ModeSocket != 0 {
				continue
			}

			// Copy over file
			if err = CopyFile(srcPath, dstPath); err != nil {
				return err
			}

		} else {

			// Copy over directory
			if err = CopyDir(srcPath, dstPath); err != nil {
				return err
			}

		}
	}

	// Success
	return nil
}
