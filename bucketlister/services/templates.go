package services

import (
	"html/template"
	"path"
	"strings"
)

type listFileInfo struct {
	Key          string
	LastModified string
	Size         string
}

type listTemplateInput struct {
	Path          string
	PrefixListing *PrefixListing
}

func (l *listTemplateInput) Parent() string {
	if l.Path == "/" {
		return "/"
	}

	d := path.Dir(strings.TrimSuffix(l.Path, "/"))
	if !strings.HasSuffix(d, "/") {
		d += "/"
	}
	return d
}

var listTemplate = template.Must(template.New("List").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Directory Listing: {{.Path}}</title>
	</head>
	<body>
		<h1>Index of {{.Path}}</h1>
		<table>
			<tr>
				<th>Type</th>
				<th>Name</th>
				<th>Size</th>
				<th>Last Modified</th>
			</tr>
			{{if ne $.Path "/"}}
			<tr>
				<td>Dir</td>
				<td><a href="{{$.Parent}}">..</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range $dir := .PrefixListing.Prefixes}}
			<tr>
				<td>Dir</td>
				<td><a href="{{$.Path}}{{$dir}}">{{.}}</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range $file := .PrefixListing.Files}}
			<tr>
				<td>File</th>
				<td><a href="{{$.Path}}{{$file.Base}}">{{$file.Base}}</a></td>
				<td>{{$file.Size}}</td>
				<td>{{$file.LastModified}}</td>
			</tr>
			{{end}}
		</table>
	</body>
</html>
`))
