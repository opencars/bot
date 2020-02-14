{{ range $i, $registration := .Registrations -}}
<b>Номер: </b>{{ .Number }}
<b>Номер документа: </b>{{ .Code }}
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>VIN: </b> {{ .VIN }}
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
<b>Категорія: </b>{{ .RankCategory }}
{{- if .NumStanding }}
<b>Кількість стоячих місць: </b>{{ .NumStanding }}
{{- end }}
{{- if .NumSeating }}
<b>Кількість сидячих місць: </b>{{ .NumSeating }}
{{- end }}
<b>Дата першої реєстрації: </b>{{ .FirstRegDate }}
<b>Дата реєстрації: </b>{{ .Date }}

{{ else }}
Дані за номером <b>{{ .Number }}</b> не знайдені.
{{ end }}
