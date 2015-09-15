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
		BucketMount{"pub/firefox/bundles/", "archive"},
		BucketMount{"pub/firefox/try-builds/", "archive"},
		BucketMount{"pub/firefox/", "firefox"},
		BucketMount{"pub/labs/", "contrib"},
		BucketMount{"pub/opus/", "contrib"},
		BucketMount{"pub/webtools/", "contrib"},
	},
}
