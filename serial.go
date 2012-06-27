package serial

import (
	"os"
	"syscall"
)

type Serial struct {
	f *os.File
}

// Defaults:
//  9600 8N1, soft/hard flow controll off, raw mode, vmin=1, vtime=0
func Open(name string) (*Serial, error) {
	f, err := os.OpenFile(
		name,
		os.O_RDWR|syscall.O_NONBLOCK|syscall.O_NOCTTY,
		0600,
	)
	if err != nil {
		return nil, err
	}
	s := &Serial{f}
	err = s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Serial) Close() error {
	err := s.f.Close()
	s.f = nil
	return err
}

func (s *Serial) Read(b []byte) (int, error) {
	return s.f.Read(b)
}

func (s *Serial) WriteString(str string) (int, error) {
	return s.f.WriteString(str)
}

func (s *Serial) Write(b []byte) (int, error) {
	return s.f.Write(b)
}

func (s *Serial) WriteByte(c byte) error {
	_, e := s.f.Write([]byte{c})
	return e
}

func (s *Serial) ReadByte() (byte, error) {
	buf := []byte{0}
	n, e := s.f.Read(buf)
	if n == 1 {
		return buf[0], nil
	}
	return 0, e
}
func (s *Serial) Name() string {
	return s.f.Name()
}

func (s *Serial) File() *os.File {
	return s.f
}

func (s *Serial) SetSpeed(boud int) error {
	return s.setSpeed(boud)
}

func (s *Serial) SetFlowCtrl(hw, soft bool) error {
	return s.setFlowCtrl(hw, soft)
}

// Sets canonical/raw mode
func (s *Serial) SetMode(canon bool) error {
	return s.setMode(canon)
}

// Sets Read behavior for noncanonical mode.
//  vmin  - minimum number of characters for Read,
//  vtime - timeout in deciseconds,
//  vmin == 0 && vtime == 0 : non-blocking Read,
//  vmin == 0 && vtime > 0  : Read returns buffered charcters or waits
//                            vtime for new charcters,
//  vmin > 0  && vtime > 0  : Read returns n >= vmin charcters or
//                            0 < n < vmin if vtime expires after n-th char
//  vmin > 0  && vtime == 0 : Read returns at least vmin characters
func (s *Serial) SetRawRead(vmin, vtime int) error {
	return s.setRawRead(vmin, vtime)
}
