{{- if .Request.Number }}
<b>{{ .Request.Number }}</b>
{{- end }}
{{- if .Request.VIN }}
<b>{{ .Request.VIN }}</b>
{{- end }}

Всього: <b>{{ len .Vehicles }}</b> транспортних засобів
{{ range $vehicle := .Vehicles -}}
<b>{{ .Brand }} {{ .Model }} {{ .Year }}</b>
{{- if .VIN }}
<b>VIN: </b><a href="https://www.opencars.app/vin/{{ .VIN }}">{{ .VIN }}</a>
{{- end }}
{{- if not .FirstRegDate.IsZero }}
<b>Перша реєстрація: </b>{{ .FirstRegDate.Format "02.01.2006" }}
{{ end }}

{{- range $no := .Actions }}
<b>Номер: </b><a href="https://www.opencars.app/number/{{ .Number }}">{{ .Number }}</a>
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>Колір: </b>{{ .Color }}
<b>Тип: </b>{{ .Kind }}
<b>Рік випуску: </b>{{ .Year }}
<b>Повна маса: </b>{{ .TotalWeight }}
<b>Маса без навантаження: </b>{{ .OwnWeight }}
{{- if .Capacity }}
<b>Об'єм двигуна: </b>{{ .Capacity }}
{{- end }}
{{- if .Fuel }}<b>
Тип пального: </b>{{ .Fuel }}
{{- end }}
<b>Дата реєстрації: </b>{{ .Date.Format "02.01.2006" }}
{{ end }}

{{ end }}