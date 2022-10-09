package usecase

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type VpnUseCase struct {
	sessionSetters       SessionSetters
	sessionManager       SessionManager
	cfgSetters           ConfigSetters
	cfgBody              ConfigBody
	cfgRemover           ConfigRemover
	cfgCheck             ConfigChecker
	cfgMerg              ConfigMerger
	cfgTools             ConfigTools
	uiLoger              UiLoger
	uiLogFormManager     UiLogFormManager
	uiButtonsManager     UiButtonsManager
	uiListConfigsManager UiListConfigsManager
	sysTrayIconsManager  SysTrayIconsManager
	fileSetters          FileSetters
	fileReader           FileReader
	dnsSetters           DnsSetters
}

func NewVpnUseCase(

	sessionSetters SessionSetters,
	sessionManager SessionManager,
	cfgSetters ConfigSetters,
	cfgBody ConfigBody,
	cfgRemover ConfigRemover,
	cfgCheck ConfigChecker,
	cfgMerg ConfigMerger,
	cfgTools ConfigTools,
	uiLoger UiLoger,
	uiLogFormManager UiLogFormManager,
	uiButtonsManager UiButtonsManager,
	uiListConfigsManager UiListConfigsManager,
	sysTrayIconsManager SysTrayIconsManager,
	fileSetters FileSetters,
	fileReader FileReader,
	dnsSetters DnsSetters,
) (obj *VpnUseCase, err error) {
	obj = &VpnUseCase{
		sessionSetters:       sessionSetters,
		sessionManager:       sessionManager,
		cfgSetters:           cfgSetters,
		cfgBody:              cfgBody,
		cfgRemover:           cfgRemover,
		cfgCheck:             cfgCheck,
		cfgMerg:              cfgMerg,
		cfgTools:             cfgTools,
		uiLoger:              uiLoger,
		uiLogFormManager:     uiLogFormManager,
		uiButtonsManager:     uiButtonsManager,
		uiListConfigsManager: uiListConfigsManager,
		sysTrayIconsManager:  sysTrayIconsManager,
		fileSetters:          fileSetters,
		fileReader:           fileReader,
		dnsSetters:           dnsSetters,
	}
	return
}

func (v *VpnUseCase) GetChanVpnLog() chan string {
	logChan := v.uiLoger.ChanVpnLog()
	return logChan
}

func (vp *VpnUseCase) CaseSetLogsInTextWidget(text string) {
	if runtime.GOOS != "windows" {
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
	originalText := vp.uiLogFormManager.GetTextFromLogForm() + text + "\n"
	vp.uiLogFormManager.SetTextInLogForm(originalText)
}

func (v *VpnUseCase) OffComboBoxAndClear() {
	v.uiListConfigsManager.DisableListConfigsBox()
	v.uiLogFormManager.ClearLogForm()
}

func (v *VpnUseCase) TurnOnConfigsBox() {
	v.uiListConfigsManager.EnableListConfigsBox()
}

func (v *VpnUseCase) DisableConnectionButton() {
	v.uiButtonsManager.ButtonConnectDisable()
}

func (v *VpnUseCase) DisableDisconnectButton() {
	v.uiButtonsManager.ButtonDisconnectDisable()
}

func (v *VpnUseCase) EnableDisconnectButton() {
	v.uiButtonsManager.ButtonDisconnectEnable()
}

func (v *VpnUseCase) EnableConnectButton() {
	v.uiButtonsManager.ButtonConnectEnable()
}

func (v *VpnUseCase) TraySetImageDisconnect() error {
	err := v.sysTrayIconsManager.SetDisconnectIcon()
	if err != nil {
		return fmt.Errorf("don't set image tray: [%w]", err)
	}
	return nil
}

func (v *VpnUseCase) TraySetImageConnect() error {
	err := v.sysTrayIconsManager.SetConnectIcon()
	if err != nil {
		return fmt.Errorf("don't set image tray: [%w]", err)
	}
	return nil
}

func (v *VpnUseCase) SetPhyseInterface(i string) {
	v.dnsSetters.SetInterface(i)
}

func (vp *VpnUseCase) CaseFlickeringIcon() error {

	if runtime.GOOS != "windows" {
		time.Sleep(time.Duration(250) * time.Millisecond)
	}
	err := vp.sysTrayIconsManager.SetBlinkIcon()
	if err != nil {
		return fmt.Errorf("don't set image blink [%w]", err)
	}
	if runtime.GOOS != "windows" {
		time.Sleep(time.Duration(30) * time.Millisecond)
	}
	err = vp.sysTrayIconsManager.SetOpenIcon()
	if err != nil {
		return fmt.Errorf("don't set image open [%w]", err)
	}
	if runtime.GOOS != "windows" {
		time.Sleep(time.Duration(60) * time.Millisecond)
	}
	err = vp.sysTrayIconsManager.SetDisconnectIcon()
	if err != nil {
		return fmt.Errorf("don't set image disconnect [%w]", err)
	}
	return nil
}

func (vp *VpnUseCase) GetTextFromComboBox() string {
	return *vp.uiListConfigsManager.SelectedCfgFromListConfigs()
}

func (vp *VpnUseCase) ReadFile() ([]byte, error) {
	body, err := vp.fileReader.ReadFileAsByte()
	if err != nil {
		return nil, fmt.Errorf("don't read file: [%w]", err)
	}
	return body, nil
}

func (vp *VpnUseCase) GetOvpnAuthPathFileName() string {
	return vp.cfgTools.GetAuthpathFileName()
}

func (vp *VpnUseCase) SetProfileBody(profileBody string) {
	vp.sessionSetters.SetConfig(profileBody)
	vp.cfgSetters.SetBody(profileBody)
}

func (vp *VpnUseCase) CheckOvpnUseAuthUserPass() bool {
	return vp.cfgCheck.CheckStringAuthUserPass()
}

func (vp *VpnUseCase) RunSession(ctx context.Context) {
	vp.sessionManager.StartSession(ctx)
}

func (vp *VpnUseCase) NewSession() {
	vp.sessionSetters.SetSession()
}

func (vp *VpnUseCase) SetSessionCread(u, p string) {
	vp.sessionSetters.SetCread(u, p)
}

func (vp *VpnUseCase) SetPathToFile(path string) {
	vp.fileSetters.SetPath(path)
}

func (vp *VpnUseCase) SetBodyToCfg(path string) {
	vp.cfgSetters.SetBody(path)
}

func (vp *VpnUseCase) GetVpnCread() (string, string) {
	return vp.cfgTools.GetUserAndPass()
}
