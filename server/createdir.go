package main

import (
	pb "alex/gofigure/proto"
	"errors"
	"golang.org/x/net/context"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

// TODO add tests for invalid modes
func parseFileMode(s string) (os.FileMode, error) {
	mode, err := strconv.ParseUint(s, 8, 16)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}

func (s *GofigureServer) GofigureDirectory(ctx context.Context, req *pb.FileRequest) (*pb.Result, error) {
	props := req.Properties

	// Get a FileMode object
	mode, err := parseFileMode(props.Mode)
	if err != nil {
		return nil, err
	}

	// enforce full path
	if !filepath.IsAbs(props.Path) {
		return nil, errors.New("Must be an absolute path.")
	}

	// Get user ID
	owner, err := user.Lookup(props.Owner)
	if err != nil {
		owner, err = user.LookupId(props.Owner)
		if err != nil {
			return nil, err
		}
	}
	uid, err := strconv.Atoi(owner.Uid)
	if err != nil {
		return nil, err
	}

	// Get group ID
	group, err := user.LookupGroup(props.Group)
	if err != nil {
		group, err = user.LookupGroupId(props.Group)
		if err != nil {
			return nil, err
		}
	}
	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return nil, err
	}

	// Create the directory, ignoring errors if it already exists
	err = os.Mkdir(props.Path, mode)
	if (err != nil) && !os.IsExist(err) {
		return nil, err
	}

	// chown
	err = os.Lchown(props.Path, uid, gid)
	if err != nil {
		return nil, err
	}

	// TODO find Lchmod?
	// chmod
	err = os.Chmod(props.Path, mode)
	if err != nil {
		return nil, err
	}

	// how to hand enum
	//switch x := m.Avatar.(type) {
	//case *account.Profile_ImageUrl:
	//	// Load profile image based on URL
	//	// using x.ImageUrl
	//case *account.Profile_ImageData:
	//	// Load profile image based on bytes
	//	// using x.ImageData
	//case nil:
	//	// The field is not set.
	//default:
	//	return fmt.Errorf("Profile.Avatar has unexpected type %T", x)
	//}
	return &pb.Result{Success: true}, nil
}
