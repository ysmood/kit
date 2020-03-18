package os

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hectane/go-acl"
	"github.com/karrick/godirwalk"
	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
	"github.com/ysmood/kit/pkg/utils"
)

// Copy file or dir recursively
var Copy = copy.Copy

// Chmod ...
var Chmod = acl.Chmod

// HomeDir returns the current user's home dir path
func HomeDir() string {
	p, _ := homedir.Dir()
	return p
}

// MkdirOptions ...
type MkdirOptions struct {
	Perm os.FileMode
}

// Mkdir makes dir recursively
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

// OutputFile auto creates file if not exists, it will try to detect the data type and
// auto output binary, string or json
func OutputFile(p string, data interface{}, options *OutputFileOptions) error {
	if options == nil {
		options = &OutputFileOptions{0775, 0664, "", "    "}
	}

	dir := filepath.Dir(p)
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

// ReadFile reads file as bytes
func ReadFile(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}

// ReadString reads file as string
func ReadString(p string) (string, error) {
	bin, err := ioutil.ReadFile(p)
	return string(bin), err
}

// ReadJSON reads file as json
func ReadJSON(p string, data interface{}) error {
	bin, err := ReadFile(p)

	if err != nil {
		return err
	}

	return json.Unmarshal(bin, data)
}

// Move file or folder to another location, create path if needed
func Move(from, to string, perm *os.FileMode) error {
	err := Mkdir(filepath.Dir(to), nil)

	if err != nil {
		return err
	}

	return os.Rename(from, to)
}

// Remove dirs, files, patterns as expected.
func Remove(patterns ...string) error {
	return RemoveWithDir("", patterns...)
}

// RemoveWithDir is the low level of Remove
func RemoveWithDir(dir string, patterns ...string) error {
	return Walk(patterns...).
		Dir(dir).
		PostChildrenCallback(func(dir string, info *godirwalk.Dirent) error {
			return os.RemoveAll(dir)
		}).
		Do(func(p string, info *godirwalk.Dirent) error {
			if info.IsDir() {
				return nil
			}
			return os.Remove(p)
		})
}

// Exists checks if file or dir exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileExists checks if file exists, only for file, not for dir
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

// DirExists checks if file exists, only for dir, not for file
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

// CD to a dir and return a function to cd back to the preivous dir
func CD(dir string) func() {
	curr, err := os.Getwd()
	utils.E(err)

	utils.E(os.Chdir(dir))

	return func() {
		utils.E(os.Chdir(curr))
	}
}
