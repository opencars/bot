<b>{{ len . }}</b> транспортних засобів знайдено.

{{ range $i, $car := . -}}
<b>Модель: </b>{{ .MarkName }} {{ .ModelName }}
<b>Рік Випуску: </b>{{ .Car.Year }}
<b>Місто: </b>{{ .LocationCityName }}
<b>Пробіг: </b>{{ .Car.Race }}
<b>Трансміссія: </b>{{ .Car.GearboxName }}
<b>Паливо: </b>{{ .Car.FuelName }}
<b>Аналіз: </b>/auto_{{ .Car.AutoID }}
<a href="{{ .LinkToView }}">Детальніше</a>

{{ end }}