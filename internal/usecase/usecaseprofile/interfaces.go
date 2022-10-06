package usecaseprofile

import (
	"io/fs"
	"os"

	"github.com/Vai3soh/goovpn/entity"
)

type FileSetters interface {
	SetBody([]byte)
	SetPath(path string)
	SetDestPath(destPath string)
}

type FileToolsManager interface {
	FileNameWithoutExtension() *string
	FileOpen() (*os.File, error)
	SetPermissonFile(fs.FileMode)
	CreateFile() (*os.File, error)
	AbsolutePath() (*string, error)
	CopyFile() error
	CheckFileExists() bool
}

type FileReader interface {
	ReadFileAsByte() ([]byte, error)
	ReadFileAsString() (*string, error)
}

type ConfigSetters interface {
	SetBody(body string)
	SetPath(path string)
}

type ConfigMerger interface {
	MergeCertsAndKeys(key string) string
}

type ConfigRemover interface {
	RemoveSpaceLines()
	RemoveCommentLines()
	RemoveEmptyString()
	RemoveNotCertsAndKeys()
	RemoveCertsAndKeys()
}

type ConfigBody interface {
	GetBody() string
}
type ProfileRepository interface {
	Store(p entity.Profile)
	Find(key string) entity.Profile
	Delete(p entity.Profile)
}

type ConfigTools interface {
	AddStringToConfig(inFile *os.File) string
	SearchFilesPaths() map[string]string
	GetAuthpathFileName() string
	GetUserAndPass() (string, string)
}

type ConfigChecker interface {
	CheckConfigUseFiles() bool
	CheckStringAuthUserPass() bool
}

type WithCfgfileProfile interface {
	GetProfileFromDisk() error
	SetPath(path string)
	SaveProfile() error
}
