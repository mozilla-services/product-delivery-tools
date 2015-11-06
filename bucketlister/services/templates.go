package services

import (
	"html/template"
	"path"
	"strings"
)

var s3Escaper = strings.NewReplacer("+", `%2B`, ":", "%3A")
var urlEscaper = strings.NewReplacer(":", "%3A")

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
	return s3Escaper.Replace(l.Path)
}

func (l *listTemplateInput) PathURLEscaped() string {
	return urlEscaper.Replace(l.Path)
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
				<td><a href="{{$.PathURLEscaped}}{{$dir.URLEscaped}}">{{.}}</a></td>
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
