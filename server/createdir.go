package main

import (
	pb "alex/gofigure/proto"
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
		return &pb.Result{Success: false, Msg: "Failed to parse mode."}, err
	}

	// enforce full path
	if !filepath.IsAbs(props.Path) {
		return &pb.Result{Success: false, Msg: "Path is not an absolute path."}, err
	}

	// Create the directory, ignoring errors if it already exists
	err = os.Mkdir(props.Path, mode)
	if (err != nil) && !os.IsExist(err) {
		return &pb.Result{Success: false, Msg: "Failed to create directory."}, err
	}

	// Get user and group from OS
	owner, err := user.Lookup(props.Owner)
	if err != nil {
		owner, err = user.LookupId(props.Owner)
		if err != nil {
			return &pb.Result{Success: false, Msg: "Failed to find user."}, err
		}
	}
	uid, err := strconv.Atoi(owner.Uid)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to parse uid."}, err
	}
	group, err := user.LookupGroup(props.Group)
	if err != nil {
		group, err = user.LookupGroupId(props.Group)
		if err != nil {
			return &pb.Result{Success: false, Msg: "Failed to find group."}, err
		}
	}
	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to parse gid for."}, err
	}

	// chown
	err = os.Lchown(props.Path, uid, gid)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to chown directory."}, err
	}

	// chmod
	err = os.Chmod(props.Path, mode)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to chmod directory."}, err
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
