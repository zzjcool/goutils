package castwait

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCWBefore 完成条件在wait之前
func TestCWBefore(t *testing.T) {
	n := 100
	c := New()
	exErr := fmt.Errorf("testcond")
	c.Done(exErr)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i += 1 {
		go func() {
			err := c.Wait()
			assert.Equal(t, exErr, err)
			wg.Done()
		}()
	}

	wg.Wait()
	time.Sleep(time.Second)
}

// TestConCWAfter 完成条件在wait之后
func TestConCWAfter(t *testing.T) {
	n := 100000
	c := New()
	exErr := fmt.Errorf("testcond")

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i += 1 {
		go func() {
			err := c.Wait()
			assert.Equal(t, exErr, err)
			wg.Done()
		}()
	}

	c.Done(exErr)
	wg.Wait()
}

// TestCondBefore 完成条件在wait之前
func TestCondBefore(t *testing.T) {
	n := 100
	c := NewCond()
	exErr := fmt.Errorf("testcond")
	c.Done(exErr)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i += 1 {
		go func() {
			err := c.Wait()
			assert.Equal(t, exErr, err)
			wg.Done()
		}()
	}

	wg.Wait()
	time.Sleep(time.Second)
}

// TestCondAfter 完成条件在wait之后
func TestCondAfter(t *testing.T) {
	n := 100000
	c := NewCond()
	exErr := fmt.Errorf("testcond")

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i += 1 {
		go func() {
			err := c.Wait()
			assert.Equal(t, exErr, err)
			wg.Done()
		}()
	}

	c.Done(exErr)
	wg.Wait()
}
