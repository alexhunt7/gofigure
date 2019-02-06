package minion

import (
	pb "github.com/alexhunt7/gofigure/proto"
	"golang.org/x/net/context"
	"os"
	"os/user"
	"strconv"
	"testing"
	"time"
)

func TestParseFileMode(t *testing.T) {
	tests := []struct {
		smode string
		fmode os.FileMode
	}{
		{"0777", os.FileMode(uint(511))},
		{"777", os.FileMode(uint(511))},
		{"0644", os.FileMode(uint(420))},
		{"0000", os.FileMode(uint(0))},
		{"0002", os.FileMode(uint(2))},
		{"0400", os.FileMode(uint(256))},
		{"0600", os.FileMode(uint(384))},
	}

	for _, tt := range tests {
		omode, _ := parseFileMode(tt.smode)
		if omode != tt.fmode {
			t.Errorf("TestParseFileMode, %s: %v != %v", tt.smode, omode, tt.fmode)
		}
	}
}

func TestParseFileModeError(t *testing.T) {
	tests := []string{
		"asdf",
		"a",
		"",
		"!",
		"10.4",
		"9",
	}

	for _, tt := range tests {
		_, err := parseFileMode(tt)
		if err == nil {
			t.Errorf("TestParseFileModeError did not error as expected")
		}
	}
}

func TestOwnMod(t *testing.T) {
	// TODO
}

func TestParseFileProperties(t *testing.T) {
	tests := []struct {
		fileRequest  *pb.FileRequest
		expectedPath string
		expectedMode os.FileMode
		expectedUID  int
		expectedGID  int
		shouldErr    bool
	}{
		// TODO find a way to test this stuff cross environment. Mock? Docker?
		{
			// Simplest case
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "/home/alex/git",
					Owner: "alex",
					Group: "alex",
					Mode:  "700",
				},
			},
			expectedPath: "/home/alex/git",
			expectedUID:  1000,
			expectedGID:  1000,
			expectedMode: os.FileMode(uint(448)),
		},
		{
			// Check that it handles uid/gid numbers
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "/home/alex",
					Owner: "1000",
					Group: "1000",
					Mode:  "700",
				},
			},
			expectedPath: "/home/alex",
			expectedUID:  1000,
			expectedGID:  1000,
			expectedMode: os.FileMode(uint(448)),
		},
		// TODO make it trim relative paths
		//{
		//	// Check that it handles full relative paths
		//	fileRequest: &pb.FileRequest{
		//		Properties: &pb.FileProperties{
		//			Path:  "/home/alex/../alex/git",
		//			Owner: "alex",
		//			Group: "alex",
		//			Mode:  "700",
		//		},
		//	},
		//	expectedPath: "/home/alex/git",
		//	expectedUID:  1000,
		//	expectedGID:  1000,
		//	expectedMode: os.FileMode(uint(448)),
		//},
		{
			// Check that it errors on ~/
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "~/git",
					Owner: "alex",
					Group: "alex",
					Mode:  "700",
				},
			},
			shouldErr: true,
		},
		{
			// Check that it errors on relative paths
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "git",
					Owner: "alex",
					Group: "alex",
					Mode:  "700",
				},
			},
			shouldErr: true,
		},
		{
			// Check that it errors on non-existent uid
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "/home/alex",
					Owner: "32487",
					Group: "1000",
					Mode:  "700",
				},
			},
			shouldErr: true,
		},
		{
			// Check that it errors on non-existent gid
			fileRequest: &pb.FileRequest{
				Properties: &pb.FileProperties{
					Path:  "/home/alex",
					Owner: "1000",
					Group: "32487",
					Mode:  "700",
				},
			},
			shouldErr: true,
		},
	}
	for _, tt := range tests {
		path, mode, uid, gid, err := parseFileProperties(tt.fileRequest)
		if tt.shouldErr {
			if err == nil {
				t.Errorf("TestParseFileProperties did not error as expected")
			}
			continue
		} else {
			if err != nil {
				t.Errorf("TestParseFileProperties got unexpected error: %v", err)
			}
		}
		if path != tt.expectedPath {
			t.Errorf("TestParseFileProperties path, %v != %v", path, tt.expectedPath)
		}
		if mode != tt.expectedMode {
			t.Errorf("TestParseFileProperties mode, %v != %v", mode, tt.expectedMode)
		}
		if uid != tt.expectedUID {
			t.Errorf("TestParseFileProperties uid, %v != %v", uid, tt.expectedUID)
		}
		if gid != tt.expectedGID {
			t.Errorf("TestParseFileProperties gid, %v != %v", gid, tt.expectedGID)
		}
	}
}

func TestFile(t *testing.T) {
	minion := &Minion{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get working directory: %v", err)
	}
	testdir := wd + "/testdata/dir"
	testfile := testdir + "/testfile"
	err = os.RemoveAll(testdir)
	if err != nil {
		t.Errorf("Failed to remove testdir: %v", err)
	}

	statReq := &pb.FilePath{Path: testfile}
	stat, err := minion.Stat(ctx, statReq)
	if err != nil {
		t.Errorf("Failed to run Stat: %v", err)
	}
	if stat.Exists {
		t.Errorf("Failed to remove testdir")
	}

	user, err := user.Current()
	if err != nil {
		t.Errorf("Failed to get current user: %v", err)
	}

	req := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  testdir,
			Owner: user.Username,
			Group: user.Gid,
			Mode:  "700",
		},
	}

	_, err = minion.Directory(ctx, req)
	if err != nil {
		t.Errorf("Failed to run Directory: %v", err)
	}

	req.Content = []byte("Hello world!\n")
	req.Properties.Path = testfile

	_, err = minion.File(ctx, req)
	if err != nil {
		t.Errorf("Failed to run File: %v", err)
	}

	stat, err = minion.Stat(ctx, statReq)
	if err != nil {
		t.Errorf("Failed to run Stat: %v", err)
	}
	if !stat.Exists {
		t.Errorf("Does not exist")
	}
	if stat.Mode != "700" {
		t.Errorf("Mode not 700: %s", stat.Mode)
	}
	uid, _ := strconv.Atoi(user.Uid)
	if stat.Uid != uint32(uid) {
		t.Errorf("Uid mismatch: %d != %s", stat.Uid, user.Uid)
	}
	gid, _ := strconv.Atoi(user.Gid)
	if stat.Gid != uint32(gid) {
		t.Errorf("Gid mismatch: %d != %s", stat.Gid, user.Gid)
	}
	if stat.Size != 13 {
		t.Errorf("Size mismatch: %d", stat.Size)
	}

	req.Properties.Mode = "0600"
	_, err = minion.File(ctx, req)
	if err != nil {
		t.Errorf("Failed to run File: %v", err)
	}
	stat, err = minion.Stat(ctx, statReq)
	if err != nil {
		t.Errorf("Failed to run Stat: %v", err)
	}
	if stat.Mode != "600" {
		t.Errorf("Mode not 700: %s", stat.Mode)
	}
}

func BenchmarkFile(b *testing.B) {
	minion := &Minion{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil {
		b.Errorf("Failed to get working directory: %v", err)
	}
	testdir := wd + "/testdata/dir"
	err = os.RemoveAll(testdir)
	if err != nil {
		b.Errorf("Failed to remove testdir: %v", err)
	}

	user, err := user.Current()
	if err != nil {
		b.Errorf("Failed to get current user: %v", err)
	}

	req := &pb.FileRequest{
		Properties: &pb.FileProperties{
			Path:  testdir,
			Owner: user.Username,
			Group: user.Gid,
			Mode:  "700",
		},
	}

	_, err = minion.Directory(ctx, req)
	if err != nil {
		b.Errorf("Failed to run Directory: %v", err)
	}

	b.ResetTimer()
	req.Content = []byte("Hello world!\n")
	for i := 0; i < b.N; i++ {
		//req.Properties.Path = fmt.Sprintf("%s/%d", testdir, i)
		req.Properties.Path = testdir + "/" + strconv.Itoa(i)

		_, err = minion.File(ctx, req)
		if err != nil {
			b.Errorf("Failed to run File: %v", err)
		}
	}
}

func TestDelete(t *testing.T) {
	// TODO test delete file
	// TODO test delete empty dir
	// TODO test delete dir with file in it
	minion := &Minion{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get working directory: %v", err)
	}

	testdir := wd + "/testdata/dir"
	testfile := testdir + "/afile"
	err = os.RemoveAll(testdir)
	if err != nil {
		t.Errorf("Failed to remove testdir: %v", err)
	}

	err = os.MkdirAll(testdir, 0700)
	if err != nil {
		t.Errorf("Failed to create testdir: %v", err)
	}

	f, err := os.OpenFile(testfile, os.O_RDONLY|os.O_CREATE, 0700)
	if err != nil {
		t.Errorf("Failed to create testfile: %v", err)
	}
	f.Close()

	req := &pb.DeleteRequest{
		Path: testfile,
	}
	_, err = minion.Delete(ctx, req)
	if err != nil {
		t.Errorf("Failed to run Delete: %v", err)
	}
	req = &pb.DeleteRequest{
		Path: testdir,
	}
	_, err = minion.Delete(ctx, req)
	if err != nil {
		t.Errorf("Failed to run Delete: %v", err)
	}

	err = os.MkdirAll(testdir, 0700)
	if err != nil {
		t.Errorf("Failed to create testdir: %v", err)
	}

	f, err = os.OpenFile(testfile, os.O_RDONLY|os.O_CREATE, 0700)
	if err != nil {
		t.Errorf("Failed to create testfile: %v", err)
	}
	f.Close()
	_, err = minion.Delete(ctx, req)
	if err == nil {
		t.Errorf("Succeeded running Delete when it should have failed.")
	}
	req = &pb.DeleteRequest{
		Path:      testdir,
		Recursive: true,
	}
	_, err = minion.Delete(ctx, req)
	if err != nil {
		t.Errorf("Failed to run Delete: %v", err)
	}
}
