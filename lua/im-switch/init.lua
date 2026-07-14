---@type ImSwitch
local M = {}

---@return boolean
local function is_wsl()
	if vim.fn.has("wsl") == 1 then
		return true
	end

	return vim.env.WSL_DISTRO_NAME ~= nil or vim.env.WSL_INTEROP ~= nil
end

---@return ImSwitchConfig
local function get_default_config()
	local is_mac = vim.fn.has("mac") == 1 or vim.fn.has("macunix") == 1
	local wsl = is_wsl()
	local is_linux = vim.fn.has("unix") == 1 and not is_mac

	if is_mac then
		return {
			binary_path = "im-switch",
			default_input = "com.apple.keylayout.ABC",
			auto_switch = true,
			auto_capslock_off = true,
			debug = false,
		}
	elseif wsl then
		return {
			binary_path = "im-switch",
			default_input = "en-US",
			auto_switch = true,
			auto_capslock_off = true,
			debug = false,
		}
	elseif is_linux then
		return {
			binary_path = "im-switch",
			default_input = "us",
			auto_switch = true,
			auto_capslock_off = true,
			debug = false,
		}
	else
		return {
			binary_path = "im-switch",
			default_input = "en-US",
			auto_switch = true,
			auto_capslock_off = true,
			debug = false,
		}
	end
end

---@type ImSwitchConfig
local config = get_default_config()

---@type boolean
local enabled = true

---@param msg string
local function log(msg)
	if config.debug then
		print("[im-switch] " .. msg)
	end
end

---@param args string|nil Command arguments
---@return string[] cmd Argv list for vim.system
local function build_cmd(args)
	local cmd = { config.binary_path }
	if args then
		table.insert(cmd, args)
	end
	return cmd
end

---@param result vim.SystemCompleted
---@return string|nil output Trimmed stdout or nil on failure
local function command_output(result)
	if result.code ~= 0 then
		return nil
	end
	return (result.stdout or ""):gsub("%s+$", "")
end

---Run the binary and block until it exits. Only for API functions that
---must return a value; editor event handlers use execute_command_async.
---@param args string|nil Command arguments
---@return string|nil result Command output or nil if failed
local function execute_command(args)
	local cmd = build_cmd(args)

	local ok, proc = pcall(vim.system, cmd, { text = true })
	if not ok then
		log("Failed to execute: " .. table.concat(cmd, " "))
		return nil
	end

	local output = command_output(proc:wait())
	if not output then
		log("Command failed: " .. table.concat(cmd, " "))
	end
	return output
end

---Run the binary without blocking the editor. The callback (if any) runs
---on the main loop via vim.schedule, so it may safely use the vim API.
---@param args string|nil Command arguments
---@param callback? fun(output: string|nil)
local function execute_command_async(args, callback)
	local cmd = build_cmd(args)

	local ok = pcall(vim.system, cmd, { text = true }, function(result)
		vim.schedule(function()
			local output = command_output(result)
			if not output then
				log("Command failed: " .. table.concat(cmd, " "))
			end
			if callback then
				callback(output)
			end
		end)
	end)

	if not ok then
		log("Failed to execute: " .. table.concat(cmd, " "))
		if callback then
			callback(nil)
		end
	end
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

local function turn_off_capslock()
	if config.auto_capslock_off then
		execute_command_async("--capslock-off")
	end
end

---Switch to the default input method without blocking. This runs on every
---insert-mode exit, so the query and the switch are both async; a blocking
---call here would freeze the editor for the binary's runtime on each Esc.
local function switch_to_default()
	if not (enabled and config.auto_switch) then
		return
	end

	turn_off_capslock()

	execute_command_async(nil, function(current)
		-- The mode may have changed while the query was in flight. Never
		-- force-switch under the user's fingers if they are typing again.
		local mode = vim.fn.mode()
		if not enabled or mode:find("^i") or mode:find("^R") then
			return
		end

		if current and current ~= config.default_input then
			execute_command_async(config.default_input, function(output)
				if output then
					log("Switched to: " .. config.default_input)
				end
			end)
		end
	end)
end

local function on_mode_changed()
	if not enabled then
		return
	end

	local mode = vim.fn.mode()

	if mode == "n" or mode == "c" then
		switch_to_default()
	end
end

local function on_focus_gained()
	local mode = vim.fn.mode()

	if enabled and mode ~= "i" then
		switch_to_default()
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

	vim.api.nvim_create_autocmd({ "VimEnter" }, {
		group = group,
		callback = switch_to_default,
	})
end

local function create_user_commands()
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
		local wsl_candidates = {
			"/mnt/c/Tools/im-switch.exe",
			"/mnt/c/Windows/im-switch.exe",
			"/mnt/c/Program Files/im-switch/im-switch.exe",
		}

		if vim.fn.executable(local_binary) == 1 then
			config.binary_path = local_binary
		elseif vim.fn.executable(system_binary) == 1 then
			config.binary_path = system_binary
		elseif is_wsl() then
			for _, candidate in ipairs(wsl_candidates) do
				if vim.fn.executable(candidate) == 1 then
					config.binary_path = candidate
					break
				end
			end
			if not config.binary_path or config.binary_path == "" then
				config.binary_path = "im-switch"
			end
		else
			config.binary_path = "im-switch" -- fallback to PATH
		end
	end

	if vim.fn.executable(config.binary_path) ~= 1 then
		vim.notify("[im-switch] Binary not found: " .. config.binary_path, vim.log.levels.WARN)
		return
	end

	-- Probe the binary asynchronously so setup never blocks startup.
	-- Autocmds and commands are only registered once the probe succeeds,
	-- matching the old behavior of bailing out on a broken binary.
	execute_command_async(nil, function(current)
		if not current or current == "" then
			vim.notify("[im-switch] Failed to get current input method", vim.log.levels.WARN)
			return
		end

		log("Current input method: " .. current)
		log("Default input method: " .. config.default_input)

		setup_autocmds()
		create_user_commands()

		log("im-switch plugin initialized")
	end)
end

function M.switch_to_english()
	switch_to_default()
end

function M.turn_off_capslock()
	turn_off_capslock()
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
