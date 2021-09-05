## Prompt

In the language of your choosing, create a web server with a single REST API endpoint. This endpoint should accept a JSON list representing dates. For each date, the endpoint will return information about the holiday falling on that date (if a holiday falls on that date).

### Example Implementation

See `main.go` for an example implementation of the above prompt in Golang.

### Example Requests and Responses
#### GET request on http://localhost:8080/areTheseHolidaysJSON/US with the following JSON: `{"dates":["2021-07-05"]}`

```shell
response Body: Yes 2021-07-05 is a holiday in US
{"date":"2021-07-05","localName":"Independence Day","name":"Independence Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]}
```

#### GET request on http://localhost:8080/areTheseHolidaysJSON/US with the following JSON: `{"dates":["2021-07-06"]}`

```shell
no holidays on 2021-07-06 in US
```