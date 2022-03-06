package usecase

import (
	"io/fs"
	"os"
	"os/exec"

	"github.com/Vai3soh/goovpn/entity"
)

type (
	LoggerInteractor interface {
		Fatal(...interface{})
		Debugf(string, ...interface{})
		Fatalf(string, ...interface{})
	}

	ManagerInteractor interface {
		Connect(ConfigsPath, Level string, StopTimeout int, UseSystemd bool) func()
		Disconnect(StopTimeout int, UseSystemd bool) func()
	}

	GlueConfigInteractor interface {
		SetBody(body string)
		SetPath(path string)
		GetBody() string
		RemoveSpaceLines()
		RemoveCommentLines()
		RemoveEmptyString()
		RemoveNotCertsAndKeys()
		RemoveCertsAndKeys()
		CheckConfigUseFiles() bool
		AddStringToConfig(inFile *os.File) string
		SearchFilesPaths() map[string]string
		MergeCertsAndKeys(key string) string
		GetAuthpathFileName() string
		GetUserAndPass() (string, string)
		CheckStringAuthUserPass() bool
	}

	SessionInteractor interface {
		SetConfig(config string)
		SetSession(username, password string)
		StartSession()
		StopSession()
		StopSessionWithTimeout(timeout int)
	}

	CloseAppInteractor interface {
		CloseApp()
		SetBind(bind func())
	}

	UiInteractor interface {
		DisableComboBox()
		EnableComboBox()
		ButtonConnectEnable()
		ButtonConnectDisable()
		ButtonDisconnectEnable()
		ButtonDisconnectDisable()
		GetTextFromTextEdit() string
		SetTextInTextEdit(text string)
		SelectedFromComboBox() *string
		ClearTextEdit()
		ChanVpnLog() chan string
		CloseChanVpnLog()
		Log(text string)
	}

	SysTrayInteractor interface {
		SetIcon(path string)
		SetDisconnectIcon() error
		SetConnectIcon() error
		SetOpenIcon() error
		SetBlinkIcon() error
		SearchKeyInMap(s string) (*string, error)
		Image() map[string][]byte
	}

	FileInteractor interface {
		FileNameWithoutExtension() *string
		ReadFileAsByte() ([]byte, error)
		ReadFileAsString() (*string, error)
		SetBody([]byte)
		SetPath(path string)
		Path() string
		FileOpen() (*os.File, error)
		SetPermissonFile(fs.FileMode)
		WriteByteFile() error
		CreateFile() (*os.File, error)
		WriteStringToFile(file *os.File, data string) error
		AbsolutePath() (*string, error)
		SetDestPath(destPath string)
		CopyFile() error
		Body() []byte
		CheckFileExists() bool
	}

	DnsInteractor interface {
		CmdSystemdResolv() (*string, error)
		CmdResolvConf() (*string, *string, error)
		CmdDownResolvConf() (*string, error)
		SetAddress(addr []string)
		SetInterface(iface string)
	}

	CommandInteractor interface {
		SetCommand(command string)
		SplitCmd() ([]string, error)
		PassArgumentsToExec([]string) *exec.Cmd
		SetToProc(*exec.Cmd)
		StartProc() error
		Proc() *exec.Cmd
		RunCmdWithPipe(args1, args2 []string) error
	}

	MemoryInteractor interface {
		Save(cfgPath, body string)
		GetProfile(cfgPath string) entity.Profile
	}
)
