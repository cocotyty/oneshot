package oneshot

import (
	"sync/atomic"
	"unsafe"
	_ "unsafe"
)

type Shot struct {
	gp     uintptr
}

func getG() uintptr

//go:linkname gopark runtime.gopark
func gopark(unlockf func(uintptr, unsafe.Pointer) bool, lock unsafe.Pointer, reason byte, traceEv byte, traceskip int)

//go:linkname goready runtime.goready
func goready(gp uintptr, traceskip int)

func commit(ug uintptr, xp unsafe.Pointer) bool {
	x := (*Shot)(xp)
	return atomic.CompareAndSwapUintptr(&x.gp, 0, ug)
}

func (x *Shot) Wait() {
	gopark(commit, unsafe.Pointer(x), 1, 0, 0)
}

func (x *Shot) WakeUp() {
	ugn := atomic.LoadUintptr(&x.gp)
	if ugn != 0 {
		goready(ugn, 0)
	} else {
		if !atomic.CompareAndSwapUintptr(&x.gp, 0, 1) {
			ugn := atomic.LoadUintptr(&x.gp)
			goready(ugn, 0)
		}
	}
}
