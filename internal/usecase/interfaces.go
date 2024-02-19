package usecase

import (
	"io/fs"
	"os"

	"github.com/Vai3soh/goovpn/entity"
)

type (
	SessionSetters interface {
		SetConfig(config string)
		SetCread(u, p string) error
	}

	SessionLoger interface {
		ChanVpnLog() chan string
	}

	SessionManager interface {
		StartSession() error
		StopSession()
		DestroyClient()
	}

	ConfigSetters interface {
		SetBody(body string)
		SetPath(path string)
	}

	ConfigBody interface {
		GetBody() string
	}

	ConfigRemover interface {
		RemoveSpaceLines()
		RemoveCommentLines()
		RemoveEmptyString()
		RemoveNotCertsAndKeys()
		RemoveCertsAndKeys()
	}

	ConfigChecker interface {
		CheckConfigUseFiles() bool
		CheckStringAuthUserPass() bool
	}

	ConfigMerger interface {
		MergeCertsAndKeys(key string) string
	}

	ConfigTools interface {
		AddStringToConfig()
		SearchFilesPaths() map[string]string
		GetAuthpathFileName() string
		GetUserAndPass() (string, string)
	}

	FileSetters interface {
		SetBody([]byte)
		SetPath(path string)
		SetDestPath(destPath string)
	}

	FileGetters interface {
		Path() string
		Body() []byte
	}

	FileReader interface {
		ReadFileAsByte() ([]byte, error)
		ReadFileAsString() (*string, error)
	}

	FileWriter interface {
		WriteStringToFile(file *os.File, data string) error
		WriteByteFile() error
	}

	FileToolsManager interface {
		FileNameWithoutExtension() *string
		FileOpen() (*os.File, error)
		SetPermissonFile(fs.FileMode)
		CreateFile() (*os.File, error)
		AbsolutePath() (*string, error)
		CopyFile() error
		CheckFileExists() bool
	}

	DnsSetters interface {
		SetInterface(iface string)
	}

	ProfileRepository interface {
		Store(p entity.Profile)
		Find(key string) entity.Profile
	}
)
