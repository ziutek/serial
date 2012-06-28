package serial

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

// bits/termios.h
const nccs = 32

type cc_t byte
type speed_t uint
type tcflag_t uint
type termios struct {
	c_iflag  tcflag_t   // input specific flags (bitmask)
	c_oflag  tcflag_t   // output specific flags (bitmask)
	c_cflag  tcflag_t   // control flags (bitmask)
	c_lflag  tcflag_t   // local flags (bitmask)
	c_cc     [nccs]cc_t // special characters
	c_ispeed speed_t    // input speed 
	c_ospeed speed_t    // output speed 
}

// bits/termios.h
const (
	ignbrk = 1 << iota
	brkint
	ignpar
	parmrk
	inpck
	istrip
	inlcr
	igncr
	icrnl
	iuclc
	ixon
	ixany
	ixoff
	imaxbeL
	iutf8
)

// bits/termios.h
const (
	opost = 1 << iota
	olcuc
	onlcr
	ocrnl
	onocr
	onlret
	ofill
	ofdel
	nldly
)

// bits/termios.h
const (
	isig = 1 << iota
	icanon
	xcase
	echo
	echoe
	echok
	echonl
	noflsh
	tostop
	echoctl
	echoprt
	echoke
	flusho
	_undef_
	pendin
	iexten
	extproc
)

// bits/termios.h
var boud = map[int]tcflag_t{
	0:       0000000,
	50:      0000001,
	75:      0000002,
	110:     0000003,
	134:     0000004,
	150:     0000005,
	200:     0000006,
	300:     0000007,
	600:     0000010,
	1200:    0000011,
	1800:    0000012,
	2400:    0000013,
	4800:    0000014,
	9600:    0000015,
	19200:   0000016,
	38400:   0000017,
	57600:   0010001,
	115200:  0010002,
	230400:  0010003,
	460800:  0010004,
	500000:  0010005,
	576000:  0010006,
	921600:  0010007,
	1000000: 0010010,
	1152000: 0010011,
	1500000: 0010012,
	2000000: 0010013,
	2500000: 0010014,
	3000000: 0010015,
	3500000: 0010016,
	4000000: 0010017,
}

// bits/termios.h
var bits = map[int]tcflag_t{
	5: 0000000,
	6: 0000020,
	7: 0000040,
	8: 0000060,
}

// bits/termios.h
const (
	cstopb = 0000100 << iota
	cread
	parenb
	parodd
	hupcl
	clocal
)

// bits/termios.h
const (
	vintr = iota
	vquit
	verase
	vkill
	veof
	vtime
	vmin
	vswtc
	vstart
	vstop
	vsusp
	veol
	vreprint
	vdiscard
	vwerase
	vlnext
	veol2
)

// bits/termios.h
const (
	cbaud   = 0010017
	cbaudex = 0010000
	crtscts = 020000000000
)

// asm-generic/ioctls.h
const (
	tcgets = 0x5401
	tcsets = 0x5402
)

func (s *Serial) tcGetAttr(cfg *termios) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		tcgets,
		uintptr(unsafe.Pointer(cfg)),
	)
	if e != 0 {
		return os.NewSyscallError("tcgetattr", e)
	}
	return nil
}

func (s *Serial) tcSetAttr(cfg *termios) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		tcsets,
		uintptr(unsafe.Pointer(cfg)),
	)
	if e != 0 {
		return os.NewSyscallError("tcsetattr", e)
	}
	return nil
}

func (s *Serial) init() error {
	var t termios
	t.c_iflag = 0
	t.c_oflag = 0
	t.c_lflag = 0
	t.c_cflag = boud[9600] | bits[8] | clocal | cread | hupcl
	t.c_cc[vmin] = 1
	t.c_cc[vtime] = 0
	t.c_ispeed = speed_t(boud[9600])
	t.c_ospeed = speed_t(boud[9600])
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	// Clear non-blocking flag (we need nonblocking only for Open)
	return syscall.SetNonblock(int(s.f.Fd()), false)
}

func (s *Serial) setSpeed(b int) error {
	var t termios
	bb, ok := boud[b]
	if !ok {
		return errors.New("Unknown boud rate")
	}
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	t.c_cflag &^= cbaud | cbaudex
	t.c_cflag |= bb
	t.c_ispeed = speed_t(b)
	t.c_ospeed = speed_t(b)
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setParity(parity, odd bool) error {
	var t termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if parity {
		t.c_cflag |= parenb
	} else {
		t.c_cflag &^= parenb
	}
	if odd {
		t.c_cflag |= parodd
	} else {
		t.c_cflag &^= parodd
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setStopBits2(two bool) error {
	var t termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if two {
		t.c_cflag |= cstopb
	} else {
		t.c_cflag &^= cstopb
	}

	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setFlowCtrl(hw, soft bool) error {
	var t termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if hw {
		t.c_cflag |= crtscts
	} else {
		t.c_cflag &^= crtscts
	}
	if soft {
		t.c_iflag |= (ixon | ixoff | ixany)
	} else {
		t.c_iflag &^= (ixon | ixoff | ixany)
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setMode(canon bool) error {
	var t termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if canon {
		t.c_iflag |= icanon
	} else {
		t.c_iflag &^= icanon
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setRawRead(vmin, vtime int) error {
	var t termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	t.c_cc[vmin] = cc_t(vmin)
	t.c_cc[vtime] = cc_t(vtime)
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil

}
