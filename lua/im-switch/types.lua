---@meta

---@class ImSwitchConfig
---@field binary_path? string Path to the im-switch binary
---@field default_input? string Default input method ID to switch to
---@field auto_switch? boolean Automatically switch to default input in normal mode
---@field auto_restore? boolean Automatically restore previous input in insert mode (experimental)
---@field debug? boolean Enable debug logging

---@class ImSwitch
local ImSwitch = {}

---Initialize the im-switch plugin
---@param opts? ImSwitchConfig Configuration options
function ImSwitch.setup(opts) end

---Switch to the default input method (English)
function ImSwitch.switch_to_english() end

---Restore the previously saved input method
function ImSwitch.restore_input() end

---Get the current input method
---@return string|nil current_input Current input method ID or nil if failed
function ImSwitch.get_current() end

---Set input method to the specified ID
---@param input_id string Input method ID to switch to
---@return boolean success True if successful, false otherwise
function ImSwitch.set_input(input_id) end

---List all available input methods
---@return string[] inputs Array of available input method IDs
function ImSwitch.list_inputs() end

---Enable the plugin
function ImSwitch.enable() end

---Disable the plugin
function ImSwitch.disable() end

---Toggle the plugin state
---@return boolean enabled New enabled state
function ImSwitch.toggle() end

---Check if the plugin is enabled
---@return boolean enabled True if enabled, false otherwise
function ImSwitch.is_enabled() end

return ImSwitch