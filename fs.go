package gokit

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// MkdirOptions ...
type MkdirOptions struct {
	perm os.FileMode
}

// Mkdir make dir recursively
func Mkdir(path string, options *MkdirOptions) error {
	if options == nil {
		options = &MkdirOptions{
			perm: 0775,
		}
	}

	return os.MkdirAll(path, options.perm)
}

// OutputFileOptions ...
type OutputFileOptions struct {
	DirPerm    os.FileMode
	FilePerm   os.FileMode
	JSONPrefix string
	JSONIndent string
}

// OutputFile auto create file if not exists, it will try to detect the data type and
// auto output binary, string or gob
func OutputFile(p string, data interface{}, options *OutputFileOptions) error {
	if options == nil {
		options = &OutputFileOptions{0775, 0664, "", "    "}
	}

	dir := path.Dir(p)
	err := Mkdir(dir, &MkdirOptions{perm: options.DirPerm})

	if err != nil {
		return err
	}

	var bin []byte

	switch t := data.(type) {
	case []byte:
		bin = t
	case string:
		bin = []byte(t)
	default:
		bin, err = json.MarshalIndent(data, options.JSONPrefix, options.JSONIndent)

		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(p, bin, options.FilePerm)
}

// ReadFile ...
func ReadFile(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}

// ReadStringFile ...
func ReadStringFile(p string) (string, error) {
	bin, err := ioutil.ReadFile(p)
	return string(bin), err
}

// ReadJSON ...
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

// Remove ...
func Remove(p string) error {
	return os.RemoveAll(p)
}

// Exists check if file or dir exists
func Exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
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
