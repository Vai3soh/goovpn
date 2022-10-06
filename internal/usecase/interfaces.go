package usecase

import (
	"context"
	"io/fs"
	"os"

	"github.com/Vai3soh/goovpn/entity"
)

type (
	SessionSetters interface {
		SetConfig(config string)
		SetCread(user, pwd string)
		SetSession()
	}

	SessionManager interface {
		StartSession(ctx context.Context)
		StopSession()
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
		AddStringToConfig(inFile *os.File) string
		SearchFilesPaths() map[string]string
		GetAuthpathFileName() string
		GetUserAndPass() (string, string)
	}

	UiLoger interface {
		ChanVpnLog() chan string
	}

	UiLogFormManager interface {
		GetTextFromLogForm() string
		SetTextInLogForm(text string)
		ClearLogForm()
	}

	UiButtonsManager interface {
		ButtonConnectEnable()
		ButtonConnectDisable()
		ButtonDisconnectEnable()
		ButtonDisconnectDisable()
	}

	UiListConfigsManager interface {
		DisableListConfigsBox()
		EnableListConfigsBox()
		SelectedCfgFromListConfigs() *string
	}

	SysTrayIconsManager interface {
		SetIcon(path string)
		SetDisconnectIcon() error
		SetConnectIcon() error
		SetOpenIcon() error
		SetBlinkIcon() error
	}

	SysTrayImagesManager interface {
		SearchKeyInMap(s string) (*string, error)
		Image() map[string][]byte
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
