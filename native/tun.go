package main

//#include "bridge.h"
import "C"
import (
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"unsafe"
	"wg/native/conn"
	"wg/native/device"
	"wg/native/tun"
)

//export startTun
func startTun(fd C.int, devicePrivateKey, listenPort, peerPublicKey, allowedIps, endpoint C.c_string, callback unsafe.Pointer) C.c_string {
	var foreground bool
	var interfaceName string

	foreground = false

	if !foreground {
		foreground = os.Getenv(ENV_WG_PROCESS_FOREGROUND) == "1"
	}

	// get log level (default: info)

	logLevel := func() int {
		switch os.Getenv("LOG_LEVEL") {
		case "verbose", "debug":
			return device.LogLevelVerbose
		case "error":
			return device.LogLevelError
		case "silent":
			return device.LogLevelSilent
		}
		return device.LogLevelError
	}()

	// open TUN device (or use supplied fd)

	tdev, err := func() (tun.Device, error) {
		//tunFdStr := os.Getenv(ENV_WG_TUN_FD)
		//if tunFdStr == "" {
		return tun.WgCreateTun(device.DefaultMTU, int(fd))
		//}

		// construct tun device from supplied fd

		//fd, err := strconv.ParseUint(tunFdStr, 10, 32)
		//if err != nil {
		//	return nil, err
		//}
		//
		//err = unix.SetNonblock(int(fd), true)
		//if err != nil {
		//	return nil, err
		//}
		//
		//file := os.NewFile(uintptr(fd), "")
		//return tun.CreateTUNFromFile(file, device.DefaultMTU)
	}()

	if err == nil {
		realInterfaceName, err2 := tdev.Name()
		if err2 == nil {
			interfaceName = realInterfaceName
		}
	}

	logger := device.NewLogger(
		logLevel,
		fmt.Sprintf("(%s) ", interfaceName),
	)

	logger.Verbosef("Starting wireguard-go version %s", Version)

	if err != nil {
		logger.Errorf("Failed to create TUN device: %v", err)
		return C.CString(err.Error())
	}

	// open UAPI file (or use supplied fd)

	//fileUAPI, err := func() (*os.File, error) {
	//	uapiFdStr := os.Getenv(ENV_WG_UAPI_FD)
	//	if uapiFdStr == "" {
	//		return ipc.UAPIOpen(interfaceName)
	//	}
	//
	//	// use supplied fd
	//
	//	fd, err := strconv.ParseUint(uapiFdStr, 10, 32)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return os.NewFile(uintptr(fd), ""), nil
	//}()
	//if err != nil {
	//	logger.Errorf("UAPI listen error: %v", err)
	//	os.Exit(ExitSetupFailed)
	//	return -1
	//}
	// daemonize the process

	//if !foreground {
	//	env := os.Environ()
	//	env = append(env, fmt.Sprintf("%s=3", ENV_WG_TUN_FD))
	//	env = append(env, fmt.Sprintf("%s=4", ENV_WG_UAPI_FD))
	//	env = append(env, fmt.Sprintf("%s=1", ENV_WG_PROCESS_FOREGROUND))
	//	files := [3]*os.File{}
	//	if os.Getenv("LOG_LEVEL") != "" && logLevel != device.LogLevelSilent {
	//		files[0], _ = os.Open(os.DevNull)
	//		files[1] = os.Stdout
	//		files[2] = os.Stderr
	//	} else {
	//		files[0], _ = os.Open(os.DevNull)
	//		files[1], _ = os.Open(os.DevNull)
	//		files[2], _ = os.Open(os.DevNull)
	//	}
	//	attr := &os.ProcAttr{
	//		Files: []*os.File{
	//			files[0], // stdin
	//			files[1], // stdout
	//			files[2], // stderr
	//			tdev.File(),
	//			fileUAPI,
	//		},
	//		Dir: ".",
	//		Env: env,
	//	}
	//
	//	path, err := os.Executable()
	//	if err != nil {
	//		logger.Errorf("Failed to determine executable: %v", err)
	//		os.Exit(ExitSetupFailed)
	//	}
	//
	//	process, err := os.StartProcess(
	//		path,
	//		os.Args,
	//		attr,
	//	)
	//	if err != nil {
	//		logger.Errorf("Failed to daemonize: %v", err)
	//		os.Exit(ExitSetupFailed)
	//	}
	//	process.Release()
	//	return -1
	//}

	device := device.NewDevice(tdev, conn.NewDefaultBind(), logger)

	logger.Verbosef("Device started")

	errs := make(chan error)
	term := make(chan os.Signal, 1)

	//uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
	//if err != nil {
	//	logger.Errorf("Failed to listen on uapi socket: %v", err)
	//	os.Exit(ExitSetupFailed)
	//}

	//go func() {
	//	for {
	//		conn, err := uapi.Accept()
	//		if err != nil {
	//			errs <- err
	//			return
	//		}
	//		go device.IpcHandle(conn)
	//	}
	//}()
	ppk:=[]byte{220,239,61,112,192,52,204,107,74,161,131,117,5,226,246,65,150,94,134,223,226,251,200,34,243,207,134,192,251,89,143,75}
	dpk:=[]byte{216,112,171,239,161,147,57,95,165,230,33,66,196,253,217,82,31,210,101,174,128,233,107,172,19,146,75,230,247,104,128,127}
	lines := []string{
		fmt.Sprintf("private_key=%s", hex.EncodeToString(dpk),
		fmt.Sprintf("listen_port=%s", string(C.GoString(listenPort))),
		fmt.Sprintf("public_key=%s", hex.EncodeToString(ppk),
		"replace_allowed_ips=true",
		fmt.Sprintf("allowed_ip=%s", string(C.GoString(allowedIps))),
		fmt.Sprintf("endpoint=%s", string(C.GoString(endpoint))),
	}
	err = device.IpcSetOperationByString(lines)
	if err != nil {
		return C.CString(err.Error())
	}
	err = device.Up()
	if err != nil {
		return C.CString(err.Error())
	}
	//logger.Verbosef("UAPI listener started")

	// wait for program to terminate
	signal.Notify(term, unix.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	case <-device.Wait():
	}

	// clean up

	//uapi.Close()
	device.Close()

	logger.Verbosef("Shutting down")
	return C.CString("success")
}

////export stopTun
//func stopTun() {
//	rTunLock.Lock()
//	defer rTunLock.Unlock()
//
//	if rTun != nil {
//		rTun.close()
//		rTun = nil
//	}
//}
