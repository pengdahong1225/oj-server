package signal

// syscall.SIGUSR1 = 10
// syscall.SIGUSR2 = 12
func SignalListen(cb func()) {
	//sigs := make(chan os.Signal)
	//
	//signal.Notify(sigs, syscall.SIGUSR1)
	//for {
	//	sig := <-sigs
	//	if int(sig.(syscall.Signal)) == 10 {
	//		fmt.Println("User has logged in (SIGUSR1(10) received).")
	//		cb()
	//	}
	//}

	// signal.Notify(sigs, syscall.SIGALRM)
	// for {
	// 	sig := <-sigs
	// 	if sig == syscall.SIGALRM {
	// 		fmt.Println("User has logged in (SIGALRM received).")
	// 		cb()
	// 	}
	// }
}
