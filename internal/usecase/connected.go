package usecase

import (
	"os"
	"sync"
)

var Wg sync.WaitGroup

func (vp *VpnUseCase) Connect(ConfigsPath, Level string, StopTimeout int, UseSystemd bool) func() {

	f := func() {

		Wg.Add(1)
		go vp.ReadLogsAndStartProcess(&UseSystemd, &Wg)
		vp.UiRepo.DisableComboBox()
		vp.UiRepo.ClearTextEdit()
		vp.UiRepo.ButtonConnectDisable()
		cfg := ConfigsPath + *vp.UiRepo.SelectedFromComboBox()
		profile := vp.MemoryRepo.GetProfile(cfg)

		if profile.Body == "" {
			for _, f := range [2]string{cfg + ".ovpn", cfg + ".conf"} {
				vp.FileRepo.SetPath(f)
				if vp.FileRepo.CheckFileExists() {
					break
				}
			}
			body, err := vp.FileRepo.ReadFileAsByte()
			if err != nil {
				vp.LogRepo.Fatalf("Don't read file: %s", err)
			}
			vp.FileRepo.SetBody(body)
			s := string(body)
			vp.GlueRepo.SetBody(s)
			vp.GlueRepo.RemoveSpaceLines()
			vp.GlueRepo.RemoveCommentLines()
			vp.GlueRepo.RemoveEmptyString()
			vp.MemoryRepo.Save(cfg, vp.GlueRepo.GetBody())
		}
		os.Chdir(ConfigsPath)
		if vp.GlueRepo.CheckConfigUseFiles() {
			if profile.Body == "" {
				infile, err := vp.FileRepo.FileOpen()
				if err != nil {
					vp.LogRepo.Fatal(err)
				}
				profileBody := vp.GlueRepo.AddStringToConfig(infile)
				profile.Body = profileBody
			}
			vp.caseRunOpenvpnWithFiles(profile.Body, StopTimeout, UseSystemd)
		} else {
			if profile.Body == "" {
				profileBody, err := vp.FileRepo.ReadFileAsString()
				if err != nil {
					vp.LogRepo.Fatal(err)
				}
				profile.Body = *profileBody
			}
			vp.caseRunOpenvpn(profile.Body, StopTimeout, UseSystemd)
		}
	}
	return f

}
