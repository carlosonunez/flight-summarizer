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

### Web Scraping Notice

Functionality provided by this application relies on scraping one or more
websites. The owners of these websites will likely change how they works in the
future and/or might add additional anti-web-scraping defenses that make this
project unviable in the future.

I will do my best to keep this source code current as these changes happen on a
best-available basis. However, you might experience degraded performance or
unavailability depending on where this application runs.

**AS SUCH, THIS APPLICATION AND ITS SOURCE CODE ARE PROVIDED ON AN AS-IS BASIS.
I AM NOT RESPONSIBLE FOR ANY TECHNICAL OR LEGAL ACTIONS THE UPSTREAM PROVIDER
MAY TAKE IN RESPONSE TO THE OPERATION THIS SOFTWARE.**

## Roadmap

### v1

- [X] Retrieve just-in-time takeoff and arrival info for a flight
- [X] Local server mode

### Future

- [ ] Deployment instructions for AWS Lambda
- [ ] Retrieve takeoff and arrival info for a flight on any date
- [ ] Move to [carlosonunez/serverless-stack](https://github.com/carlosonunez/serverless-stack) (upcoming)
- [ ] Support headless browser scraping
- [ ] Summarizer healthchecks
