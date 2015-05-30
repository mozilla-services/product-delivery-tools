package bucketlister

import "html/template"

type listFileInfo struct {
	Key          string
	LastModified string
	Size         string
}

type listTemplateInput struct {
	Path        string
	Directories []string
	Files       []listFileInfo
}

var listTemplate = template.Must(template.New("List").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Directory Listing: /{{.Path}}</title>
	</head>
	<body>
		<h1>Index of /{{.Path}}</h1>
		<table>
			<tr>
				<th>Type</th>
				<th>Name</th>
				<th>Size</th>
				<th>Last Modified</th>
			</tr>
			{{range .Directories}}
			<tr>
				<td>Dir</th>
				<td><a href="/{{.}}">/{{.}}</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range .Files}}
			<tr>
				<td>File</th>
				<td><a href="/{{.Key}}">/{{.Key}}</a></td>
				<td>{{.Size}}</td>
				<td>{{.LastModified}}</td>
			</tr>
			{{end}}
		</table>
	</body>
</html>
`))
