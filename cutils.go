package main

/*
#include <stdlib.h>
#include <sys/stat.h>
#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

func DropPrivileges(name string) {
	cname := C.CString(name)
	home := C.CString("HOME")
	slash := C.CString("/")
	defer C.cfree(unsafe.Pointer(home))
	defer C.cfree(unsafe.Pointer(cname))
	defer C.cfree(unsafe.Pointer(slash))
	cpw := C.getpwnam(cname)
	C.setgid(cpw.pw_gid)
	C.setuid(cpw.pw_uid)
	C.setenv(home, cpw.pw_dir, 1)
	C.setsid()
	C.chdir(slash)
	C.umask(022)
}
