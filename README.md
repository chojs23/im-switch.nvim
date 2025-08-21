# im-switch.nvim

A Neovim plugin that automatically switches keyboard input method to English when Neovim is focused or in normal mode, and restores the previous input method when entering insert mode.

Perfect for users who type in multiple languages and want seamless input method switching in Neovim.

## Features

- 🎯 **Auto-switch to English** when Neovim gains focus
- 🔄 **Smart mode switching**: English in normal/command mode, restore in insert mode
- 🛠️ **Auto-build**: Binary builds automatically on first plugin load
- ⚙️ **Configurable**: Custom input methods and behavior
- 🐛 **Debug mode** for troubleshooting
- 🍎 **macOS native**: Uses macOS Text Input Source APIs

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

  -- Default input method ID (default: US English)
  default_input = 'com.apple.keylayout.ABC',

  -- Auto-switch to default input in normal mode (default: true)
  auto_switch = true,

  -- Auto-restore previous input in insert mode (default: true)
  auto_restore = true,

  -- Enable debug logging (default: false)
  debug = false,
})
```

## Usage

The plugin works automatically once installed. However, you can also control it manually:

```lua
local im_switch = require('im-switch')

-- Switch to default input
im_switch.switch_to_english()

-- Restore previous input method
im_switch.restore_input()

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
   - Insert mode (from normal) → restores previous input method
3. **Session Management**:
   - Startup → switches to default input
   - Exit → restores previous input method

## Finding Input Method IDs

To discover available input method IDs on your system:

```bash
# Build the binary first
make build

# List all input methods
./build/im-switch -l
```

### Common macOS Input Methods

- `com.apple.keylayout.ABC` - US English
- `com.apple.inputmethod.Korean.2SetKorean` - Korean (2-Set)
- `com.apple.inputmethod.Korean` - Korean
- `com.apple.inputmethod.SCIM.ITABC` - Chinese (Simplified)
- `com.apple.inputmethod.TCIM.Cangjie` - Chinese (Traditional)

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

- **Neovim** (uses Neovim-specific APIs)
- **macOS** (uses macOS Text Input Source framework)
- **Go 1.19+** (for building the binary)
- **CGO enabled** (uses macOS system APIs)

## How It Differs from Other Solutions

- **Native Integration**: Uses macOS APIs directly, no external dependencies
- **Smart Restoration**: Remembers and restores your previous input method
- **Focus Aware**: Handles window focus changes intelligently
- **Auto-build**: No manual compilation needed
- **Lightweight**: Single binary with minimal overhead

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

Issues and feature requests are welcome
