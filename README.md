# redash2pagerduty

This app is a redash webhook proxy to pagerduty v2 api.

https://tutorialedge.net/post/golang/creating-restful-api-with-golang/

```
[redash webhook] --(http)--> {http async} --> enqueue

worker
```
