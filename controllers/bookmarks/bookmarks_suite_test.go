package bookmarks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBookmarks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bookmarks Suite")
}
