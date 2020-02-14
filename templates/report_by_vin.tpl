Пошук інформації за <b>{{ .Registration.VIN }}</b>

<b>Номер: </b>{{ .Registration.Number }}
<b>Номер документа: </b>{{ .Registration.Code }}
<b>Марка: </b>{{ .Registration.Brand }}
<b>Модель: </b>{{ .Registration.Model }}
<b>VIN: </b> {{ .Registration.VIN }}
<b>Колір: </b>{{ .Registration.Color }}
<b>Тип: </b>{{ .Registration.Kind }}
<b>Рік випуску: </b>{{ .Registration.Year }}
<b>Повна маса: </b>{{ .Registration.TotalWeight }}
<b>Маса без навантаження: </b>{{ .Registration.OwnWeight }}
{{- if .Registration.Capacity }}
<b>Об'єм двигуна: </b>{{ .Registration.Capacity }}
{{- end }}
{{- if .Registration.Fuel }}<b>
Тип пального: </b>{{ .Registration.Fuel }}
{{- end }}
<b>Категорія: </b>{{ .Registration.RankCategory }}
{{- if .Registration.NumStanding }}
<b>Кількість стоячих місць: </b>{{ .Registration.NumStanding }}
{{- end }}
{{- if .Registration.NSeating }}
<b>Кількість сидячих місць: </b>{{ .Registration.NumSeating }}
{{- end }}
<b>Дата першої реєстрації: </b>{{ .Registration.FirstRegDate }}
<b>Дата реєстрації: </b>{{ .Registration.Date }}

<b>Операції</b>
{{ if len .Operations }}
{{ range $i, $operation := .Operations -}}
<b>Номер: </b>{{ .Number }}
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>Колір: </b>{{ .Color }}
<b>Тип: </b>{{ .Kind }} {{ .Body }}
<b>Рік випуску: </b>{{ .Year }}
<b>Повна маса: </b>{{ .TotalWeight }}
<b>Маса без навантаження: </b>{{ .OwnWeight }}
{{ if .Capacity -}}
<b>Об'єм двигуна: </b>{{ .Capacity }}
{{ end -}}
{{ if .Fuel -}}
<b>Тип пального: </b>{{ .Fuel }}
{{ end -}}
<b>Дата реєстрації: </b>{{ .Date }}
<b>Власник: </b>{{ .Person }}
{{ end }}
{{ else }}
Дані за номером {{ .Registration.NRegNew }} не знайдені в реєстраційній інформації з 1 січня 2013 та транспортних ліцензіях.
{{ end }}
