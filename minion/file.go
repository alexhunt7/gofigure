/*
Copyright 2018 Alex Hunt

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package minion

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"

	pb "github.com/alexhunt7/gofigure/proto"
)

// TODO add tests for invalid modes
func parseFileMode(s string) (os.FileMode, error) {
	mode, err := strconv.ParseUint(s, 8, 16)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}

func ownMod(path string, mode os.FileMode, uid int, gid int) error {
	// chown
	err := os.Lchown(path, uid, gid)
	if err != nil {
		return err
	}

	// TODO find Lchmod?
	// chmod
	err = os.Chmod(path, mode)
	if err != nil {
		return err
	}

	return nil
}

func parseFileProperties(props *pb.FileProperties) (string, os.FileMode, int, int, error) {
	path := props.Path
	var mode os.FileMode
	uid := 0
	gid := 0

	// Get a FileMode object
	mode, err := parseFileMode(props.Mode)
	if err != nil {
		return path, mode, uid, gid, err
	}

	// enforce full path
	if !filepath.IsAbs(path) {
		return path, mode, uid, gid, errors.New("must be an absolute path")
	}

	// Get user ID
	owner, err := user.Lookup(props.Owner)
	if err != nil {
		owner, err = user.LookupId(props.Owner)
		if err != nil {
			return path, mode, uid, gid, err
		}
	}
	uid, err = strconv.Atoi(owner.Uid)
	if err != nil {
		return path, mode, uid, gid, err
	}

	// Get group ID
	group, err := user.LookupGroup(props.Group)
	if err != nil {
		group, err = user.LookupGroupId(props.Group)
		if err != nil {
			return path, mode, uid, gid, err
		}
	}
	gid, err = strconv.Atoi(group.Gid)
	if err != nil {
		return path, mode, uid, gid, err
	}

	return path, mode, uid, gid, nil
}

func safeWrite(path string, content []byte, mode os.FileMode) error {
	f, err := ioutil.TempFile(filepath.Dir(path), ".gofig.")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	err = ioutil.WriteFile(f.Name(), content, mode)
	if err != nil {
		return err
	}

	err = os.Rename(f.Name(), path)
	if err != nil {
		return err
	}

	return nil
}

// Stat retrieves properties of a file or directory.
func (s *Minion) Stat(ctx context.Context, req *pb.FilePath) (*pb.StatResult, error) {
	fstat, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return &pb.StatResult{Exists: false}, nil
		}
		return nil, err
	}

	statt := fstat.Sys().(*syscall.Stat_t)
	uid := statt.Uid
	gid := statt.Gid

	result := pb.StatResult{
		Size:   fstat.Size(),
		Mode:   strconv.FormatUint(uint64(fstat.Mode().Perm()), 8),
		IsDir:  fstat.IsDir(),
		Exists: true,
		Uid:    uid,
		Gid:    gid,
	}

	owner, err := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	if err == nil {
		result.Owner = owner.Username
	}

	group, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))
	if err == nil {
		result.Group = group.Name
	}

	return &result, nil
}

// File creates a file with the specified mode, uid, gid, and contents.
func (s *Minion) File(ctx context.Context, req *pb.FileRequest) (*pb.FileResult, error) {
	path, mode, uid, gid, err := parseFileProperties(req.Properties)
	if err != nil {
		return nil, err
	}

	needsWrite, err := fileNeedsWrite(path, req.Content)
	if err != nil {
		return nil, err
	}

	if needsWrite {
		err = safeWrite(path, req.Content, mode)
		if err != nil {
			return nil, err
		}
	}

	// TODO compare stat with mode, uid, gid
	//fstat, err := os.Stat(path)
	//if err != nil {
	//	return nil, err
	//}

	err = ownMod(path, mode, uid, gid)
	if err != nil {
		return nil, err
	}

	return &pb.FileResult{}, nil
}

func fileNeedsWrite(path string, content []byte) (bool, error) {
	needsWrite := false

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			needsWrite = true
		} else {
			return false, err
		}
	} else {
		defer f.Close()

		// Compare existing file content
		hasher := sha256.New()
		if _, err := io.Copy(hasher, f); err != nil {
			return false, err
		}
		existingSum := hasher.Sum(nil)

		hasher.Reset()
		// TODO handle streaming?
		// TODO hash on the controller side?
		_, err = hasher.Write(content)
		if err != nil {
			return false, err
		}
		contentSum := hasher.Sum(nil)

		if !bytes.Equal(existingSum, contentSum) {
			needsWrite = true
		}
	}

	return needsWrite, nil
}

func (s *Minion) LineInFile(ctx context.Context, req *pb.LineInFileRequest) (*pb.LineInFileResult, error) {
	path, mode, uid, gid, err := parseFileProperties(req.Properties)
	if err != nil {
		return nil, err
	}

	existingContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var newContent []byte
	if req.Regex == "" {
		if bytes.Contains(existingContent, []byte(req.Line)) {
			// TODO shortcut out here instead? needs ownMod
			newContent = existingContent
		} else {
			newContent = append(newContent, existingContent...)
			newContent = append(newContent, []byte("\n")...)
			newContent = append(newContent, []byte(req.Line)...)
			newContent = append(newContent, []byte("\n")...)
		}
	} else {
		regex, err := regexp.Compile(fmt.Sprintf("(?m).*%s.*", req.Regex))
		if err != nil {
			return nil, fmt.Errorf("failed to compile regex: %v", err)
		}

		// sanity check that regex matches line
		if !regex.MatchString(req.Line) {
			return nil, fmt.Errorf("regular expression does not match line")
		}

		found := regex.Find(existingContent)
		if found != nil {
			newContent = bytes.Replace(existingContent, found, []byte(req.Line), 1)
		} else {
			// TODO DRY this up
			newContent = append(newContent, existingContent...)
			newContent = append(newContent, []byte("\n")...)
			newContent = append(newContent, []byte(req.Line)...)
			newContent = append(newContent, []byte("\n")...)
		}
	}

	needsWrite, err := fileNeedsWrite(path, newContent)
	if err != nil {
		return nil, err
	}

	if needsWrite {
		err = safeWrite(path, newContent, mode)
		if err != nil {
			return nil, err
		}
	}

	err = ownMod(path, mode, uid, gid)
	if err != nil {
		return nil, err
	}

	return &pb.LineInFileResult{}, nil
}

// Directory creates a directory with the specified mode, uid, and gid.
func (s *Minion) Directory(ctx context.Context, req *pb.FileRequest) (*pb.DirectoryResult, error) {
	path, mode, uid, gid, err := parseFileProperties(req.Properties)
	if err != nil {
		return nil, err
	}

	// Create the directory, ignoring errors if it already exists
	err = os.Mkdir(path, mode)
	if (err != nil) && !os.IsExist(err) {
		return nil, err
	}

	err = ownMod(path, mode, uid, gid)
	if err != nil {
		return nil, err
	}

	return &pb.DirectoryResult{}, nil
}

// Delete removes the specified file or directory, optionally recursively.
func (s *Minion) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResult, error) {
	var err error
	if req.Recursive {
		err = os.RemoveAll(req.Path)
	} else {
		err = os.Remove(req.Path)
	}
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResult{}, nil
}
