package app

import (
	"context"
	"runtime"

	"github.com/Vai3soh/goovpn/config"
	"github.com/Vai3soh/goovpn/embedfile"
	"github.com/Vai3soh/goovpn/entity"
	"github.com/Vai3soh/goovpn/internal/adapters/db/memory"
	"github.com/Vai3soh/goovpn/internal/cli"
	"github.com/Vai3soh/goovpn/internal/close"

	"github.com/Vai3soh/goovpn/internal/dns"
	"github.com/Vai3soh/goovpn/internal/fileextended"
	"github.com/Vai3soh/goovpn/internal/gui"
	"github.com/Vai3soh/goovpn/internal/parser"
	"github.com/Vai3soh/goovpn/internal/session"
	transport "github.com/Vai3soh/goovpn/internal/transport/openvpn"
	"github.com/Vai3soh/goovpn/internal/usecase"
	"github.com/Vai3soh/goovpn/internal/usecase/usecasedns"
	"github.com/Vai3soh/goovpn/internal/usecase/usecaseprofile"
	"github.com/Vai3soh/goovpn/pkg/logger"

	"os"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type fileEmbbed interface {
	usecase.FileSetters
	usecase.FileGetters
	usecase.FileToolsManager
	usecase.FileWriter
}

func CreateDir(path string, file fileEmbbed) {
	file.SetPath(path)
	if _, err := os.Stat(file.Path()); os.IsNotExist(err) {
		os.MkdirAll(file.Path(), 0755)
	}
}

func getFiles(l *logger.Logger, file *fileextended.File, path string) []string {
	files, err := file.FilesInDir(path)
	if err != nil {
		l.Fatalf("don't read dir: %s\n", err)
	}
	return files
}

func CopyImages(file fileEmbbed, stray usecase.SysTrayImagesManager) {

	mode := os.FileMode(int(0644))
	for key, value := range stray.Image() {
		file.SetPath(key)
		file.SetBody(value)
		file.SetPermissonFile(mode)
		file.WriteByteFile()
	}
}

func readEmbed(path string, l *logger.Logger) []byte {
	dataImage, err := embedfile.ReadFs(path)
	if err != nil {
		l.Fatalf("don't read png file: %s", err)
	}
	return dataImage
}

func Run(cfg *config.Config) {

	l := logger.NewLogger(
		logger.WithLogTextFormatter(), logger.WithLogLevel(&cfg.Log.Level),
	)

	ImageMap := make(map[string][]byte)
	dataPng := readEmbed(cfg.AppImagePathConnected, l)
	ImageMap[cfg.TempDir+strings.Split(cfg.AppImagePathConnected, "/")[1]] = dataPng

	dataPng = readEmbed(cfg.AppImagePathDisconnected, l)
	ImageMap[cfg.TempDir+strings.Split(cfg.AppImagePathDisconnected, "/")[1]] = dataPng

	dataPng = readEmbed(cfg.AppImagePathBlink, l)
	ImageMap[cfg.TempDir+strings.Split(cfg.AppImagePathBlink, "/")[1]] = dataPng

	dataPng = readEmbed(cfg.AppImagePathOpen, l)
	ImageMap[cfg.TempDir+strings.Split(cfg.AppImagePathOpen, "/")[1]] = dataPng

	dataPng = readEmbed(cfg.AppIcon, l)
	ImageMap[cfg.TempDir+strings.Split(cfg.AppIcon, "/")[1]] = dataPng

	var files []string
	file := fileextended.NewFile()
	if runtime.GOOS == "windows" {
		path := os.Getenv(`USERPROFILE`)
		files = getFiles(l, file, path+"\\"+cfg.ConfigsPath)
	} else {
		files = getFiles(l, file, cfg.ConfigsPath)
	}
	command := cli.NewCmd()
	cmdResolver := cli.NewResolver(cli.WithCliResolver(command))

	logVpn := make(chan string)
	app := widgets.NewQApplication(len(os.Args), os.Args)
	mainWindown := widgets.NewQMainWindow(nil, 0)
	centralwidget := widgets.NewQWidget(mainWindown, core.Qt__Widget)
	comboBox := widgets.NewQComboBox(centralwidget)
	gridLayout := widgets.NewQGridLayout(centralwidget)
	horizontalLayout := widgets.NewQHBoxLayout()
	pushButtonClear := widgets.NewQPushButton2("Clear", centralwidget)
	pushButtonConnect := widgets.NewQPushButton2("Connect", centralwidget)
	pushButtonDisconnect := widgets.NewQPushButton2("Disconnect", centralwidget)
	pushButtonExit := widgets.NewQPushButton2("Exit", centralwidget)
	textEditReadOnly := widgets.NewQTextEdit(centralwidget)
	verticalLayout := widgets.NewQVBoxLayout()

	mainUiWindow := gui.NewUiMainWindow(

		gui.WithApp(mainWindown),
		gui.WithCentralwidget(centralwidget),
		gui.WithComboBox(comboBox),
		gui.WithGridLayout(gridLayout),
		gui.WithHorizontalLayout(horizontalLayout),
		gui.WithPushButtonClear(pushButtonClear),
		gui.WithPushButtonConnect(pushButtonConnect),
		gui.WithPushButtonDisconnect(pushButtonDisconnect),
		gui.WithPushButtonExit(pushButtonExit),
		gui.WithTextEditReadOnly(textEditReadOnly),
		gui.WithVerticalLayout(verticalLayout),
		gui.WithChanVpnLog(&logVpn),
	)
	trayIcon := widgets.NewQSystemTrayIcon(nil)
	menu := widgets.NewQMenu(nil)

	stray := gui.NewSysTray(
		gui.WithSystemTrayIcon(trayIcon),
		gui.WithSystemTrayMenu(menu),
		gui.WithImage(ImageMap),
	)

	gl := parser.NewConfig()
	cl := &close.ShutdownApp{}

	sessOvpn := session.NewOpenvpnClient(
		session.WithCompressionMode(cfg.App.CompressionMode),
		session.WithDisableClientCert(cfg.App.CheckDisableClientCert),
		session.WithTimeout(cfg.App.ConnectTimeout),
		session.WithVerboseLog(cfg.App.VerbLogs),
		session.WithClockTicks(cfg.App.ClockTicks),
		session.WithTunPersist(cfg.App.TunPersist),
		session.WithUi(mainUiWindow),
	)

	memory := memory.NewDb(memory.WithMapMemory(
		make(map[string]entity.Profile)),
	)

	profileUseCase, err := usecaseprofile.NewProfileUseCase(
		file, file, file, gl, gl, gl, gl, gl, gl, memory,
	)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	dnsUseCase, err := usecasedns.NewDnsUseCase(
		command, command, cmdResolver, command, cmdResolver,
	)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	vpnUseCase, _ := usecase.NewVpnUseCase(
		sessOvpn, sessOvpn, gl, gl, gl, gl, gl, gl,
		mainUiWindow, mainUiWindow, mainUiWindow, mainUiWindow,
		stray, file, file, cmdResolver,
	)

	sys, err := dns.NewSystem(runtime.GOOS, dnsUseCase, dnsUseCase, cfg.App.UseSystemd)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}

	names, err := dns.NewNames(dnsUseCase)
	if err != nil {
		l.Fatalf("don't get constructor [%w]\n", err)
	}
	names.SetGoos(*sys)

	var tr *transport.TransportOvpnClient
	if runtime.GOOS != "windows" {
		tr = transport.New(
			cfg.App.ConfigsPath, cfg.StopTimeout,
			vpnUseCase, profileUseCase, names,
		)
	} else {
		path := os.Getenv(`USERPROFILE`)
		tr = transport.New(
			path+"\\"+cfg.ConfigsPath, cfg.StopTimeout,
			vpnUseCase, profileUseCase, names,
		)
	}

	appIconPath, err := stray.SearchKeyInMap("app")
	if err != nil {
		l.Fatalf("don't found key in map [%w]\n", err)
	}
	mainUiWindow.SetupUI(
		app, cfg.ConfigsPath,
		cfg.Log.Level, cfg.TempDir,
		cfg.Name,
		*appIconPath,
		files,
	)

	exit, main, updateCfgs, err := stray.SetupSysTray()
	if err != nil {
		l.Fatalf("setup systray failed: [%w]\n", err)
	}
	mainWindown.SetWindowFlags(core.Qt__Dialog)
	app.SetQuitOnLastWindowClosed(false)

	main.ConnectTriggered(func(bool) { mainWindown.Show() })
	updateCfgs.ConnectTriggered(func(bool) {
		if mainUiWindow.IsEnableCombo() {
			files, err := file.FilesInDir(cfg.ConfigsPath)
			if err != nil {
				l.Fatalf("don't read dir: [%s]\n", err)
			}
			mainUiWindow.UpdateComboBox(files)
		}
	})

	CreateDir(cfg.TempDir, file)
	CopyImages(file, stray)

	path, err := stray.SearchKeyInMap(`disconnect`)
	if err != nil {
		l.Fatalf("don't found key in map: [%w]\n", err)
	}
	stray.SetIcon(*path)

	trayIcon.SetContextMenu(menu)
	trayIcon.Show()
	mainWindown.Show()

	mainUiWindow.PushButtonClear.ConnectClicked(func(_ bool) {
		mainUiWindow.ClearLogForm()
	})

	mainUiWindow.PushButtonConnect.ConnectClicked(func(_ bool) {

		ctx, cancel := context.WithCancel(context.Background())

		cl.SetBind(tr.Disconnect(cancel))
		cl.Binder()

		err := tr.Connect(ctx)()
		if err != nil {
			l.Fatalf("connect failed: [%w]\n", err)
		}

		mainUiWindow.PushButtonDiscconnect.ConnectClicked(func(_ bool) {
			mainUiWindow.ButtonDisconnectDisable()
			cancel()
		})

		exit.ConnectTriggered(func(bool) {
			tr.Disconnect(cancel)()
			app.Exit(0)
		})

		mainUiWindow.PushButtonExit.ConnectClicked(func(_ bool) {
			tr.Disconnect(cancel)()
			app.Exit(0)
		})

	})

	exit.ConnectTriggered(func(bool) {
		app.Exit(0)
	})

	mainUiWindow.PushButtonExit.ConnectClicked(func(_ bool) {
		app.Exit(0)
	})

	mainUiWindow.ButtonDisconnectDisable()
	widgets.QApplication_Exec()
}
