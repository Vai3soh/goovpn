package usecase

func (vp *VpnUseCase) ReturnResolvConf() {
	cmd, err := vp.DnsRepo.CmdDownResolvConf()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	vp.CmdRepo.SetCommand(*cmd)
	cmdArg, err := vp.CmdRepo.SplitCmd()
	if err != nil {
		vp.LogRepo.Fatal(err)
	}
	cmdExec := vp.CmdRepo.PassArgumentsToExec(cmdArg)
	vp.CmdRepo.SetToProc(cmdExec)
	vp.CmdRepo.StartProc()
}

func (vp *VpnUseCase) CloserSession(time int, UseSystemd bool) func() {
	return func() {
		vp.SessionRepo.StopSessionWithTimeout(time)
		vp.UiRepo.EnableComboBox()
		vp.UiRepo.ButtonDisconnectDisable()
		vp.UiRepo.ButtonConnectEnable()
		if !UseSystemd {
			vp.ReturnResolvConf()
		}
	}
}

func (vp *VpnUseCase) Disconnect(StopTimeout int, UseSystemd bool) func() {
	return vp.CloserSession(StopTimeout, UseSystemd)
}
