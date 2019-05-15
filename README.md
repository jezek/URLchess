URLchess - simple anonymus chess played via url
-----------------------------------------------
### Demo
Try at [URLchess project pages](https://jezek.github.io/URLchess).

### Why?
To play chess via mail (or messenger app) without need to register somewhere

### How to play?
- 1st move: Go to [URLchess page](https://jezek.github.io/URLchess), make your move, copy and send generated link to your oponent (via email, messenger, sms, ...).
- Reply to move: Click on link, you got from your oponent, make move, copy and send generated link back.

### Dependencies
- [gopherjs](https://github.com/gopherjs/gopherjs) to generate js
- [Multipurpose chess package for Go/Golang](https://github.com/andrewbackes/chess) for chess logic

### Instalation
Get or update URLchess and dependencies with:
```
go get -u -v github.com/jezek/URLchess
```

Go to URLchess project directory and generate javascript content using:
```
gopherjs build -m
```

Copy `index.html`, `URLchess.js` and `URLchess.css` to your server accessible from internet (or just run it in browser) and feel free to play

### Roadmap
This is an early relase. Improvements will be done soon. Some of them:
- player should be able to ask for draw and if accepted, then draw the game
- player should be able to resign and loose
- previous moves explorer
- minimize js file size (don't use some packages like fmt, ...)
- more touch friendly ui (dragable pieces?)
- multilanguage?
- transition to wasm?
- ...

### Contribution
Test, fork, make pull requests, file issues or improvement ideas, discuss. Everything within reasons is welcome.

### Why golang + gopherjs?
To test golang's web development capabilities and learn something from it. 
