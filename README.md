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

Go to URLchess project directory and generate javascript content using:
```
gopherjs build -m
```

Copy `index.html` and `URLchess.js` to your server accessible from internet and feel free to play

### Demo
Try at [URLchess project pages](https://jezek.github.io/URLchess).

### Why?
To play chess via mail (or messanger app) without need to register somewhere

### How to play?
Go to URLchess page, make your move and send generated link to your oponent.

### Roadmap
This is an very early relase. Improvements will be done soon. Some of them:
- previous moves explorer
- chess board orientation depending on moving player (really is it good?)
- move by clicking (or dragging) piece
- next move arrow
- copy generated link to clipboard
- ...

### Contribution
Test, fork, make pull requests, file issues or improvement ideas, discuss. Everything within reasons is welcome.
