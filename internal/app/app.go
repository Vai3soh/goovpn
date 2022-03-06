package app

import (
	"github.com/Vai3soh/goovpn/config"
	"github.com/Vai3soh/goovpn/embedfile"
	"github.com/Vai3soh/goovpn/entity"
	"github.com/Vai3soh/goovpn/internal/cache"
	"github.com/Vai3soh/goovpn/internal/close"
	"github.com/Vai3soh/goovpn/internal/cmdextended"
	"github.com/Vai3soh/goovpn/internal/dns"
	"github.com/Vai3soh/goovpn/internal/fileextended"
	"github.com/Vai3soh/goovpn/internal/glue"
	"github.com/Vai3soh/goovpn/internal/gui"
	"github.com/Vai3soh/goovpn/internal/session"
	"github.com/Vai3soh/goovpn/internal/usecase"
	"github.com/Vai3soh/goovpn/pkg/logger"

	"os"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func readEmbed(path string, l *logger.Logger) []byte {
	dataImage, err := embedfile.ReadFs(path)
	if err != nil {
		l.Fatalf("don't read png file: %s", err)
	}
	return dataImage
}

func closeAppTrigger(s *session.Openvpn,
	manager *gui.Manager,
	cfg *config.Config,
	VpnUseCase *usecase.VpnUseCase,
	app *widgets.QApplication,
) {
	if !s.SessionIsClose() {
		manager.Disconnect(cfg.StopTimeout, cfg.App.UseSystemd)()
	}

	app.Exit(0)
	usecase.Wg.Wait()
}

func Run(cfg *config.Config) {

	l := logger.NewLogger(logger.WithLogTextFormatter(), logger.WithLogLevel(&cfg.Log.Level))

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

	file := fileextended.NewFile()
	files, err := file.FilesInDir(cfg.ConfigsPath)
	if err != nil {
		l.Fatalf("don't read dir: %s", err)
	}

	command := cmdextended.NewCmd()

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

	gl := glue.NewConfig()
	cl := &close.ShutdownApp{}

	session := session.NewOpenvpn(
		session.WithCompressionMode(cfg.App.CompressionMode),
		session.WithDisableClientCert(cfg.App.CheckDisableClientCert),
		session.WithTimeout(cfg.App.ConnectTimeout),
		session.WithUi(mainUiWindow),
	)

	dns := dns.NewDns()
	memory := cache.NewDb(cache.WithMapMemory(
		make(map[string]entity.Profile)),
	)
	VpnUseCase := usecase.New(
		l, gl, cl,
		session, mainUiWindow,
		stray, file,
		command, dns, memory,
	)

	m := gui.Manager{ManagerInteractor: VpnUseCase}

	appIconPath, err := stray.SearchKeyInMap("app")
	if err != nil {
		l.Fatal(err)
	}
	mainUiWindow.SetupUI(
		app, cfg.ConfigsPath,
		cfg.Log.Level, cfg.TempDir,
		cfg.Name,
		*appIconPath,
		*files,
	)

	exit, main, err := stray.SetupSysTray()
	if err != nil {
		l.Fatal(err)
	}
	mainWindown.SetWindowFlags(core.Qt__Dialog)
	app.SetQuitOnLastWindowClosed(false)

	main.ConnectTriggered(func(bool) { mainWindown.Show() })

	VpnUseCase.FileRepo.SetPath(cfg.TempDir)
	VpnUseCase.CreateDir()
	VpnUseCase.CopyImages()

	path, err := VpnUseCase.TrayRepo.SearchKeyInMap("disconnect")
	if err != nil {
		VpnUseCase.LogRepo.Fatal(err)
	}
	VpnUseCase.TrayRepo.SetIcon(*path)

	trayIcon.SetContextMenu(menu)
	trayIcon.Show()
	mainWindown.Show()

	VpnUseCase.CloseRepo.SetBind(m.Disconnect(cfg.StopTimeout, cfg.App.UseSystemd))
	VpnUseCase.CloseRepo.CloseApp()

	mainUiWindow.PushButtonDiscconnect.ConnectClicked(func(_ bool) {
		m.Disconnect(cfg.StopTimeout, cfg.App.UseSystemd)()
	})
	mainUiWindow.PushButtonClear.ConnectClicked(func(_ bool) {
		mainUiWindow.ClearTextEdit()
	})
	mainUiWindow.PushButtonConnect.ConnectClicked(func(_ bool) {
		m.Connect(cfg.ConfigsPath, cfg.Level, cfg.StopTimeout, cfg.UseSystemd)()
	})
	mainUiWindow.PushButtonExit.ConnectClicked(func(_ bool) {
		closeAppTrigger(session, &m, cfg, VpnUseCase, app)
	})
	exit.ConnectTriggered(func(bool) {
		closeAppTrigger(session, &m, cfg, VpnUseCase, app)
	})

	stray.Tray.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		if reason == widgets.QSystemTrayIcon__Trigger {
			m.Connect(cfg.ConfigsPath, cfg.Level, cfg.StopTimeout, cfg.UseSystemd)
		}
	})

	VpnUseCase.UiRepo.ButtonDisconnectDisable()
	widgets.QApplication_Exec()
}
