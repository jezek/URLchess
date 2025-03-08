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

### Web application files

These files are needed to run the app: `index.*.html` and `assets/*`.

To run locally open the `index.html` in a browser.

Although hhis runs only the js version, which is slower than wasm version.

To run the wasm version, you'll need to run it through some server. E.g. run python's file server `$ python3 -m http.server 8080` in the URLchess directory and open `http://0.0.0.0:8080/` in browser.

### Building assets
Go to URLchess project directory and generate all assets content using:
```
go generate
```

This generates `assets/URLchess.js`, `assets/URLchess.js.map`, `assets/URLchess.wasm`, `assets/URLchess.tinygo.wasm` and copies `assets/wasm_exec.js` and `assets/wasm_exec.tinygo.js` from local sources.

You'll need [`gopherjs`](https://github.com/gopherjs/gopherjs) and [`tinygo`](https://tinygo.org/) to have installed to successfully generate all assets.


### Roadmap
This is an early relase. Improvements will be done soon. Some of them:
- player should be able to ask for draw and if accepted, then draw the game
- player should be able to resign and loose
- multilanguage?
- use smaller tinygo wasm as main wasm version
- import a game
- 2 player mode
- ...

### Contribution
Test, fork, make pull requests, file issues or improvement ideas, discuss. Everything within reasons is welcome.

### Why golang + wasm/gopherjs?
To test golang's web development capabilities and learn something from it. 
