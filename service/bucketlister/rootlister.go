package bucketlister

import (
	"log"
	"net/http"
)

// RootLister lists bucketlisters are directories
type RootLister struct {
	listers []*BucketLister
}

// AddBucketLister adds a lister to the root lister
func (r *RootLister) AddBucketLister(b *BucketLister) {
	if r.listers == nil {
		r.listers = []*BucketLister{b}
		return
	}
	r.listers = append(r.listers, b)
}

// ServeHTTP implements http.Handler
func (r *RootLister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tmplParams := &listTemplateInput{
		Path:        "/",
		Directories: []string{},
	}
	for _, lister := range r.listers {
		if empty, err := lister.Empty(); empty {
			if err != nil {
				log.Printf("Error checking empty: %s", err)
			}
			continue
		}
		tmplParams.Directories = append(tmplParams.Directories, lister.mountedAt)
	}

	w.Header().Set("Content-Type", "text/html")
	err := listTemplate.Execute(w, tmplParams)
	if err != nil {
		log.Printf("Error executing template err: %s", err)
	}

}
