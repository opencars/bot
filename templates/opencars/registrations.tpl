{{ if len .Registrations }}
{{ range $i, $registration := .Registrations -}}
<b>Номер свідоцтва: </b>{{ .Code }}

<b>Реєстраційний номер: </b>{{ .Number }}
<b>Дата першої реєстрації: </b>{{ .FirstReg }}
<b>Дата реєстрації: </b>{{ .Date }}
<b>Рік випуску: </b>{{ .Year }}
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>Тип: </b>{{ .Kind }}
<b>Маса без навантаження: </b>{{ .OwnWeight }}
<b>Категорія: </b>{{ .Category }}
<b>Об'єм двигуна: </b>{{ .Capacity }}
<b>Тип палива: </b>{{ .Fuel }}
<b>Колір: </b>{{ .Color }}
<b>VIN-код: </b> {{ .VIN }}

{{ end }}
{{ else }}
Дані за номером {{ .Code }} не знайдені.
{{ end }}
