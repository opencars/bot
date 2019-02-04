{{ if .Number }}
{{ if len .Cars }}
{{ range $i, $car := .Cars -}}
<b>Номер: </b>{{ .Number }}
<b>Модель: </b>{{ .Model }}
<b>Рік випуску: </b>{{ .Year }}
<b>Дата: </b>{{ .Date }}
<b>Колір: </b>{{ .Color }}
<b>Об'єм двигуна: </b>{{ .Capacity }}
<b>Вага: </b>{{ .OwnWeight }}
<b>Тип авто: </b>{{ .Kind }}
<b>Тип кузова: </b>{{ .Body }}
{{ end }}
{{ else }}
Інформацію щодо автомобілів з номером {{ .Number }} не знайдено!
{{ end }}
{{ else }}
Помилковий запит!
{{ end }}
