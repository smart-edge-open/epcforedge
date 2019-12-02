package ngcnef_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNef(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nef Suite")
}
