package bucketlister

import (
	"html/template"

	"github.com/awslabs/aws-sdk-go/service/s3"
)

type listTemplateInput struct {
	Prefixes []*s3.CommonPrefix
	Objects  []*s3.Object
}

var listTemplate = template.Must(template.New("List").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Directory Listing</title>
	</head>
	<body>
		<table>
			<tr>
				<th>Name</th>
				<th>Size</th>
				<th>Last Modified</th>
			</tr>
			{{range .Prefixes}}
			<tr>
				<td><a href="/{{.Prefix}}">/{{.Prefix}}</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range .Objects}}
			<tr>
				<td><a href="/{{.Key}}">/{{.Key}}</a></td>
				<td>{{.Size}}</td>
				<td>{{.LastModified}}</td>
			</tr>
			{{end}}
		</table>
	</body>
</html>
`))
