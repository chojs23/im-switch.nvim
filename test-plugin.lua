-- Quick test script for the nvim-hangul plugin
-- Run with: nvim --headless -c "luafile test-plugin.lua" -c "qa"

package.path = package.path .. ";./lua/?.lua;./lua/?/init.lua"

local nvim_hangul = require("nvim-hangul-go.lua.nvim-hangul.init")

print("Testing nvim-hangul plugin...")

-- Test setup with debug enabled
nvim_hangul.setup({
	binary_path = "./nvim-hangul-go",
	debug = true,
})

-- Test getting current input
local current = nvim_hangul.get_current()
print("Current input method: " .. (current or "nil"))

-- Test listing inputs
local inputs = nvim_hangul.list_inputs()
print("Available input methods:")
for i, input in ipairs(inputs) do
	print("  " .. i .. ": " .. input)
end

-- Test switching to English
print("Switching to English...")
nvim_hangul.switch_to_english()

-- Check if it switched
local new_current = nvim_hangul.get_current()
print("After switch: " .. (new_current or "nil"))

print("Plugin test completed!")
