package monitor

func StartMonitot() {
	go StartClusterStatusM()
}
