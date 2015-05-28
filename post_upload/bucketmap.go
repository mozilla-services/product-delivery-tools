package main

import "strings"

var pathPrefixesToBucket = [][]string{
	[]string{"firefox/try-builds/", "firefox-try"},
	[]string{"firefox/", "firefox"},
	[]string{"mobile/try-builds/", "firefox-android-try"},
	[]string{"mobile/", "firefox-android"},
	[]string{"opus/", "opus"},
	[]string{"thunderbird/try-builds/", "thunderbird-try"},
	[]string{"thunderbird/", "thunderbird"},
	[]string{"xulrunner/try-builds/", "xulrunner-try"},
	[]string{"xulrunner/", "xulrunner"},
}

func destToBucket(dest string) string {
	for _, pathPrefix := range pathPrefixesToBucket {
		if strings.HasPrefix(dest, pathPrefix[0]) {
			return pathPrefix[1]
		}
	}

	return "ftp-archive"
}
