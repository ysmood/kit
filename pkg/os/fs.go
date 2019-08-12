package os

import (
	"encoding/json"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/hectane/go-acl"
	"github.com/karrick/godirwalk"
	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
)

// Copy copy file or dir recursively
var Copy = copy.Copy

// Chmod ...
var Chmod = acl.Chmod

// HomeDir current user's home dir path
func HomeDir() string {
	p, _ := homedir.Dir()
	return p
}

// ThisFilePath get the current file path
func ThisFilePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// ThisDirPath get the current file directory path
func ThisDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

// GoPath get the current GOPATH properly
func GoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// MkdirOptions ...
type MkdirOptions struct {
	Perm os.FileMode
}

// Mkdir make dir recursively
func Mkdir(path string, options *MkdirOptions) error {
	if options == nil {
		options = &MkdirOptions{
			Perm: 0775,
		}
	}

	return os.MkdirAll(path, options.Perm)
}

// OutputFileOptions ...
type OutputFileOptions struct {
	DirPerm    os.FileMode
	FilePerm   os.FileMode
	JSONPrefix string
	JSONIndent string
}

// OutputFile auto create file if not exists, it will try to detect the data type and
// auto output binary, string or json
func OutputFile(p string, data interface{}, options *OutputFileOptions) error {
	if options == nil {
		options = &OutputFileOptions{0775, 0664, "", "    "}
	}

	dir := path.Dir(p)
	_ = Mkdir(dir, &MkdirOptions{Perm: options.DirPerm})

	var bin []byte

	switch t := data.(type) {
	case []byte:
		bin = t
	case string:
		bin = []byte(t)
	default:
		var err error
		bin, err = json.MarshalIndent(data, options.JSONPrefix, options.JSONIndent)

		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(p, bin, options.FilePerm)
}

// ReadFile read file as bytes
func ReadFile(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}

// ReadString read file as string
func ReadString(p string) (string, error) {
	bin, err := ioutil.ReadFile(p)
	return string(bin), err
}

// ReadJSON read file as json
func ReadJSON(p string, data interface{}) error {
	bin, err := ReadFile(p)

	if err != nil {
		return err
	}

	return json.Unmarshal(bin, data)
}

// Move move file or folder to another location, create path if needed
func Move(from, to string, perm *os.FileMode) error {
	err := Mkdir(path.Dir(to), nil)

	if err != nil {
		return err
	}

	return os.Rename(from, to)
}

// Remove remove dirs, files, patterns as expected.
// The pattern cannot be absolute path.
func Remove(patterns ...string) error {
	// if any of the patterns is a raw folder path not a pattern remove all children of it
	for _, p := range patterns {
		if DirExists(p) {
			err := Remove(p + "/**")
			if err != nil {
				return err
			}
		}
	}

	return Walk(patterns...).PostChildrenCallback(func(p string, info *godirwalk.Dirent) error {
		return os.Remove(p)
	}).Do(func(p string, info *godirwalk.Dirent) error {
		if info.IsDir() {
			return nil
		}
		return os.Remove(p)
	})
}

// Exists check if file or dir exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileExists check if file exists, only for file, not for dir
func FileExists(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// DirExists check if file exists, only for dir, not for file
func DirExists(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}
