## Использование
`.\<executable-name>` `-config` `<config-path>`

## Методы API
* `POST /create_event`
* `POST /update_event`
* `POST /delete_event`
* `GET /events_for_day`
* `GET /events_for_week`
* `GET /events_for_month`
  Параметры передаются в виде `www-url-form-encoded` (т.е. обычные `user_id=3&date=2019-09-09`).
  В `GET` методах параметры передаются через `queryString`, в `POST` через тело запроса.