package sloppy_duration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSloppyDuration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SloppyDuration Suite")
}
