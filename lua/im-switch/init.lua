local M = {}

local config = {
	binary_path = "im-switch",
	english_input = "com.apple.keylayout.ABC",
	auto_switch = true,
	auto_restore = true,
	debug = false,
}

local saved_input = nil
local last_mode = nil

local function log(msg)
	if config.debug then
		print("[im-switch] " .. msg)
	end
end

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

	return result:gsub("%s+$", "") -- trim whitespace
end

-- Get current input method
local function get_current_input()
	return execute_command()
end

-- Set input method
local function set_input(input_id)
	if input_id and input_id ~= "" then
		local result = execute_command(input_id)
		log("Switched to: " .. input_id)
		return result ~= nil
	end
	return false
end

-- Switch to English input
local function switch_to_english()
	if config.auto_switch then
		local current = get_current_input()
		if current and current ~= config.english_input then
			saved_input = current
			log("Saving current input: " .. current)
			set_input(config.english_input)
		end
	end
end

-- Restore previous input
local function restore_input()
	if config.auto_restore and saved_input then
		log("Restoring input: " .. saved_input)
		set_input(saved_input)
		saved_input = nil
	end
end

-- Handle mode changes
local function on_mode_changed()
	local mode = vim.fn.mode()

	if mode == "n" or mode == "c" then
		-- Normal or command mode - switch to English
		switch_to_english()
	elseif mode == "i" and last_mode == "n" then
		-- Entering insert mode from normal mode - restore previous input
		restore_input()
	end

	last_mode = mode
end

-- Handle focus events
local function on_focus_gained()
	switch_to_english()
end

local function on_focus_lost()
	-- Save current state when losing focus
	local current = get_current_input()
	if current and current ~= config.english_input then
		saved_input = current
	end
end

-- Setup autocommands
local function setup_autocmds()
	local group = vim.api.nvim_create_augroup("ImSwitch", { clear = true })

	-- Mode change events
	vim.api.nvim_create_autocmd({ "ModeChanged" }, {
		group = group,
		callback = on_mode_changed,
	})

	-- Focus events
	vim.api.nvim_create_autocmd({ "FocusGained" }, {
		group = group,
		callback = on_focus_gained,
	})

	vim.api.nvim_create_autocmd({ "FocusLost" }, {
		group = group,
		callback = on_focus_lost,
	})

	-- Also switch to English when entering Neovim
	vim.api.nvim_create_autocmd({ "VimEnter" }, {
		group = group,
		callback = switch_to_english,
	})

	-- Restore input when leaving Neovim
	vim.api.nvim_create_autocmd({ "VimLeavePre" }, {
		group = group,
		callback = restore_input,
	})
end

-- Setup function
function M.setup(opts)
	opts = opts or {}
	config = vim.tbl_deep_extend("force", config, opts)

	-- If no binary_path specified, use the one from plugin directory
	if not opts or not opts.binary_path then
		local plugin_dir = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h:h:h")
		config.binary_path = plugin_dir .. "/im-switch"
	end

	-- Validate binary exists
	if vim.fn.executable(config.binary_path) ~= 1 then
		vim.notify("[im-switch] Binary not found: " .. config.binary_path, vim.log.levels.WARN)
		return
	end

	-- Test if we can get current input
	local current = get_current_input()
	if not current or current == "" then
		vim.notify("[im-switch] Failed to get current input method", vim.log.levels.WARN)
		return
	end

	log("Current input method: " .. current)
	log("English input method: " .. config.english_input)

	setup_autocmds()

	log("im-switch plugin initialized")
end

-- Expose functions for manual use
function M.switch_to_english()
	switch_to_english()
end

function M.restore_input()
	restore_input()
end

function M.get_current()
	return get_current_input()
end

function M.set_input(input_id)
	return set_input(input_id)
end

function M.list_inputs()
	local result = execute_command("-l")
	if result then
		return vim.split(result, "\n")
	end
	return {}
end

return M
