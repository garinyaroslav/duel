# Duel

A networked top-down duel game for two players. Written in Go using [Ebitengine](https://ebitengine.org/). One player starts the server, the other connects to it by IP.

![Game screenshot](docs/screenshot.png)

## Rules

- Two players on a symmetrical map with static walls.

## Controls

| Action | Keys                                       |
| ------ | ------------------------------------------ |
| Move   | `W` `A` `S` `D` or arrow keys              |
| Aim    | mouse                                      |
| Shoot  | `Space` or left mouse button (can be held) |

## Running

Requires Go 1.21 or newer. On Linux, Ebitengine needs some system libraries; the list is in the [documentation](https://ebitengine.org/en/documents/install.html).

```bash
git clone https://github.com/garinyaroslav/duel.git
cd duel
make all
```
