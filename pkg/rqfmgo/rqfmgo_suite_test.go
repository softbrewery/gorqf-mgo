package rqfmgo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRqfmgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rqfmgo Suite")
}
