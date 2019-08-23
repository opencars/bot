<b>{{ .Package }}</b>

<b>{{ .Resource.Name }}</b>
<b>Номер</b>: {{ .Resource.ID }}
<b>Формат</b>: {{ .Resource.Format }}

Дані були востаннє оновлені об {{ .Resource.LastModified.Format "15:04:05" }}.

<a href="{{ .Resource.URL }}">Завантажити</a>
