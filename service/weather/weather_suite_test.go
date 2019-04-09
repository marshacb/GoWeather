package weather_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWeather(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Weather Suite")
}
