package aws

import "strings"

func SkipS3Upload(name string) bool {
	if strings.HasSuffix(name, ".log") {
		return true
	} else if strings.HasSuffix(name, ".dbtmp") {
		return true
	}
	return false
}
