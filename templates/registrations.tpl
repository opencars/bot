{{ range $i, $registration := .Registrations -}}
<b>Номер: </b>{{ .NRegNew }}
<b>Номер документа: </b>{{ .SDoc }}{{ .NDoc }}
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>VIN: </b> {{ .VIN }}
<b>Колір: </b>{{ .Color }}
<b>Тип: </b>{{ .Kind }}
<b>Рік випуску: </b>{{ .MakeYear }}
<b>Повна маса: </b>{{ .TotalWeight }}
<b>Маса без навантаження: </b>{{ .OwnWeight }}
{{- if .Capacity }}
<b>Об'єм двигуна: </b>{{ .Capacity }}
{{- end }}
{{- if .Fuel }}<b>
Тип пального: </b>{{ .Fuel }}
{{- end }}
<b>Категорія: </b>{{ .RankCategory }}
{{- if .NStanding }}
<b>Кількість стоячих місць: </b>{{ .NStanding }}
{{- end }}
{{- if .NSeating }}
<b>Кількість сидячих місць: </b>{{ .NSeating }}
{{- end }}
<b>Дата першої реєстрації: </b>{{ .DFirstReg }}
<b>Дата реєстрації: </b>{{ .DReg }}

{{ else }}
Дані за номером <b>{{ .Number }}</b> не знайдені.
{{ end }}
