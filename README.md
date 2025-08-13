# pokedexCLI

A colorful, interactive command-line Pokédex inspired by the Pokémon games, built in Go.  
Browse areas, explore wild Pokémon, catch and inspect them, and view your growing collection—all from your terminal, styled with a nostalgic Pokémon flair!

---

## Features

- Explore location areas using live data from the PokéAPI
- Catch wild Pokémon (with real catch odds!)
- Build your personal Pokédex
- Inspect stats, types, and details of your caught Pokémon
- Colorful CLI output inspired by classic game palettes
- Simple REPL interface (just like a game console)

---

## Demo

![Demo](/demo.gif)

*To display your own GIF here:*
1. Create a demonstration GIF of your app (e.g., using [peek](https://github.com/phw/peek) or [asciinema → gif](https://github.com/asciinema/asciicast2gif)).
2. Save it as `demo.gif` in the root of your repository.
3. Git add/commit/push the file.  
   It will auto-display with the markdown above!

---

## How to Run

Make sure you have [Go](https://golang.org/doc/install) installed.

```bash
git clone https://github.com/YOUR_USERNAME/pokedexCLI.git
cd pokedexCLI
go run .
```
---

## How It's Built
- Language: Go
- External APIs: PokeAPI
- Caching: In-memory, with automatic expiry
- Color: fatih/color
- Design: REPL (Read-Eval-Print Loop) dispatches user commands to handlers

---
## License
#### *MIT*
---

*Ready to catch them all in your terminal?*