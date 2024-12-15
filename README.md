# Flight Summarizer ✈️

Returns summarized information about a current flight. Useful for simple
integrations with other apps or services, like iOS Shortcuts.

```sh
# Summarize a flight happening today.
# $: curl https://flight-summarizer.example/summarize?flightNumber=AAL1
{
  "flight_number": "AAL1",
  "origin": "JFK",
  "destination": "LAX",
  "times": {
     "scheduled": {
        "departure_time": "2024-12-15T06:58:00-05:00",
        "arrival_time": "2024-12-15T10:20:00-08:00",
     },
     "actual" {
        "departure_time": "2024-12-15T07:10:00-05:00",
        "arrival_time": "2024-12-15T09:58:00-08:00",
     }
   }
}
```

## Roadmap

### v1

[ ] Retrieve just-in-time takeoff and arrival info for a flight
[ ] Deployment instructions for AWS Lambda

### Future

[ ] Local server mode
[ ] Retrieve takeoff and arrival info for a flight on any date
[ ] Move to [carlosonunez/serverless-stack](https://github.com/carlosonunez/serverless-stack) (upcoming)
[ ] Support headless browser scraping
