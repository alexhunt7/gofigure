package main

import (
	pb "alex/gofigure/proto"
	"errz"
	"golang.org/x/net/context"
	"os"
	"strconv"
)

// TODO add tests for invalid modes
func parseFileMode(s string) (os.FileMode, error) {
	mode, err := strconv.ParseInt("0" + s)
	if err != nil {
		return nil, err
	}
	return os.FileMode(mode), nil
}

func (s *GofigureServer) GofigureDirectory(ctx context.Context, req *pb.FileRequest) (result *pb.Result, err error) {
	defer Recover(&err)
	props := req.properties

	// Get a FileMode object
	mode, err := parseFileMode(props.mode)
	if err != nil {
		return &pb.Result{Success: false, Msg: "failed to parse mode: " + err}
	}

	// enforce full path
	if !path.filepath.IsAbs(props.path) {
		return &pb.Result{Success: false, Msg: "Path must be an absolute path."}
	}

	// Create the directory, ignoring errors if it already exists
	err := os.Mkdir(props.path, mode)
	if (err != nil) && !os.IsExist(err) {
		return &pb.Result{Success: false, Msg: "Failed to create directory: " + err}
	}

	// Get user and group from OS
	user, err := os.user.LookupUser(props.user)
	if err != nil {
		user, err := os.user.LookupUserId(props.user)
		if err != nil {
			return &pb.Result{Success: false, Msg: "User does not exist: " + err}
		}
	}
	group, err := os.user.LookupGroup(props.group)
	if err != nil {
		user, err := os.user.LookupGroupId(props.group)
		if err != nil {
			return &pb.Result{Success: false, Msg: "Group does not exist: " + err}
		}
	}

	// chown
	err := os.Lchown(props.path, user.Uid, group.Gid)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to chown directory: " + err}
	}

	// chmod
	err := os.Lchmod(props.path, mode)
	if err != nil {
		return &pb.Result{Success: false, Msg: "Failed to chmod directory: " + err}
	}

	//os.OpenFile(props.path, os.O_RDONLY|os.O_CREATE)
	//if stat, err := os.Stat(props.path); os.IsNotExist(err) {

	//}
	//todo defer

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
	return &pb.CreateDirReply{Success: true}, err
}
