package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oam Suite")
}
