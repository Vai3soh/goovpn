package fileextended

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	path     string
	body     []byte
	destPath string
	perm     fs.FileMode
}

type Option func(*File)

func NewFile(opts ...Option) *File {
	f := &File{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithPath(path string) Option {
	return func(a *File) {
		a.path = path
	}
}

func WithBody(body []byte) Option {
	return func(a *File) {
		a.body = body
	}
}

func WithDestPath(destPath string) Option {
	return func(a *File) {
		a.destPath = destPath
	}
}

func WithPermisson(perm fs.FileMode) Option {
	return func(a *File) {
		a.perm = perm
	}
}

func (f *File) FilesInDir(dir string) ([]string, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var arr []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".ovpn" || filepath.Ext(file.Name()) == ".conf" {
			arr = append(arr, strings.TrimSuffix(file.Name(),
				filepath.Ext(file.Name())))
		}
	}
	return arr, nil
}

func (f *File) ReadFileAsByte() ([]byte, error) {
	body, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (f *File) ReadFileAsString() (*string, error) {
	body, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	s := string(body)
	return &s, nil
}

func (f *File) FileOpen() (*os.File, error) {
	inFile, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	return inFile, nil
}

func (f *File) CopyFile() error {
	err := ioutil.WriteFile(f.destPath, f.body, f.perm)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) AbsolutePath() (*string, error) {
	full, err := filepath.Abs(f.path)
	if err != nil {
		return nil, err
	}
	return &full, nil
}

func (f *File) SetBody(body []byte) {
	f.body = body
}

func (f *File) Body() []byte {
	return f.body
}

func (f *File) Path() string {
	return f.path
}

func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) SetDestPath(destPath string) {
	f.destPath = destPath
}

func (f *File) SetPermissonFile(perm fs.FileMode) {
	f.perm = perm
}

func (f *File) FileNameWithoutExtension() *string {
	n := (f.path)[:len(f.path)-len(filepath.Ext(f.path))]
	name := filepath.Base(n)
	return &name
}

func (f *File) CreateFile() (*os.File, error) {
	file, err := os.Create(f.path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *File) WriteStringToFile(file *os.File, data string) error {
	defer file.Close()
	_, err := file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) WriteByteFile() error {
	err := ioutil.WriteFile(f.path, f.body, f.perm)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) CheckFileExists() bool {
	if _, err := os.Stat(f.path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
