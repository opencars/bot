Всього: <b>{{ len .Vehicles }}</b> транспортних засобів

{{ range $vehicle := .Vehicles -}}
<b>{{ .Brand }} {{ .Model }} {{ .Year }}</b>
<b>VIN: </b><a href="https://www.opencars.app/vin/{{ .Vin }}">{{ .Vin }}</a>
<b>Перша реєстрація: </b>{{ .FirstRegDate.Day | printf "%02d" }}.{{ .FirstRegDate.Month | printf "%02d" }}.{{ .FirstRegDate.Year }}

{{ range $no := .Registrations -}}
<b>Номер: </b><a href="https://www.opencars.app/number/{{ .Number }}">{{ .Number }}</a>
<b>Номер документа: </b><a href="https://www.opencars.app/code/{{ .Code }}">{{ .Code }}</a>
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
<b>Категорія: </b>{{ .Category }}
{{- if .NumSeating }}
<b>Кількість сидячих місць: </b>{{ .NumSeating }}
{{- end }}
<b>Дата реєстрації: </b>{{ .Date.Day | printf "%02d" }}.{{ .Date.Month | printf "%02d" }}.{{ .Date.Year }}
{{ end }}

{{ end }}