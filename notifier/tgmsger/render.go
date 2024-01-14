package notifier

import (
	"bytes"
	"cf_ddns/model"
	"html/template"
	"time"
)

const htmlTplt = `
<b>{{ .Title }}</b>
Domain: {{ .Domain }}
New IP: {{ .NewIP }}
Old IP: {{ .OldIP }}
Refresher: {{ .Refresher }}
`

func tsPrint(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func renderMsg(executor *template.Template, nt *model.Notification) (string, error) {
	buf := bytes.Buffer{}
	if err := executor.Execute(&buf, nt); err != nil {
		return "", err
	}
	return buf.String(), nil
}
