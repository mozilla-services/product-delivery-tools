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

func (l *listTemplateInput) PathEscaped() string {
	parts := strings.Split(l.Path, "/")
	escapedParts := make([]string, len(parts))
	for i, p := range parts {
		escapedParts[i] = template.URLQueryEscaper(p)
	}
	return strings.Join(escapedParts, "/")
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
			{{range $dir := .PrefixListing.PrefixStructs}}
			<tr>
				<td>Dir</td>
				<td><a href="{{$.PathEscaped}}{{$dir.Escaped}}">{{.}}</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range $file := .PrefixListing.Files}}
			{{if ne $file.Base "."}}
			<tr>
				<td>File</td>
				<td><a href="{{$.PathEscaped}}{{$file.BaseEscaped}}">{{$file.Base}}</a></td>
				<td>{{$file.SizeString}}</td>
				<td>{{$file.LastModifiedString}}</td>
			</tr>
			{{end}}
			{{end}}
		</table>
	</body>
</html>
`))
