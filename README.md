# Recurlit

Recurlit (recURLit) is a lightweight terminal UI for building, executing, and saving HTTP requests — a keyboard-driven alternative to Postman that lives entirely in your terminal. Inspired by lazygit's pane-based UX and vim keybindings, it keeps only the features you actually need: request building, live curl preview, execution, and save/load.

Construct a curl commands once, recurlit anytime.

---

## Features

- **7 HTTP methods** — GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Headers & query params** — key-value pair editors with add/delete support
- **Request body** — multi-line textarea for JSON, form data, etc.
- **Authentication** — None, Bearer token, or Basic auth (username/password)
- **Live curl preview** — auto-generated curl command updates as you type
- **Response viewer** — status code, duration, headers, and pretty-printed JSON body
- **Save & load requests** — persist requests by name to `~/.config/recurlit/requests/`
- **Help overlay** — press `?` at any time to see all keybindings
- **Vim-inspired navigation** — modal editing with normal/insert modes

---

## Building

Requires **Go 1.21+**. The project uses [BubbleTea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss) from [Charm](https://charm.land/).

```bash
git clone https://github.com/bugsy6/recurlit.git
cd recurlit
go build
```

This produces a `recurlit` binary in the current directory. Optionally install it to your PATH:

```bash
go install
```

---

## Running

```bash
./recurlit
```

Or if installed:

```bash
recurlit
```

---

## Usage

Recurlit uses a **modal editing model** similar to vim. On launch you start in **Normal mode**, where you navigate between panes and trigger actions. Press `i` or `Enter` on a pane to enter **Insert mode** and edit its content. Press `Esc` to return to Normal mode.

### Typical workflow

1. Tab to the **Method** pane — press `i`, then `j`/`k` to select a method.
2. Tab to the **URL** pane — press `i` and type your endpoint.
3. Tab to **Headers** or **Params** to add key-value pairs as needed.
4. Tab to **Body** to enter a request body (JSON, etc.).
5. Tab to **Auth** to configure authentication if required.
6. Press `r` (in Normal mode) to execute the request.
7. Tab to **Results** to scroll through the response.
8. Press `s` to save the request, or `o` to load a previously saved one.

---

## Controls

### Normal Mode (navigation & actions)

| Key               | Action                                |
| ----------------- | ------------------------------------- |
| `Tab` / `l`       | Focus next pane                       |
| `Shift+Tab` / `h` | Focus previous pane                   |
| `i` / `Enter`     | Enter Insert mode for the active pane |
| `Esc`             | Return to Normal mode                 |
| `r`               | Execute the request                   |
| `s`               | Save request (opens Save modal)       |
| `o`               | Open saved request (opens Load modal) |
| `?`               | Toggle help overlay                   |
| `q` / `Ctrl+C`    | Quit                                  |

### Method Pane (Insert mode)

| Key       | Action                            |
| --------- | --------------------------------- |
| `j` / `↓` | Cycle to next HTTP method         |
| `k` / `↑` | Cycle to previous HTTP method     |
| `Enter`   | Confirm and return to Normal mode |

### Headers & Params Panes (Insert mode)

| Key         | Action                             |
| ----------- | ---------------------------------- |
| `Tab`       | Next field (key → value → new row) |
| `Shift+Tab` | Previous field                     |
| `Enter`     | Confirm field / add new row        |
| `↑` / `↓`   | Move between rows                  |
| `Ctrl+D`    | Delete current row                 |

### Auth Pane (Insert mode)

| Key      | Action                                         |
| -------- | ---------------------------------------------- |
| `Tab`    | Cycle through fields for the current auth type |
| `Ctrl+T` | Cycle auth type (None → Bearer → Basic)        |

### Results Pane (Normal mode)

| Key       | Action                |
| --------- | --------------------- |
| `j` / `↓` | Scroll down           |
| `k` / `↑` | Scroll up             |
| `g`       | Jump to top           |
| `G`       | Jump to bottom        |
| `Ctrl+d`  | Scroll half page down |
| `Ctrl+u`  | Scroll half page up   |

---

## Feedback & Contributing

If you run into bugs, have feature requests, or just want to say something nice, please open an issue — all feedback is welcome. If you'd like to contribute a fix or a new feature, pull requests are very much welcomed and appreciated.
