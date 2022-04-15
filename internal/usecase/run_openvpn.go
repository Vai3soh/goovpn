package usecase

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func (vp *VpnUseCase) caseSetupDnsWithUseSystemd() {
	cmdResolv, err := vp.DnsRepo.CmdSystemdResolv()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	vp.CmdRepo.SetCommand(*cmdResolv)
	cmdArg, err := vp.CmdRepo.SplitCmd()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	cmdExec := vp.CmdRepo.PassArgumentsToExec(cmdArg)
	vp.CmdRepo.SetToProc(cmdExec)
	vp.CmdRepo.StartProc()
}

func (vp *VpnUseCase) caseSetupDnsNotUseSystemd() {
	cmdPrintf, cmdResolv, err := vp.DnsRepo.CmdResolvConf()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	vp.CmdRepo.SetCommand(*cmdPrintf)
	cmdArg1, err := vp.CmdRepo.SplitCmd()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	vp.CmdRepo.SetCommand(*cmdResolv)
	cmdArg2, err := vp.CmdRepo.SplitCmd()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	if err := vp.CmdRepo.RunCmdWithPipe(cmdArg1, cmdArg2); err != nil {
		vp.LogRepo.Fatal(err)
	}
}

func (vp *VpnUseCase) caseAddDnsAddress(text string) {
	reg := regexp.MustCompile(`(?m)DNS Servers:\s+[^*]+\d$`)
	matches := reg.FindAllString(text, -1)
	matches = strings.Split(matches[0], "\n")
	matches = append(matches[:0], matches[0+1:]...)
	for i := range matches {
		matches[i] = strings.TrimSpace(matches[i])
	}
	vp.DnsRepo.SetAddress(matches)
}

func (vp *VpnUseCase) caseFlickeringIcon(text string) {

	time.Sleep(time.Duration(250) * time.Millisecond)
	err := vp.TrayRepo.SetBlinkIcon()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	time.Sleep(time.Duration(30) * time.Millisecond)
	err = vp.TrayRepo.SetOpenIcon()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	time.Sleep(time.Duration(60) * time.Millisecond)
	err = vp.TrayRepo.SetDisconnectIcon()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
}

func (vp *VpnUseCase) caseOpenvpnSessionEnd(UseSystemd *bool) {
	vp.UiRepo.ButtonConnectEnable()
	vp.UiRepo.EnableComboBox()
	if !*UseSystemd {
		vp.ReturnResolvConf()
	}
}

func (vp *VpnUseCase) caseSetLogsInTextWidget(text string) {
	time.Sleep(time.Duration(50) * time.Millisecond)
	originalText := vp.UiRepo.GetTextFromTextEdit() + text + "\n"
	vp.UiRepo.SetTextInTextEdit(originalText)
}

func (vp *VpnUseCase) ReadLogsAndStartProcess(UseSystemd *bool, wg *sync.WaitGroup) {

	logChan := vp.UiRepo.ChanVpnLog()
	defer wg.Done()
	iterSkip := false
	for text := range logChan {
		text += "\n"
		vp.caseSetLogsInTextWidget(text)

		if strings.Contains(text, "DNS Servers:") {
			vp.caseAddDnsAddress(text)
		}
		if strings.Contains(text, "Openvpn3 session ended") {
			vp.caseOpenvpnSessionEnd(UseSystemd)
			vp.UiRepo.ButtonDisconnectDisable()
			err := vp.TrayRepo.SetDisconnectIcon()
			if err != nil {
				vp.LogRepo.Fatal(err)
			}
			break
		}
		if strings.Contains(text, "Connected via") {
			err := vp.TrayRepo.SetConnectIcon()
			if err != nil {
				vp.LogRepo.Fatal(err)
			}
			continue
		}
		if iterSkip {
			continue
		}
		if regexp.MustCompile(`\w+[0-9]\s+opened|\w+\s+opened`).MatchString(text) {

			r := regexp.MustCompile(`(?P<int>\w+[0-9]|\w+)`)
			iface := r.FindStringSubmatch(text)

			vp.DnsRepo.SetInterface(iface[0])
			if !*UseSystemd {
				vp.caseSetupDnsNotUseSystemd()
				err := vp.TrayRepo.SetConnectIcon()
				if err != nil {
					vp.LogRepo.Fatal(err)
				}
				iterSkip = true
				continue
			} else {
				vp.caseSetupDnsWithUseSystemd()
				err := vp.TrayRepo.SetConnectIcon()
				if err != nil {
					vp.LogRepo.Fatal(err)
				}
				iterSkip = true
				continue
			}
		} else {
			vp.caseFlickeringIcon(text)
		}
	}
}

func (vp *VpnUseCase) caseRunOpenvpnWithFiles(profileBody string, time int, UseSystemd bool) {

	profileCurrent := profileBody
	vp.GlueRepo.SetBody(profileBody)
	vp.GlueRepo.RemoveNotCertsAndKeys()

	filesMap := vp.GlueRepo.SearchFilesPaths()
	for key, file := range filesMap {
		file = strings.Trim(file, `'`)
		file = strings.Trim(file, `"`)
		vp.FileRepo.SetPath(file)
		fileAbs, _ := vp.FileRepo.AbsolutePath()
		vp.FileRepo.SetPath(*fileAbs)
		s, _ := vp.FileRepo.ReadFileAsString()
		vp.GlueRepo.SetBody(*s)
		merged := vp.GlueRepo.MergeCertsAndKeys(key)
		vp.GlueRepo.SetBody(profileCurrent)
		vp.GlueRepo.RemoveCertsAndKeys()
		profileCurrent = strings.Trim(profileCurrent, "\n") + merged
	}

	vp.SessionRepo.SetConfig(profileCurrent)
	vp.GlueRepo.SetBody(profileCurrent)
	vp.NewSession()
	vp.SessionRepo.StartSession()
	vp.UiRepo.ButtonDisconnectEnable()
}

func (vp *VpnUseCase) NewSession() {
	ok := vp.GlueRepo.CheckStringAuthUserPass()
	username, password := "", ""
	vp.SessionRepo.SetSession(username, password)
	if ok {
		fileAuth := vp.GlueRepo.GetAuthpathFileName()
		if fileAuth == "" {
			vp.LogRepo.Fatal(`not get path from directive auth-user-pass, edit config - auth-user-pass auth.txt`)
		}
		vp.FileRepo.SetPath(fileAuth)
		absPathCredFile, err := vp.FileRepo.AbsolutePath()
		if err != nil {
			vp.LogRepo.Fatal(err)
		}
		vp.FileRepo.SetPath(*absPathCredFile)
		CredFileBody, err := vp.FileRepo.ReadFileAsString()
		if err != nil {
			vp.LogRepo.Fatal(err)
		}
		vp.GlueRepo.SetBody(*CredFileBody)
		username, password := vp.GlueRepo.GetUserAndPass()
		vp.SessionRepo.SetSession(username, password)
	}
}

func (vp *VpnUseCase) CreateDir() {
	if _, err := os.Stat(vp.FileRepo.Path()); os.IsNotExist(err) {
		os.MkdirAll(vp.FileRepo.Path(), 0755)
	}
}

func (vp *VpnUseCase) CopyImages() {

	mode := os.FileMode(int(0644))
	for key, value := range vp.TrayRepo.Image() {
		vp.FileRepo.SetPath(key)
		vp.FileRepo.SetBody(value)
		vp.FileRepo.SetPermissonFile(mode)
		vp.FileRepo.WriteByteFile()
	}
}

func (vp *VpnUseCase) caseRunOpenvpn(profileBody string, time int, UseSystemd bool) {
	vp.SessionRepo.SetConfig(profileBody)
	vp.GlueRepo.SetBody(profileBody)
	vp.NewSession()
	vp.SessionRepo.StartSession()
	vp.UiRepo.ButtonDisconnectEnable()
}
