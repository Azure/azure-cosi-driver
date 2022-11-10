package driver

import (
	"sync"
	"testing"

	"google.golang.org/grpc"
)

func TestWait(t *testing.T) {
	s := COSIServer{}
	s.server = grpc.NewServer()
	s.waitGroup = sync.WaitGroup{}
	s.Wait()
}

func TestGracefulStop(t *testing.T) {
	s := COSIServer{}
	s.server = grpc.NewServer()
	s.GracefulStop()
}

func TestStop(t *testing.T) {
	s := COSIServer{}
	s.server = grpc.NewServer()
	s.Stop()
}
