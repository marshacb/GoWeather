package bookmark_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBookmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bookmark Suite")
}
