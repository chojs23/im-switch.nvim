# im-switch.nvim

A Neovim plugin that automatically switches keyboard input method to English when Neovim is focused or in normal mode, and restores the previous input method when entering insert mode.

Perfect for users who type in multiple languages and want seamless input method switching in Neovim.

## Features

- **Auto-switch to English** when Neovim gains focus
- **Smart mode switching**: English in normal/command mode, restore in insert mode
- **macOS native**: Uses macOS Text Input Source APIs
- **Linux support**: Works with IBus, Fcitx, Fcitx5, and XKB layouts

## Installation

### Using [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  "your-username/im-switch.nvim", -- Replace with your repo path
  build = "make build", -- Async build handled by LazyVim
  config = function()
    require('im-switch').setup({
      -- Configuration options (see below)
    })
  end,
  event = "VeryLazy", -- Load after UI is ready
}
```

### Using [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  'your-username/im-switch.nvim', -- Replace with your repo path
  config = function()
    require('im-switch').setup()
  end
}
```

## Configuration

```lua
require('im-switch').setup({
  -- Path to the binary (auto-detected if not specified)
  binary_path = 'im-switch',

  -- Default input method ID (platform-specific defaults)
  -- macOS: 'com.apple.keylayout.ABC'
  -- Linux: 'us' (XKB), 'xkb:us::eng' (IBus), 'keyboard-us' (Fcitx)
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
3. **Session Management**:
   - Startup → switches to default input

## Finding Input Method IDs

To discover available input method IDs on your system:

```bash
# Build the binary first
make build

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

## How It Differs from Other Solutions

- **Cross-platform**: Works on both macOS and Linux
- **Native Integration**: Uses platform-specific APIs (macOS) and tools (Linux)
- **Smart Restoration**: Remembers and restores your previous input method
- **Focus Aware**: Handles window focus changes intelligently
- **Auto-detection**: Automatically detects and works with available input method frameworks

## License

MIT License - see LICENSE file for details.

## Contributing

Issues and feature requests are welcome
