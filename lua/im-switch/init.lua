---@type ImSwitch
local M = {}

---@return ImSwitchConfig
local function get_default_config()
	local is_mac = vim.fn.has("mac") == 1 or vim.fn.has("macunix") == 1
	local is_linux = vim.fn.has("unix") == 1 and not is_mac

	if is_mac then
		return {
			binary_path = "im-switch",
			default_input = "com.apple.keylayout.ABC",
			auto_switch = true,
			auto_restore = false,
			debug = false,
		}
	elseif is_linux then
		return {
			binary_path = "im-switch",
			default_input = "us",
			auto_switch = true,
			auto_restore = false,
			debug = false,
		}
	else -- TODO: window support
		-- Fallback for other systems
		return {
			binary_path = "im-switch",
			default_input = "us",
			auto_switch = true,
			auto_restore = false,
			debug = false,
		}
	end
end

---@type ImSwitchConfig
local config = get_default_config()

---@type string|nil
local saved_input = nil
---@type string|nil
local last_mode = nil
---@type boolean
local enabled = true

---@param msg string
local function log(msg)
	if config.debug then
		print("[im-switch] " .. msg)
	end
end

---@param args string|nil Command arguments
---@return string|nil result Command output or nil if failed
local function execute_command(args)
	local cmd = config.binary_path
	if args then
		cmd = cmd .. " " .. args
	end

	local handle = io.popen(cmd .. " 2>&1")
	if not handle then
		log("Failed to execute: " .. cmd)
		return nil
	end

	local result = handle:read("*a")
	local success = handle:close()

	if not success then
		log("Command failed: " .. cmd)
		return nil
	end

	return result:gsub("%s+$", "")
end

---@return string|nil current_input Current input method ID or nil if failed
local function get_current_input()
	return execute_command()
end

---@param input_id string Input method ID to switch to
---@return boolean success True if successful, false otherwise
local function set_input(input_id)
	if input_id and input_id ~= "" then
		local result = execute_command(input_id)
		log("Switched to: " .. input_id)
		return result ~= nil
	end
	return false
end

local function switch_to_default()
	if enabled and config.auto_switch then
		local current = get_current_input()
		if current and current ~= config.default_input then
			saved_input = current
			log("Saving current input: " .. current)
			set_input(config.default_input)
		end
	end
end

local function restore_input()
	if enabled and config.auto_restore and saved_input then
		log("Restoring input: " .. saved_input)
		set_input(saved_input)
		saved_input = nil
	end
end

local function on_mode_changed()
	if not enabled then
		return
	end

	local mode = vim.fn.mode()

	if mode == "n" or mode == "c" then
		switch_to_default()
	elseif mode == "i" and last_mode == "n" then
		restore_input()
	end

	last_mode = mode
end

local function on_focus_gained()
	if enabled then
		switch_to_default()
	end
end

local function on_focus_lost()
	if not enabled then
		return
	end

	local current = get_current_input()
	if current and current ~= config.default_input then
		saved_input = current
	end
end

local function setup_autocmds()
	local group = vim.api.nvim_create_augroup("ImSwitch", { clear = true })

	vim.api.nvim_create_autocmd({ "ModeChanged" }, {
		group = group,
		callback = on_mode_changed,
	})

	vim.api.nvim_create_autocmd({ "FocusGained" }, {
		group = group,
		callback = on_focus_gained,
	})

	vim.api.nvim_create_autocmd({ "FocusLost" }, {
		group = group,
		callback = on_focus_lost,
	})

	vim.api.nvim_create_autocmd({ "VimEnter" }, {
		group = group,
		callback = switch_to_default,
	})
end

---Initialize the im-switch plugin
---@param opts? ImSwitchConfig Configuration options
function M.setup(opts)
	opts = opts or {}
	config = vim.tbl_deep_extend("force", config, opts)

	if not opts or not opts.binary_path then
		local plugin_dir = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h:h:h")
		local local_binary = plugin_dir .. "/build/im-switch"
		local system_binary = "/usr/local/bin/im-switch"

		if vim.fn.executable(local_binary) == 1 then
			config.binary_path = local_binary
		elseif vim.fn.executable(system_binary) == 1 then
			config.binary_path = system_binary
		else
			config.binary_path = "im-switch" -- fallback to PATH
		end
	end

	if vim.fn.executable(config.binary_path) ~= 1 then
		vim.notify("[im-switch] Binary not found: " .. config.binary_path, vim.log.levels.WARN)
		return
	end

	local current = get_current_input()
	if not current or current == "" then
		vim.notify("[im-switch] Failed to get current input method", vim.log.levels.WARN)
		return
	end

	log("Current input method: " .. current)
	log("Default input method: " .. config.default_input)

	setup_autocmds()

	-- Create user commands
	vim.api.nvim_create_user_command("ImSwitchEnable", function()
		enabled = true
		vim.notify("[im-switch] Enabled", vim.log.levels.INFO)
		log("Plugin enabled")
		-- Switch to default immediately when enabled
		switch_to_default()
	end, { desc = "Enable im-switch plugin" })

	vim.api.nvim_create_user_command("ImSwitchDisable", function()
		enabled = false
		vim.notify("[im-switch] Disabled", vim.log.levels.INFO)
		log("Plugin disabled")
	end, { desc = "Disable im-switch plugin" })

	vim.api.nvim_create_user_command("ImSwitchToggle", function()
		enabled = not enabled
		local status = enabled and "Enabled" or "Disabled"
		vim.notify("[im-switch] " .. status, vim.log.levels.INFO)
		log("Plugin " .. status:lower())
		if enabled then
			switch_to_default()
		end
	end, { desc = "Toggle im-switch plugin" })

	vim.api.nvim_create_user_command("ImSwitchStatus", function()
		local status = enabled and "Enabled" or "Disabled"
		local current = get_current_input()
		local msg = string.format("[im-switch] Status: %s | Current input: %s", status, current or "Unknown")
		vim.notify(msg, vim.log.levels.INFO)
	end, { desc = "Show im-switch plugin status" })

	log("im-switch plugin initialized")
end

function M.switch_to_english()
	switch_to_default()
end

function M.restore_input()
	restore_input()
end

---Get the current input method
---@return string|nil current_input Current input method ID or nil if failed
function M.get_current()
	return get_current_input()
end

---Set input method to the specified ID
---@param input_id string Input method ID to switch to
---@return boolean success True if successful, false otherwise
function M.set_input(input_id)
	return set_input(input_id)
end

---List all available input methods
---@return string[] inputs Array of available input method IDs
function M.list_inputs()
	local result = execute_command("-l")
	if result then
		return vim.split(result, "\n")
	end
	return {}
end

---Enable the plugin
function M.enable()
	enabled = true
	log("Plugin enabled via API")
	switch_to_default()
end

---Disable the plugin
function M.disable()
	enabled = false
	log("Plugin disabled via API")
end

---Toggle the plugin state
---@return boolean enabled New enabled state
function M.toggle()
	enabled = not enabled
	log("Plugin " .. (enabled and "enabled" or "disabled") .. " via API")
	if enabled then
		switch_to_default()
	end
	return enabled
end

---Check if the plugin is enabled
---@return boolean enabled True if enabled, false otherwise
function M.is_enabled()
	return enabled
end

return M
