# golauncher_game
Static Sheep: simple game made using Go, primarily made as a way to learn the language.

Move character around and launch particles, like electrons, around various charged objects to hit a goal

In a rough initial state now, but getting a little better!

Not really worrying about proper error handling, screen/window sizing, or any audio and animations for now, just want to get something working.

Will eventually make executable in broswer via github.io to demo more easily

---

All assets created by me

Code initially inspired/based off guide from https://threedots.tech/post/making-games-in-go/
(actual github: https://github.com/ThreeDotsLabs/meteors/tree/master)

---

Some notes on things that could be improved on (eventually):
- use world coordinates instead of pixels. allows for varying screen size, and making it easier to modify UI
- add animations
- add music
- create better backgrounds, graphics, etc
- break down game package into more discrete packages. it's becoming a catch-all