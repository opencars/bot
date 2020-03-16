Всього: <b>{{ .Amount }}</b> транспортних засобів
Нових: <b>{{ len .Cars }}</b> транспортних засобів

{{ range $i, $car := .Cars -}}
<b>Модель: </b>{{ .MarkName }} {{ .ModelName }}
<b>Рік випуску: </b>{{ .Car.Year }}
<b>Місто: </b>{{ .LocationCityName }}
<b>Пробіг: </b>{{ .Car.Race }}
<b>Трансміссія: </b>{{ .Car.GearboxName }}
<b>Паливо: </b>{{ .Car.FuelName }}
<b>Перевірка: </b>/auto_{{ .Car.AutoID }}
<a href="{{ .LinkToView }}">Детальніше</a>

{{ end }}