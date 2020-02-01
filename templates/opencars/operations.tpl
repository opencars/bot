{{ if len .Operations }}
Реєстрації транспортних засобів за номером <b>{{ .Number }}</b> — {{ len .Operations }}

{{ range $i, $operation := .Operations -}}
<b>Номер: </b>{{ .Number }}
<b>Марка: </b>{{ .Make }}
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
{{ if eq .Person "P" -}}
<b>Власник: </b>Фізична особа
{{- else if eq .Person "J" -}}
<b>Власник: </b>Юридична особа
{{ end }}

{{ end }}
{{ else }}
Дані за номером {{ .Number }} не знайдені в реєстраційній інформації з 1 січня 2013 та транспортних ліцензіях.
{{ end }}
