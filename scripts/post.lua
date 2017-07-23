wrk.method = "POST"
wrk.body = [[{
  "url_base": "http://localhost:8080/",
  "event": "alert_state_change",
  "alert": {
    "state": "triggered",
    "query_id": 1,
    "name": "fuga: count(*) greater than 5000",
    "rearm": null,
    "updated_at": "2015-12-11T08:08:00.123456",
    "user_id": 10,
    "created_at": "2017-07-23T01:00:10.000000+00:00",
    "last_triggered_at": "2017-07-23T01:02:03.123456+00:00",
    "id": 5,
    "options": {
      "column": "count(*)",
      "value": 5000,
      "op": "greater than"
    }
  }
}]]
wrk.headers["Content-Type"] = "application/json"
