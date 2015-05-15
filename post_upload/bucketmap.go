package main

import "strings"

var pathPrefixesToBucket = [][]string{
	[]string{"firefox/try-builds/", "delivery-firefox-try"},
	[]string{"firefox/", "delivery-firefox"},
	[]string{"mobile/try-builds/", "delivery-firefox-android-try"},
	[]string{"mobile/", "delivery-firefox-android"},
	[]string{"opus/", "delivery-opus"},
	[]string{"thunderbird/try-builds/", "delivery-thunderbird-try"},
	[]string{"thunderbird/", "delivery-thunderbird"},
	[]string{"xulrunner/try-builds/", "delivery-xulrunner-try"},
	[]string{"xulrunner/", "delivery-xulrunner"},
}

func destToBucket(dest string) string {
	for _, pathPrefix := range pathPrefixesToBucket {
		if strings.HasPrefix(dest, pathPrefix[0]) {
			return pathPrefix[1]
		}
	}

	return "ftp-archive"
}
