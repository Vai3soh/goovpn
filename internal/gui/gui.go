package gui

import (
	"fmt"
	"strings"

	"github.com/Vai3soh/goovpn/internal/usecase"
	qtgui "github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Manager struct {
	usecase.ManagerInteractor
}

type sysTray struct {
	Tray  *widgets.QSystemTrayIcon
	menu  *widgets.QMenu
	image map[string][]byte
}

type uiMainWindow struct {
	mainWindow            *widgets.QMainWindow
	centralwidget         *widgets.QWidget
	gridLayout            *widgets.QGridLayout
	horizontalLayout      *widgets.QHBoxLayout
	PushButtonClear       *widgets.QPushButton
	PushButtonConnect     *widgets.QPushButton
	PushButtonDiscconnect *widgets.QPushButton
	PushButtonExit        *widgets.QPushButton
	verticalLayout        *widgets.QVBoxLayout
	textEditReadOnly      *widgets.QTextEdit
	comboBox              *widgets.QComboBox

	vpnLog chan string
}

type Option func(*uiMainWindow)

func NewUiMainWindow(opts ...Option) *uiMainWindow {
	f := &uiMainWindow{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithApp(mainWindow *widgets.QMainWindow) Option {
	return func(ui *uiMainWindow) {
		ui.mainWindow = mainWindow
	}
}

func WithCentralwidget(centralwidget *widgets.QWidget) Option {
	return func(ui *uiMainWindow) {
		ui.centralwidget = centralwidget
	}
}

func WithGridLayout(gridLayout *widgets.QGridLayout) Option {
	return func(ui *uiMainWindow) {
		ui.gridLayout = gridLayout
	}
}

func WithHorizontalLayout(horizontalLayout *widgets.QHBoxLayout) Option {
	return func(ui *uiMainWindow) {
		ui.horizontalLayout = horizontalLayout
	}
}

func WithPushButtonClear(pushButtonClear *widgets.QPushButton) Option {
	return func(ui *uiMainWindow) {
		ui.PushButtonClear = pushButtonClear
	}
}

func WithPushButtonConnect(pushButtonConnect *widgets.QPushButton) Option {
	return func(ui *uiMainWindow) {
		ui.PushButtonConnect = pushButtonConnect
	}
}

func WithPushButtonDisconnect(pushButtonDisconnect *widgets.QPushButton) Option {
	return func(ui *uiMainWindow) {
		ui.PushButtonDiscconnect = pushButtonDisconnect
	}
}

func WithPushButtonExit(pushButtonExit *widgets.QPushButton) Option {
	return func(ui *uiMainWindow) {
		ui.PushButtonExit = pushButtonExit
	}
}

func WithVerticalLayout(verticalLayout *widgets.QVBoxLayout) Option {
	return func(ui *uiMainWindow) {
		ui.verticalLayout = verticalLayout
	}
}

func WithTextEditReadOnly(textEditReadOnly *widgets.QTextEdit) Option {
	return func(ui *uiMainWindow) {
		ui.textEditReadOnly = textEditReadOnly
	}
}

func WithComboBox(comboBox *widgets.QComboBox) Option {
	return func(ui *uiMainWindow) {
		ui.comboBox = comboBox
	}
}

func WithChanVpnLog(vpnLog *chan string) Option {
	return func(g *uiMainWindow) {
		g.vpnLog = *vpnLog
	}
}

type OptionTray func(*sysTray)

func NewSysTray(opts ...OptionTray) *sysTray {
	s := &sysTray{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithImage(image map[string][]byte) OptionTray {
	return func(s *sysTray) {
		s.image = image
	}
}

func WithSystemTrayIcon(t *widgets.QSystemTrayIcon) OptionTray {
	return func(s *sysTray) {
		s.Tray = t
	}
}

func WithSystemTrayMenu(m *widgets.QMenu) OptionTray {
	return func(s *sysTray) {
		s.menu = m
	}
}

func (ui *uiMainWindow) SetupUI(
	app *widgets.QApplication, configsPath,
	level, tempDir, appName, appPath string, configsName []string,
) {

	ui.mainWindow.SetWindowTitle(appName)
	ui.mainWindow.SetGeometry2(100, 100, 600, 400)
	ui.mainWindow.SetWindowIcon(qtgui.NewQIcon5(appPath))
	ui.horizontalLayout.AddWidget(ui.PushButtonClear, 1, 0)
	ui.horizontalLayout.AddWidget(ui.PushButtonConnect, 1, 0)
	ui.horizontalLayout.AddWidget(ui.PushButtonDiscconnect, 1, 0)
	ui.horizontalLayout.AddWidget(ui.PushButtonExit, 1, 0)
	ui.gridLayout.AddLayout(ui.horizontalLayout, 2, 0, 1)
	ui.textEditReadOnly.SetReadOnly(true)
	ui.verticalLayout.AddWidget(ui.textEditReadOnly, 0, 0)
	ui.gridLayout.AddLayout(ui.verticalLayout, 1, 0, 1)

	ui.gridLayout.AddWidget2(ui.comboBox, 0, 0, 0)
	ui.comboBox.AddItems(configsName)
	ui.mainWindow.SetCentralWidget(ui.centralwidget)

}

func (ui *uiMainWindow) DisableComboBox() {
	ui.comboBox.SetDisabled(true)
}

func (ui *uiMainWindow) EnableComboBox() {
	ui.comboBox.SetEnabled(true)
}

func (ui *uiMainWindow) ButtonConnectDisable() {
	ui.PushButtonConnect.SetDisabled(true)
}

func (ui *uiMainWindow) ButtonConnectEnable() {
	ui.PushButtonConnect.SetEnabled(true)
}

func (ui *uiMainWindow) ButtonDisconnectDisable() {
	ui.PushButtonDiscconnect.SetDisabled(true)
}

func (ui *uiMainWindow) ButtonDisconnectEnable() {
	ui.PushButtonDiscconnect.SetEnabled(true)
}

func (ui *uiMainWindow) SelectedFromComboBox() *string {
	s := ui.comboBox.CurrentText()
	return &s
}

func (ui *uiMainWindow) SetTextInTextEdit(text string) {
	ui.textEditReadOnly.SetPlainText(text)
}

func (ui *uiMainWindow) GetTextFromTextEdit() string {
	return ui.textEditReadOnly.ToPlainText()
}

func (ui *uiMainWindow) ClearTextEdit() {
	ui.textEditReadOnly.Clear()
}

func (s *sysTray) SetupSysTray() (*widgets.QAction, *widgets.QAction, error) {
	path, err := s.SearchKeyInMap("disconnect")
	if err != nil {
		return nil, nil, err
	}
	s.SetIcon(*path)
	main := s.menu.AddAction("Open main window")
	exit := s.menu.AddAction("Exit")
	return exit, main, nil
}

func (s *sysTray) Image() map[string][]byte {
	return s.image
}

func (s *sysTray) SearchKeyInMap(st string) (*string, error) {
	for key := range s.Image() {
		if strings.Contains(key, st) {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("key in map not found?")
}

func (g *uiMainWindow) ChanVpnLog() chan string {
	return g.vpnLog
}

func (g *uiMainWindow) CloseChanVpnLog() {
	close(g.vpnLog)
}

func (g *uiMainWindow) Log(text string) {
	if text != "" {
		g.vpnLog <- text
	}
}

func (s *sysTray) SetIcon(path string) {
	s.Tray.SetIcon(qtgui.NewQIcon5(path))
}

func (s *sysTray) SetConnectIcon() error {
	path, err := s.SearchKeyInMap("connecting")
	if err != nil {
		return err
	}
	s.SetIcon(*path)
	return nil
}

func (s *sysTray) SetDisconnectIcon() error {
	path, err := s.SearchKeyInMap("disconnect")
	if err != nil {
		return err
	}
	s.SetIcon(*path)
	return nil
}

func (s *sysTray) SetOpenIcon() error {
	path, err := s.SearchKeyInMap("open")
	if err != nil {
		return err
	}

	s.SetIcon(*path)
	return nil
}

func (s *sysTray) SetBlinkIcon() error {
	path, err := s.SearchKeyInMap("blink")
	if err != nil {
		return err
	}
	s.SetIcon(*path)
	return nil
}
