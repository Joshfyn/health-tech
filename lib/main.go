package lib

import (
	"fmt"
	"os"
	"strings"
)

func IsProduction() bool {
	return os.Getenv("ENV") == "prod"
}

func GetServiceHostWithPort(svcName string) string {
	if !IsProduction() {
		return fmt.Sprintf("localhost:%s", os.Getenv(fmt.Sprintf("%s_SERVICE_PORT", strings.ToUpper(svcName))))
	}
	return "0.0.0.0:80"
}
