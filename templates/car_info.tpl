{{ if len .Operations }}
Реєстрації транспортних засобів за номером <b>{{ .Number }}</b> — {{ len .Operations }}

{{ range $i, $operation := .Operations -}}
<b>Номер: </b>{{ .Number }}
<b>Модель: </b>{{ .Brand }} {{ .Model }}
<b>Рік випуску: </b>{{ .Year }}
<b>Дата: </b>{{ .Date }}
<b>Колір: </b>{{ .Color }}
<b>Об'єм двигуна: </b>{{ .Capacity }}
<b>Вага: </b>{{ .Weight }}
<b>Тип авто: </b>{{ .Kind }}
<b>Тип кузова: </b>{{ .Body }}

{{ end }}
{{ else }}
Дані за номером {{ .Number }} не знайдені в реєстраційній інформації з 1 січня 2013 та транспортних ліцензіях.
{{ end }}
