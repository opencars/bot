<b>{{ len . }}</b> автомобилей добавлено!

{{ range $i, $car := . -}}
<b>Автомобиль: </b><code>{{ .MarkName }} {{ .ModelName }} {{ .Car.Year }}</code>
<b>Город: </b><code>{{ .LocationCityName }}</code>
<b>Пробег: </b><code>{{ .Car.Race }}</code>
<b>Коробка: </b><code>{{ .Car.GearboxName }}</code>
<b>Топливо: </b><code>{{ .Car.FuelName }}</code>
<a href="{{ .LinkToView }}">Детали</a>

{{ end }}