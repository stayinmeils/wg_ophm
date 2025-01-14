//go:build !windows

/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2023 WireGuard LLC. All Rights Reserved.
 */

package main

import "C"
import (
	"fmt"
	"os"
	"runtime"
)

const (
	ExitSetupSuccess = 0
	ExitSetupFailed  = 1
)

const (
	ENV_WG_TUN_FD             = "WG_TUN_FD"
	ENV_WG_UAPI_FD            = "WG_UAPI_FD"
	ENV_WG_PROCESS_FOREGROUND = "WG_PROCESS_FOREGROUND"
)

func printUsage() {
	fmt.Printf("Usage: %s [-f/--foreground] INTERFACE-NAME\n", os.Args[0])
}

func warning() {
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd":
		if os.Getenv(ENV_WG_PROCESS_FOREGROUND) == "1" {
			return
		}
	default:
		return
	}

	fmt.Fprintln(os.Stderr, "┌──────────────────────────────────────────────────────┐")
	fmt.Fprintln(os.Stderr, "│                                                      │")
	fmt.Fprintln(os.Stderr, "│   Running wireguard-go is not required because this  │")
	fmt.Fprintln(os.Stderr, "│   kernel has first class support for WireGuard. For  │")
	fmt.Fprintln(os.Stderr, "│   information on installing the kernel module,       │")
	fmt.Fprintln(os.Stderr, "│   please visit:                                      │")
	fmt.Fprintln(os.Stderr, "│         https://www.wireguard.com/install/           │")
	fmt.Fprintln(os.Stderr, "│                                                      │")
	fmt.Fprintln(os.Stderr, "└──────────────────────────────────────────────────────┘")
}

func main() {
	os.Exit(0)
	//var foreground bool
	//var interfaceName string
	//foreground = false
	//
	////err := erro.TestFunc()
	////if err != nil {
	////	return C.CString(err.Error())
	////}
	//if !foreground {
	//	foreground = os.Getenv(ENV_WG_PROCESS_FOREGROUND) == "1"
	//}
	//
	//// get log level (default: info)
	//
	//logLevel := func() int {
	//	switch os.Getenv("LOG_LEVEL") {
	//	case "verbose", "debug":
	//		return device.LogLevelVerbose
	//	case "error":
	//		return device.LogLevelError
	//	case "silent":
	//		return device.LogLevelSilent
	//	}
	//	return device.LogLevelError
	//}()
	//
	//// open TUN device (or use supplied fd)
	//
	////tdev, err := func() (tun.Device, error) {
	////	//tunFdStr := os.Getenv(ENV_WG_TUN_FD)
	////	//if tunFdStr == "" {
	////	return tun.LinuxCreateTun(device.DefaultMTU)
	////	//}
	////
	////	// construct tun device from supplied fd
	////
	////	//fd, err := strconv.ParseUint(tunFdStr, 10, 32)
	////	//if err != nil {
	////	//	return nil, err
	////	//}
	////	//
	////	//err = unix.SetNonblock(int(fd), true)
	////	//if err != nil {
	////	//	return nil, err
	////	//}
	////	//
	////	//file := os.NewFile(uintptr(fd), "")
	////	//return tun.CreateTUNFromFile(file, device.DefaultMTU)
	////}()
	//
	//if err == nil {
	//	realInterfaceName, err2 := tdev.Name()
	//	if err2 == nil {
	//		interfaceName = realInterfaceName
	//	}
	//}
	////err = erro.TestFunc()
	////if err != nil {
	////	return C.CString(err.Error())
	////}
	//logger := device.NewLogger(
	//	logLevel,
	//	fmt.Sprintf("(%s) ", interfaceName),
	//)
	//
	//logger.Verbosef("Starting wireguard-go version %s", Version)
	//
	//if err != nil {
	//	logger.Errorf("Failed to create TUN device: %v", err)
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}
	//
	//// open UAPI file (or use supplied fd)
	//
	////fileUAPI, err := func() (*os.File, error) {
	////	uapiFdStr := os.Getenv(ENV_WG_UAPI_FD)
	////	if uapiFdStr == "" {
	////		return ipc.UAPIOpen(interfaceName)
	////	}
	////
	////	// use supplied fd
	////
	////	fd, err := strconv.ParseUint(uapiFdStr, 10, 32)
	////	if err != nil {
	////		return nil, err
	////	}
	////
	////	return os.NewFile(uintptr(fd), ""), nil
	////}()
	////if err != nil {
	////	logger.Errorf("UAPI listen error: %v", err)
	////	os.Exit(ExitSetupFailed)
	////	return -1
	////}
	//// daemonize the process
	//
	////if !foreground {
	////	env := os.Environ()
	////	env = append(env, fmt.Sprintf("%s=3", ENV_WG_TUN_FD))
	////	env = append(env, fmt.Sprintf("%s=4", ENV_WG_UAPI_FD))
	////	env = append(env, fmt.Sprintf("%s=1", ENV_WG_PROCESS_FOREGROUND))
	////	files := [3]*os.File{}
	////	if os.Getenv("LOG_LEVEL") != "" && logLevel != device.LogLevelSilent {
	////		files[0], _ = os.Open(os.DevNull)
	////		files[1] = os.Stdout
	////		files[2] = os.Stderr
	////	} else {
	////		files[0], _ = os.Open(os.DevNull)
	////		files[1], _ = os.Open(os.DevNull)
	////		files[2], _ = os.Open(os.DevNull)
	////	}
	////	attr := &os.ProcAttr{
	////		Files: []*os.File{
	////			files[0], // stdin
	////			files[1], // stdout
	////			files[2], // stderr
	////			tdev.File(),
	////			fileUAPI,
	////		},
	////		Dir: ".",
	////		Env: env,
	////	}
	////
	////	path, err := os.Executable()
	////	if err != nil {
	////		logger.Errorf("Failed to determine executable: %v", err)
	////		os.Exit(ExitSetupFailed)
	////	}
	////
	////	process, err := os.StartProcess(
	////		path,
	////		os.Args,
	////		attr,
	////	)
	////	if err != nil {
	////		logger.Errorf("Failed to daemonize: %v", err)
	////		os.Exit(ExitSetupFailed)
	////	}
	////	process.Release()
	////	return -1
	////}
	//
	//device := device.NewDevice(tdev, conn.NewDefaultBind(), logger)
	////err = erro.TestFunc()
	////if err != nil {
	////	return C.CString(err.Error())
	////}
	//
	//logger.Verbosef("Device started")
	//
	//errs := make(chan error)
	//term := make(chan os.Signal, 1)
	//
	////uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
	////if err != nil {
	////	logger.Errorf("Failed to listen on uapi socket: %v", err)
	////	os.Exit(ExitSetupFailed)
	////}
	//
	////go func() {
	////	for {
	////		conn, err := uapi.Accept()
	////		if err != nil {
	////			errs <- err
	////			return
	////		}
	////	go device.IpcHandle(conn)
	////	}
	////}()
	//ppk := []byte{220, 239, 61, 112, 192, 52, 204, 107, 74, 161, 131, 117, 5, 226, 246, 65, 150, 94, 134, 223, 226, 251, 200, 34, 243, 207, 134, 192, 251, 89, 143, 75}
	//dpk := []byte{40, 136, 93, 85, 254, 103, 10, 46, 94, 2, 115, 66, 128, 37, 161, 21, 65, 209, 215, 198, 87, 250, 94, 179, 150, 134, 155, 249, 127, 65, 154, 66}
	//lines := []string{
	//	fmt.Sprintf("private_key=%s", hex.EncodeToString(dpk)),
	//	fmt.Sprintf("listen_port=%s", "51820"),
	//	fmt.Sprintf("public_key=%s", hex.EncodeToString(ppk)),
	//	"replace_allowed_ips=true",
	//	fmt.Sprintf("allowed_ip=%s", "0.0.0.0/0"),
	//	fmt.Sprintf("endpoint=%s", "8.209.253.159:51829"),
	//}
	//err = device.IpcSetOperationByString(lines)
	//if err != nil {
	//	fmt.Println(C.CString(err.Error()))
	//}
	//err = device.Up()
	//if err != nil {
	//	fmt.Println(C.CString(err.Error()))
	//}
	////logger.Verbosef("UAPI listener started")
	////err = erro.TestFunc()
	////if err != nil {
	////	return C.CString(err.Error())
	////}
	//// wait for program to terminate
	//signal.Notify(term, unix.SIGTERM)
	//signal.Notify(term, os.Interrupt)
	//
	//select {
	//case <-term:
	//	fmt.Println("term")
	//case err = <-errs:
	//	fmt.Println("error:" + err.Error())
	//case <-device.Wait():
	//	fmt.Println("closed")
	//case err = <-erro.Err:
	//	fmt.Println("error:" + err.Error())
	//}
	//
	//// clean up
	//
	////uapi.Close()
	//device.Close()
	//
	//logger.Verbosef("Shutting down")
}

//func main() {
//	if len(os.Args) == 2 && os.Args[1] == "--version" {
//		fmt.Printf("wireguard-go v%s\n\nUserspace WireGuard daemon for %s-%s.\nInformation available at https://www.wireguard.com.\nCopyright (C) Jason A. Donenfeld <Jason@zx2c4.com>.\n", Version, runtime.GOOS, runtime.GOARCH)
//		return
//	}
//
//	warning()
//
//	var foreground bool
//	var interfaceName string
//	if len(os.Args) < 2 || len(os.Args) > 3 {
//		printUsage()
//		return
//	}
//
//	switch os.Args[1] {
//
//	case "-f", "--foreground":
//		foreground = true
//		if len(os.Args) != 3 {
//			printUsage()
//			return
//		}
//		interfaceName = os.Args[2]
//
//	default:
//		foreground = false
//		if len(os.Args) != 2 {
//			printUsage()
//			return
//		}
//		interfaceName = os.Args[1]
//	}
//
//	if !foreground {
//		foreground = os.Getenv(ENV_WG_PROCESS_FOREGROUND) == "1"
//	}
//
//	// get log level (default: info)
//
//	logLevel := func() int {
//		switch os.Getenv("LOG_LEVEL") {
//		case "verbose", "debug":
//			return device.LogLevelVerbose
//		case "error":
//			return device.LogLevelError
//		case "silent":
//			return device.LogLevelSilent
//		}
//		return device.LogLevelError
//	}()
//
//	// open TUN device (or use supplied fd)
//
//	tdev, err := func() (tun.Device, error) {
//		tunFdStr := os.Getenv(ENV_WG_TUN_FD)
//		if tunFdStr == "" {
//			return tun.CreateTUN(interfaceName, device.DefaultMTU)
//		}
//
//		// construct tun device from supplied fd
//
//		fd, err := strconv.ParseUint(tunFdStr, 10, 32)
//		if err != nil {
//			return nil, err
//		}
//
//		err = unix.SetNonblock(int(fd), true)
//		if err != nil {
//			return nil, err
//		}
//
//		file := os.NewFile(uintptr(fd), "")
//		return tun.CreateTUNFromFile(file, device.DefaultMTU)
//	}()
//
//	if err == nil {
//		realInterfaceName, err2 := tdev.Name()
//		if err2 == nil {
//			interfaceName = realInterfaceName
//		}
//	}
//
//	logger := device.NewLogger(
//		logLevel,
//		fmt.Sprintf("(%s) ", interfaceName),
//	)
//
//	logger.Verbosef("Starting wireguard-go version %s", Version)
//
//	if err != nil {
//		logger.Errorf("Failed to create TUN device: %v", err)
//		os.Exit(ExitSetupFailed)
//	}
//
//	// open UAPI file (or use supplied fd)
//
//	fileUAPI, err := func() (*os.File, error) {
//		uapiFdStr := os.Getenv(ENV_WG_UAPI_FD)
//		if uapiFdStr == "" {
//			return ipc.UAPIOpen(interfaceName)
//		}
//
//		// use supplied fd
//
//		fd, err := strconv.ParseUint(uapiFdStr, 10, 32)
//		if err != nil {
//			return nil, err
//		}
//
//		return os.NewFile(uintptr(fd), ""), nil
//	}()
//	if err != nil {
//		logger.Errorf("UAPI listen error: %v", err)
//		os.Exit(ExitSetupFailed)
//		return
//	}
//	// daemonize the process
//
//	if !foreground {
//		env := os.Environ()
//		env = append(env, fmt.Sprintf("%s=3", ENV_WG_TUN_FD))
//		env = append(env, fmt.Sprintf("%s=4", ENV_WG_UAPI_FD))
//		env = append(env, fmt.Sprintf("%s=1", ENV_WG_PROCESS_FOREGROUND))
//		files := [3]*os.File{}
//		if os.Getenv("LOG_LEVEL") != "" && logLevel != device.LogLevelSilent {
//			files[0], _ = os.Open(os.DevNull)
//			files[1] = os.Stdout
//			files[2] = os.Stderr
//		} else {
//			files[0], _ = os.Open(os.DevNull)
//			files[1], _ = os.Open(os.DevNull)
//			files[2], _ = os.Open(os.DevNull)
//		}
//		attr := &os.ProcAttr{
//			Files: []*os.File{
//				files[0], // stdin
//				files[1], // stdout
//				files[2], // stderr
//				tdev.File(),
//				fileUAPI,
//			},
//			Dir: ".",
//			Env: env,
//		}
//
//		path, err := os.Executable()
//		if err != nil {
//			logger.Errorf("Failed to determine executable: %v", err)
//			os.Exit(ExitSetupFailed)
//		}
//
//		process, err := os.StartProcess(
//			path,
//			os.Args,
//			attr,
//		)
//		if err != nil {
//			logger.Errorf("Failed to daemonize: %v", err)
//			os.Exit(ExitSetupFailed)
//		}
//		process.Release()
//		return
//	}
//
//	device := device.NewDevice(tdev, conn.NewDefaultBind(), logger)
//
//	logger.Verbosef("Device started")
//
//	errs := make(chan error)
//	term := make(chan os.Signal, 1)
//
//	uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
//	if err != nil {
//		logger.Errorf("Failed to listen on uapi socket: %v", err)
//		os.Exit(ExitSetupFailed)
//	}
//
//	go func() {
//		for {
//			conn, err := uapi.Accept()
//			if err != nil {
//				errs <- err
//				return
//			}
//			go device.IpcHandle(conn)
//		}
//	}()
//
//	logger.Verbosef("UAPI listener started")
//
//	// wait for program to terminate
//
//	signal.Notify(term, unix.SIGTERM)
//	signal.Notify(term, os.Interrupt)
//
//	select {
//	case <-term:
//	case <-errs:
//	case <-device.Wait():
//	}
//
//	// clean up
//
//	uapi.Close()
//	device.Close()
//
//	logger.Verbosef("Shutting down")
//}
