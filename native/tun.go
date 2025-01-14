package main

//#include "bridge.h"
import "C"
import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"time"
	"unsafe"
	"wg/native/conn"
	"wg/native/device"
	"wg/native/erro"
	"wg/native/tun"
)

//export startTun
func startTun(fd C.int, devicePrivateKey, listenPort, peerPublicKey, allowedIps, endpoint C.c_string, callback unsafe.Pointer, socketFd *int) C.c_string {
	var foreground bool
	var interfaceName string
	erro.Errinit(int(fd))
	foreground = false

	//err := erro.TestFunc()
	//if err != nil {
	//	return C.CString(err.Error())
	//}
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
		return tun.WgCreateTun(int(fd), device.DefaultMTU)
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
	//err = erro.TestFunc()
	//if err != nil {
	//	return C.CString(err.Error())
	//}
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
	markSocketfunc := func(fd int) {
		C.mark_socket(callback, C.int(fd))
		//time.Sleep(100 * time.Millisecond)
	}
	device := device.NewDevice(tdev, conn.NewDefaultBind(markSocketfunc), logger)
	//err = erro.TestFunc()
	//if err != nil {
	//	return C.CString(err.Error())
	//}
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
	//	go device.IpcHandle(conn)
	//	}
	//}()
	dpk, err := base64.StdEncoding.DecodeString(C.GoString(devicePrivateKey))
	if err != nil {
		return C.CString(err.Error())
	}
	ppk, err := base64.StdEncoding.DecodeString(C.GoString(peerPublicKey))
	if err != nil {
		return C.CString(err.Error())
	}
	lines := []string{
		fmt.Sprintf("private_key=%s", hex.EncodeToString(dpk)),
		fmt.Sprintf("listen_port=%s", string(C.GoString(listenPort))),
		fmt.Sprintf("public_key=%s", hex.EncodeToString(ppk)),
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
	time.Sleep(time.Second)
	*socketFd = erro.Fd
	//logger.Verbosef("UAPI listener started")
	//err = erro.TestFunc()
	//if err != nil {
	//	return C.CString(err.Error())
	//}
	// wait for program to terminate
	signal.Notify(term, unix.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
		return C.CString("term")
	case err = <-errs:
		return C.CString("error:" + err.Error())
	case <-device.Wait():
		return C.CString("closed")
	case err = <-erro.Err:
		return C.CString("error:" + err.Error())
	}

	// clean up

	//uapi.Close()
	device.Close()

	logger.Verbosef("Shutting down")
	return C.CString("success")
}

func markSocket(fd int, callback unsafe.Pointer) {
	C.mark_socket(callback, C.int(fd))
}
