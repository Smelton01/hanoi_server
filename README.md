# Hanoi Server

An API server for my Tower of Hanoi App.
Accepts GET requests on `"/"` to return all entries on the scoreboard database.
POST requests with query data: _username_, _game_time_, and _date_ are added to the scoreboard database.

## Usage

```
./hanoi_server -p <port>
```

Start the server on localhost:<port>
