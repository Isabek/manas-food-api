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

```json
[
  {
    "date": "2017-09-11",
    "foods": [
      {
        "name": "Süzme Mercimek Çorbası",
        "calories": "176"
      },
      {
        "name": "Et Sote",
        "calories": "370"
      },
      ...
    ],
    "total_calories": "1020"
  },
  {
    "date": "2017-09-12",
    "foods": [
      {
        "name": "Yayla Çorbası",
        "calories": "175"
      },
      {
        "name": "Fırın Köfte-Fıri",
        "calories": "241"
      },
      ...
    ],
    "total_calories": "791"
  },
  ...
]
```

Get menu by date

```bash
curl -i -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -X GET  http://localhost:8000/menus/2017-09-28
```

```json
{
  "date": "2017-09-28",
  "foods": [
    {
      "name": "Taneli Sebze Çorbası",
      "calories": "125"
    },
    {
      "name": "Sebzeli Kebap",
      "calories": "266"
    },
    {
      "name": "Makarna",
      "calories": "360"
    },
    {
      "name": "Üzüm",
      "calories": "71"
    },
    {
      "name": "Ekmek (1 Dilim)",
      "calories": "82"
    }
  ],
  "total_calories": "822"
}
```





### Deploying on Heroku

You can deploy this service on Heroku. Just create app on Heroku and use it!
