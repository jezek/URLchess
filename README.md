URLchess - simple anonymus chess played via url
-----------------------------------------------

### Dependencies
- [gopherjs](https://github.com/gopherjs/gopherjs) to generate js
- [Multipurpose chess package for Go/Golang](https://github.com/andrewbackes/chess) for chess logic

### Instalation
Get or update URLchess and dependencies with:
```
go get -u github.com/jezek/URLchess
```

Go to URLchess project directory and use
```
gopherjs build -m
```
to generate URLchess.js

Copy `index.html` and `URLchess.js` to your server accessible from internet and feel free to play
