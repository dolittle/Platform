package admin

import (
	"github.com/dolittle/platform-router/microservices"
	"html/template"
	"net/http"
)

var page = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>Current registry</title>
    </head>
    <body>
        <h1>Current contents in registry</h1>
        <table>
            <thead>
                <tr>
                    <th>Tenant:</th>
                    <th>Application:</th>
                    <th>Environment:</th>
                    <th>Microservice:</th>
                    <th>IP address:</th>
                </tr>
            </thead>
            <tbody>
            {{ range . }}
                <tr>
                    <td>{{ .Identity.Tenant }}</td>
                    <td>{{ .Identity.Application }}</td>
                    <td>{{ .Identity.Environment }}</td>
                    <td>{{ .Identity.Microservice }}</td>
                    <td>{{ .IP }}</td>
                </tr>
            {{ end }}
            </tbody>
        </table>
    </body>
</html>
`))

type RegistryHandler struct {
	Registry *microservices.Registry
}

func (rh RegistryHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	page.Execute(w, rh.Registry.All())
}
