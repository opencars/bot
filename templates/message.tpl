<b>{{ len . }}</b> автомобилей найдено!

{{ range $i, $car := . -}}
<b>Автомобиль: </b>{{ .MarkName }} {{ .ModelName }} {{ .Car.Year }}
<b>Город: </b>{{ .LocationCityName }}
<b>Пробег: </b>{{ .Car.Race }}
<b>Коробка: </b>{{ .Car.GearboxName }}
<b>Топливо: </b>{{ .Car.FuelName }}
<b>Детальнее: </b>/auto_{{ .Car.AutoID }}
<a href="{{ .LinkToView }}">Сайт</a>

{{ end }}
