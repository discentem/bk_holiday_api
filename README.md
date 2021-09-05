## Prompt

In the language of your choosing, create a web server with a single REST API endpoint. This endpoint should accept a list of strings representing dates. For each date, the endpoint will return information about the holiday falling on that date (if a holiday falls on that date).

### Example Requests
#### GET Request to http://localhost:8080/areTheseHolidays/US/2021-07-05,2021-01-01

```shell
Yes 2021-07-05 is a holiday in US
{"date":"2021-07-05","localName":"Independence Day","name":"Independence Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]}
Yes 2021-01-01 is a holiday in US
{"date":"2021-01-01","localName":"New Year's Day","name":"New Year's Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]}
```

#### GET Request on http://localhost:8080/holidays/2021/US

```shell
[{"date":"2021-01-01","localName":"New Year's Day","name":"New Year's Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]},{"date":"2021-01-18","localName":"Martin Luther King, Jr. Day","name":"Martin Luther King, Jr. Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]}

...
```

#### GET request on http://localhost:8080/areTheseHolidaysJSON/US with a JSON blob

```go
uri := fmt.Sprintf("http://localhost:8080/areTheseHolidaysJSON/US", serverURL)
var jsonStr = []byte(`{"dates":["2021-07-05"]}`)
req, err := http.NewRequest("GET", uri, bytes.NewBuffer(jsonStr))
if err != nil {
    logger.Fatal(err)
}
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
    logger.Fatal(err)
}
body, _ := ioutil.ReadAll(resp.Body)
fmt.Println("response Body:", string(body))
```

```shell
response Body: Yes 2021-07-05 is a holiday in US
{"date":"2021-07-05","localName":"Independence Day","name":"Independence Day","countryCode":"US","fixed":false,"global":true,"counties":null,"launchYear":null,"types":["Public"]}
```