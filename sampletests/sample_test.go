package main

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

const factor = 20

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func testn(t *testing.T, n int) {
	duration := time.Duration(rand.Intn(factor)) * time.Second
	t.Logf("Test%d: Will sleep for %s", n, duration)
	time.Sleep(duration)
	t.Logf("Test%d: Some more output %s", n, duration)
	t.Logf("Environment var TF_ACC: %s", os.Getenv("TF_ACC"))
}

func Test1(t *testing.T) {
	testn(t, 1)
}

func Test2(t *testing.T) {
	testn(t, 2)
}

func Test3(t *testing.T) {
	testn(t, 3)
}

func Test4(t *testing.T) {
	testn(t, 4)
}

func TestFail(t *testing.T) {
	t.Fatalf("Test fails")
}

func TestSkip(t *testing.T) {
	t.Skip("Skipping this")
}

func TestPanic(t *testing.T) {
	t.Log("doing some work first...")
	panic("Whoops")
}
