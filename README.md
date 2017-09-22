# manas-food-api
Scraper API for Manas University Cafeteria

## Usage

### Run the service

```bash
PORT=8000 go run main.go #by default port is 8080
```

### Call Endpoints

Get all available menus

```bash
curl -i -H "Accept: application/json" 
        -H "Content-Type: application/json" 
        -X GET  http://localhost:8000/menus
```

Get menu by date

```bash
curl -i -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -X GET  http://localhost:8000/menus/2017-09-28
```

### Deploy on Heroku

You can deploy this service to Heroku. Just create app on Heroku and use it!
