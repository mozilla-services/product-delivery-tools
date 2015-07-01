package deliverytools

// BucketMap represents the mapping of prefixes to buckets
type BucketMap struct {
	// The default or fallthrough bucket
	Default string

	// Bucket mappings
	Mounts []BucketMount
}

// BucketMount is a prefix -> bucket mapping
type BucketMount struct {
	Prefix string
	Bucket string
}

// BucketMapping is the current BucketMap
var ProdBucketMap = BucketMap{
	Default: "archive",
	Mounts: []BucketMount{
		BucketMount{"pub/mozilla.org/calendar/", "calendar"},
		BucketMount{"pub/mozilla.org/firefox/bundles/", "hg-bundles"},
		BucketMount{"pub/mozilla.org/firefox/try-builds/", "firefox-try"},
		BucketMount{"pub/mozilla.org/firefox/", "firefox"},
		BucketMount{"pub/mozilla.org/labs/", "labs"},
		BucketMount{"pub/mozilla.org/mobile/try-builds/", "firefox-android-try"},
		BucketMount{"pub/mozilla.org/mobile/", "firefox-android"},
		BucketMount{"pub/mozilla.org/nspr/", "security"},
		BucketMount{"pub/mozilla.org/opus/", "opus"},
		BucketMount{"pub/mozilla.org/seamonkey/", "seamonkey"},
		BucketMount{"pub/mozilla.org/security/", "security"},
		BucketMount{"pub/mozilla.org/thunderbird/try-builds/", "thunderbird-try"},
		BucketMount{"pub/mozilla.org/thunderbird/", "thunderbird"},
		BucketMount{"pub/mozilla.org/xulrunner/try-builds/", "xulrunner-try"},
		BucketMount{"pub/mozilla.org/xulrunner/", "xulrunner"},
		BucketMount{"pub/mozilla.org/webtools/", "webtools"},
	},
}
