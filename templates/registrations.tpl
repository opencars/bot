{{ if len .Registrations }}
{{ range $i, $registration := .Registrations -}}
<b>Номер: </b><a href="https://www.opencars.app/number/{{ .Number }}">{{ .Number }}</a>
<b>Номер документа: </b><a href="https://www.opencars.app/code/{{ .Code }}">{{ .Code }}</a>
<b>Марка: </b>{{ .Brand }}
<b>Модель: </b>{{ .Model }}
<b>VIN: </b><a href="https://www.opencars.app/vin/{{ .VIN }}">{{ .VIN }}</a>
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

{{ end }}
{{ else }}
{{ if .Number }}
Дані за номером <b>{{ .Number }}</b> не знайдені.
{{ end }}
{{ if .Code }}
Дані за номером свідоцтва про реєстрацію <b>{{ .Code }}</b> не знайдені.
{{ end }}
{{ end }}
