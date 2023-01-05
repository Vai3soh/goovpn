package usecase

import (
	"fmt"
)

type VpnUseCase struct {
	sessionSetters SessionSetters
	sessionManager SessionManager
	cfgSetters     ConfigSetters
	cfgBody        ConfigBody
	cfgRemover     ConfigRemover
	cfgCheck       ConfigChecker
	cfgMerg        ConfigMerger
	cfgTools       ConfigTools
	sessionLoger   SessionLoger
	fileSetters    FileSetters
	fileReader     FileReader
	dnsSetters     DnsSetters
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
	sessionLoger SessionLoger,
	fileSetters FileSetters,
	fileReader FileReader,
	dnsSetters DnsSetters,
) (obj *VpnUseCase, err error) {
	obj = &VpnUseCase{
		sessionSetters: sessionSetters,
		sessionManager: sessionManager,
		cfgSetters:     cfgSetters,
		cfgBody:        cfgBody,
		cfgRemover:     cfgRemover,
		cfgCheck:       cfgCheck,
		cfgMerg:        cfgMerg,
		cfgTools:       cfgTools,
		sessionLoger:   sessionLoger,
		fileSetters:    fileSetters,
		fileReader:     fileReader,
		dnsSetters:     dnsSetters,
	}
	return
}

func (v *VpnUseCase) GetChanVpnLog() chan string {
	logChan := v.sessionLoger.ChanVpnLog()
	return logChan
}

func (v *VpnUseCase) SetPhyseInterface(i string) {
	v.dnsSetters.SetInterface(i)
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

func (vp *VpnUseCase) RunSession() error {
	return vp.sessionManager.StartSession()
}

func (vp *VpnUseCase) DestroyVpnClient() {
	vp.sessionManager.DestroyClient()
}

func (vp *VpnUseCase) ExitSession() {
	vp.sessionManager.StopSession()
}

func (vp *VpnUseCase) SetSessionCread(u, p string) error {
	return vp.sessionSetters.SetCread(u, p)
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
