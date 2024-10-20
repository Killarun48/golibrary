package main

import (
	"os"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	pathDB = "./users_test.db"
	s := NewServer(":8088")

	go s.Serve()

	s.Stop()
	time.Sleep(1 * time.Second)

	os.Remove("./users_test.db")
}
