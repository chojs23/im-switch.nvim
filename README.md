# im-switch.nvim

A Neovim plugin that automatically switches keyboard input method to English when Neovim is focused or in normal mode, and restores the previous input method when entering insert mode.

Perfect for users who type in multiple languages and want seamless input method switching in Neovim.

[![](https://github.com/user-attachments/assets/d7c25f49-ae59-4aaf-866d-1ab29663a018)](https://github.com/user-attachments/assets/d7c25f49-ae59-4aaf-866d-1ab29663a018)

## Features

- **Auto-switch to English** when Neovim gains focus
- **Smart mode switching**: English in normal/command mode, restore in insert mode
- **macOS native**: Uses macOS Text Input Source APIs
- **Linux support**: Works with IBus, Fcitx, Fcitx5, and XKB layouts
- **Windows support**: Switches installed layouts and CJK IME mode

## Installation

### Using Homebrew (CLI)

```bash
brew install chojs23/tap/im-switch
```

### Using [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  "chojs23/im-switch.nvim",
  build = "make build", -- or "make build-wsl-win" for WSL
  config = function()
    require('im-switch').setup({
      -- Configuration options (see below)
    })
  end,
  event = "VeryLazy",
}
```

### Using [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  'chojs23/im-switch.nvim',
  config = function()
    require('im-switch').setup()
  end
  run = 'make build', -- or 'make build-wsl-win' for WSL
}
```

### WSL build note

- In WSL, `make build` creates a Linux binary.
- For Windows IME control from WSL, build the Windows executable with `make build-wsl-win`.
- If needed, set `binary_path` explicitly to your Windows path (for example `/mnt/c/Tools/im-switch.exe`).

## Configuration

```lua
require('im-switch').setup({
  -- Path to the binary (auto-detected if not specified)
  binary_path = 'im-switch',

  -- Default input method ID (platform-specific defaults)
  -- macOS: 'com.apple.keylayout.ABC'
  -- WSL: 'en-US'
  -- Linux: 'us' (XKB), 'xkb:us::eng' (IBus), 'keyboard-us' (Fcitx)
  -- Windows: 'en-US'
  default_input = nil, -- Uses platform default

  -- Auto-switch to default input in normal mode (default: true)
  auto_switch = true,

  -- Auto-restore previous input in insert mode (default: false) (experimental)
  auto_restore = false,

  -- Enable debug logging (default: false)
  debug = false,
})
```

## Usage

The plugin works automatically once installed. However, you can also control it manually:

### Commands

```vim
  :ImSwitchEnable " Enable plugin
  :ImSwitchDisable " Disable plugin
  :ImSwitchToggle " Toggle plugin state
  :ImSwitchStatus " Show current status and input method
```

### Lua API

```lua
local im_switch = require('im-switch')

-- Enable/Disable controls
im_switch.enable()        -- Enable the plugin
im_switch.disable()       -- Disable the plugin
local enabled = im_switch.toggle()  -- Toggle and return new state
local is_on = im_switch.is_enabled() -- Check if enabled

-- Input method controls
im_switch.switch_to_english()  -- Switch to default input
im_switch.restore_input()      -- Restore previous input method

-- Get current input method
local current = im_switch.get_current()
print(current)

-- Set specific input method
im_switch.set_input('com.apple.inputmethod.Korean.2SetKorean')

-- List all available input methods
local inputs = im_switch.list_inputs()
for _, input in ipairs(inputs) do
  print(input)
end
```

## How It Works

The plugin automatically handles these events:

1. **Focus Events**: When Neovim gains focus → switches to default input
2. **Mode Changes**:
   - Normal/Command mode → switches to default input
   - Insert mode (from normal) → restores previous input method (if auto_restore is enabled)

## Finding Input Method IDs

To discover available input method IDs on your system:

```bash
# Build the binary first
make build

# Show runtime status (OS, WSL, binary path, backend, current input source)
./build/im-switch --status

# List all input methods
./build/im-switch -l
```

### Common Input Method IDs

#### macOS

- `com.apple.keylayout.ABC` - US English
- `com.apple.inputmethod.Korean.2SetKorean` - Korean (2-Set)
- `com.apple.inputmethod.Korean` - Korean
- `com.apple.inputmethod.SCIM.ITABC` - Chinese (Simplified)
- `com.apple.inputmethod.TCIM.Cangjie` - Chinese (Traditional)

#### Linux

**XKB Layouts** (setxkbmap):

- `us` - US English
- `gb` - UK English
- `de` - German
- `fr` - French
- `ru` - Russian
- `cn` - Chinese
- `jp` - Japanese
- `kr` - Korean

**IBus Engines**:

- `xkb:us::eng` - US English
- `libpinyin` - Chinese Pinyin
- `anthy` - Japanese
- `hangul` - Korean

**Fcitx/Fcitx5**:

- `keyboard-us` - US English
- `pinyin` - Chinese Pinyin
- `mozc` - Japanese
- `hangul` - Korean

#### Windows

**Keyboard layout names**:

- `en-US` - English (United States)
- `en-GB` - English (United Kingdom)
- `de-DE` - German (Germany)
- `fr-FR` - French (France)
- `ja-JP` - Japanese
- `ko-KR` - Korean
- `zh-CN` - Chinese (Simplified)
- `zh-TW` - Chinese (Traditional)
- `ru-RU` - Russian

**Official layout codes (KLID / hex)**:

- `00000409` - English (United States)
- `00000809` - English (United Kingdom)
- `00000407` - German (Germany)
- `0000040c` - French (France)
- `00000411` - Japanese
- `00000412` - Korean
- `00000804` - Chinese (Simplified)
- `00000404` - Chinese (Traditional)
- `00000419` - Russian

You can also use `0x` style values:

- `0x409` - English (United States)
- `0x411` - Japanese
- `0x412` - Korean
- `0x804` - Chinese (Simplified)

**Windows behavior**:

- `im-switch -l` shows installed keyboard layouts.
- `im-switch [id]` changes IME mode by default.
- Keyboard layout switching is done when you pass an official layout code (KLID/hex).
- Then IME mode is applied:
  - `en-US` requests English mode (IME off)
  - `zh-CN` / `ja-JP` / `ko-KR` request IME on
- For language tags, only the active layout language and `en-US` are accepted.
  - Example: active `zh-CN` allows `zh-CN` and `en-US`, rejects `ja-JP`.
- Official layout codes trigger keyboard layout switching, regardless of active layout language.
  - Example: `00000804` or `0x804` switches to Chinese layout even when active layout is not Chinese.

Examples:

```bash
# List installed layouts
im-switch -l

# Language-tag mode control (active zh-CN layout)
im-switch en-US   # English mode (IME off)
im-switch zh-CN   # Chinese mode (IME on)
im-switch ja-JP   # Error (different language)

# Official layout code switch (always accepted)
im-switch 00000804
im-switch 0x411
```

## Building Manually

If you need to build the binary manually:

```bash
# Build for development
make build

# Build optimized release version
make build-release

# Install system-wide (optional)
make install

# Run tests
make test
```

## Troubleshooting

### Binary not building automatically

- Ensure you have Go installed and CGO enabled
- Check that you have Xcode command line tools: `xcode-select --install`
- Try building manually: `make build`

### Plugin not switching inputs

1. Enable debug mode to see what's happening:
   ```lua
   require('im-switch').setup({ debug = true })
   ```
2. Check available input methods: `./build/im-switch -l`
3. Verify your `default_input` setting matches an available input method

### Permission errors

- The plugin only reads/writes input methods, no special permissions needed
- If you see permission errors, try rebuilding: `make clean && make build`

## Requirements

### macOS

- **Neovim** (uses Neovim-specific APIs)
- **macOS** (uses macOS Text Input Source framework)
- **Go 1.19+** (for building the binary)
- **CGO enabled** (uses macOS system APIs)
- **Xcode Command Line Tools** (`xcode-select --install`)

### Linux

- **Neovim** (uses Neovim-specific APIs)
- **Go 1.19+** (for building the binary)
- **Input Method Framework**: One of:
  - IBus (`ibus-daemon`)
  - Fcitx (`fcitx`)
  - Fcitx5 (`fcitx5`)
  - XKB (setxkbmap - built into X11/Wayland)

### Windows

- **Neovim** (uses Neovim-specific APIs)
- **Windows 11** (switches installed layouts and CJK IME mode)
- **Go 1.19+** (for building the binary)

### WSL

- **Neovim in WSL**
- **Go 1.19+ in WSL** (to build Windows binary)
- A Windows path mounted in WSL (for example `/mnt/c/Tools`)

For common WSL + Windows Terminal usage, use a Windows executable from WSL.

```bash
# Build a Windows binary from WSL into Windows filesystem
make build-wsl-win

# Verify from WSL
/mnt/c/Tools/im-switch.exe
/mnt/c/Tools/im-switch.exe --status
/mnt/c/Tools/im-switch.exe en-US
/mnt/c/Tools/im-switch.exe 00000804
/mnt/c/Tools/im-switch.exe -l

# Or run the built-in verification target
make test-wsl-win
```

WSL runtime notes:

- `--status` prints `OS: windows` when running the Windows binary.
- `WSL: true` is reported when the process is detected as being launched from WSL.
- Layout changes still depend on layouts installed in Windows.

You can override the output path:

```bash
make build-wsl-win WSL_WIN_OUT=/mnt/c/Users/<you>/bin/im-switch.exe
```

WSL auto-detection prefers these paths when `binary_path` is not set:

- `/mnt/c/Tools/im-switch.exe`
- `/mnt/c/Windows/im-switch.exe`
- `/mnt/c/Program Files/im-switch/im-switch.exe`

You can still set it explicitly:

```lua
require('im-switch').setup({
  binary_path = '/mnt/c/Tools/im-switch.exe',
  default_input = 'en-US',
})
```

Minimal WSL setup (copy and paste):

```lua
{
  "chojs23/im-switch.nvim",
  build = "make build-wsl-win",
  opts = {
    binary_path = "/mnt/c/Tools/im-switch.exe",
    default_input = "en-US",
  },
}
```

## License

MIT License - see LICENSE file for details.

## Contributing

Issues and feature requests are welcome
