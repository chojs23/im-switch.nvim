local M = {}

local function mock_vim()
	_G.im_switch_test_commands = {}

	local vim_mock = _G.vim or {}
	vim_mock.env = vim_mock.env or {}
	vim_mock.fn = vim_mock.fn or {}
	vim_mock.fn.has = function(feature)
		if feature == "mac" or feature == "macunix" then
			return 1
		end
		return 0
	end
	vim_mock.fn.executable = function(path)
		return 1
	end
	vim_mock.fn.fnamemodify = function(path, modifier)
		return "/test/path"
	end
	vim_mock.tbl_deep_extend = function(behavior, ...)
		local result = {}
		for _, t in ipairs({ ... }) do
			for k, v in pairs(t) do
				result[k] = v
			end
		end
		return result
	end
	vim_mock.api = vim_mock.api or {}
	vim_mock.api.nvim_create_augroup = function(name, opts)
		return 1
	end
	vim_mock.api.nvim_create_autocmd = function(events, opts) end
	vim_mock.api.nvim_create_user_command = function(name, callback, opts) end
	vim_mock.notify = function(msg, level) end
	vim_mock.log = { levels = { WARN = 1 } }
	vim_mock.split = function(str, delimiter)
		local result = {}
		for match in (str .. delimiter):gmatch("(.-)" .. delimiter) do
			table.insert(result, match)
		end
		return result
	end
	_G.vim = vim_mock

	local debug_mock = _G.debug or {}
	debug_mock.getinfo = function(level, what)
		return { source = "@/test/path/init.lua" }
	end
	_G.debug = debug_mock

	_G.io = {
		popen = function(cmd)
			table.insert(_G.im_switch_test_commands, cmd)

			local handle = {
				read = function(self, format)
					if cmd:find("im-switch$") then
						return "com.apple.keylayout.ABC"
					elseif cmd:find("im-switch -l") then
						return "com.apple.keylayout.ABC\ncom.apple.inputmethod.Korean.2SetKorean"
					elseif cmd:find("im%-switch %-%-capslock%-off") then
						return ""
					else
						return "com.apple.keylayout.ABC"
					end
				end,
				close = function(self)
					return true
				end,
			}
			return handle
		end,
	}
end

mock_vim()
local im_switch = require("im-switch.init")

local function assert_eq(actual, expected, msg)
	if actual ~= expected then
		error(string.format("%s: expected %s, got %s", msg or "assertion failed", tostring(expected), tostring(actual)))
	end
end

local function assert_true(value, msg)
	if not value then
		error(msg or "expected true")
	end
end

local function assert_not_nil(value, msg)
	if value == nil then
		error(msg or "expected non-nil value")
	end
end

function M.test_default_config()
	local config = {
		binary_path = "im-switch",
		default_input = "com.apple.keylayout.ABC",
		auto_switch = true,
		debug = false,
		auto_capslock_off = true,
	}

	assert_eq(config.binary_path, "im-switch", "binary_path should be im-switch")
	assert_eq(config.default_input, "com.apple.keylayout.ABC", "default_input should be ABC")
	assert_true(config.auto_switch, "auto_switch should be true")
	assert_true(config.auto_capslock_off, "auto_capslock_off should be true")
end

function M.test_get_current_input()
	local current = im_switch.get_current()
	assert_not_nil(current, "get_current should return a value")
end

function M.test_list_inputs()
	local inputs = im_switch.list_inputs()
	assert_not_nil(inputs, "list_inputs should return a table")
	assert_true(type(inputs) == "table", "list_inputs should return a table")
end

function M.test_switch_to_english()
	_G.im_switch_test_commands = {}
	im_switch.switch_to_english()

	local found = false
	for _, cmd in ipairs(_G.im_switch_test_commands) do
		if cmd:find("im%-switch %-%-capslock%-off") then
			found = true
			break
		end
	end

	assert_true(found, "switch_to_english should turn Caps Lock off")
end

function M.test_set_input()
	local result = im_switch.set_input("com.apple.keylayout.ABC")
	assert_true(result, "set_input should return true for valid input")
end

function M.test_turn_off_capslock()
	_G.im_switch_test_commands = {}
	im_switch.turn_off_capslock()

	assert_true(
		_G.im_switch_test_commands[1]:find("im%-switch %-%-capslock%-off") ~= nil,
		"turn_off_capslock should call the capslock command"
	)
end

function M.test_auto_capslock_off_can_be_disabled()
	im_switch.setup({ auto_capslock_off = false })
	_G.im_switch_test_commands = {}
	im_switch.switch_to_english()

	for _, cmd in ipairs(_G.im_switch_test_commands) do
		assert_true(
			cmd:find("im%-switch %-%-capslock%-off") == nil,
			"switch_to_english should not call capslock command when auto_capslock_off is false"
		)
	end

	im_switch.setup({ auto_capslock_off = true })
end

function M.run_all_tests()
	local tests = {
		"test_default_config",
		"test_get_current_input",
		"test_list_inputs",
		"test_switch_to_english",
		"test_set_input",
		"test_turn_off_capslock",
		"test_auto_capslock_off_can_be_disabled",
	}

	local passed = 0
	local failed = 0

	for _, test_name in ipairs(tests) do
		local success, err = pcall(M[test_name])
		if success then
			print("✓ " .. test_name)
			passed = passed + 1
		else
			print("✗ " .. test_name .. ": " .. err)
			failed = failed + 1
		end
	end

	print(string.format("\nTests: %d passed, %d failed", passed, failed))
	return failed == 0
end

if not package.loaded["busted"] then
	M.run_all_tests()
end

return M
