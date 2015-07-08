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
		BucketMount{"pub/calendar/", "calendar"},
		BucketMount{"pub/firefox/bundles/", "hg-bundles"},
		BucketMount{"pub/firefox/try-builds/", "firefox-try"},
		BucketMount{"pub/firefox/", "firefox"},
		BucketMount{"pub/labs/", "labs"},
		BucketMount{"pub/mobile/try-builds/", "firefox-android-try"},
		BucketMount{"pub/mobile/", "firefox-android"},
		BucketMount{"pub/nspr/", "security"},
		BucketMount{"pub/opus/", "opus"},
		BucketMount{"pub/seamonkey/", "seamonkey"},
		BucketMount{"pub/security/", "security"},
		BucketMount{"pub/thunderbird/try-builds/", "thunderbird-try"},
		BucketMount{"pub/thunderbird/", "thunderbird"},
		BucketMount{"pub/xulrunner/try-builds/", "xulrunner-try"},
		BucketMount{"pub/xulrunner/", "xulrunner"},
		BucketMount{"pub/webtools/", "webtools"},
	},
}
