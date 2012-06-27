package serial

import (
	"os"
	"syscall"
)

type Serial struct {
	f *os.File
}

// Defaults: 9600 8N1 soft/hard flow controll off
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

func (s *Serial) Speed(boud int) error {
	return s.speed(boud)
}
