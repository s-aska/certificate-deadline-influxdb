Certificate deadline influxdb
=============

Write the Deadline of SSL certificate to InfluxDB.

#### Docker

```sh
docker pull aska/certificate-deadline-influxdb

docker run \
    -e DOMAINS="example.com,example.org" \
    -e INFLUXDB_WRITE_URL="http://example.com:8086/write?db=mydb" \
    --publish 8080:8080 \
    --name test \
    --rm aska/certificate-deadline-influxdb
```

#### Arukas

- Image: aska/certificate-deadline-influxdb:latest
- Port: 8080
- ENV:
    - DOMAINS
    - INFLUXDB_WRITE_URL 

#### Grafana

```sh
SELECT mean("value") FROM "deadline" WHERE $timeFilter GROUP BY time($interval), "domain" fill(null)
```
