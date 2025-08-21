#!/usr/bin/env lua

package.path = package.path .. ";./lua/?.lua;./lua/?/init.lua"

local function run_lua_tests()
	print("Running Lua tests for im-switch plugin...")
	local test_module = require("im-switch.init_test")
	local success = test_module.run_all_tests()
	return success
end

local function run_integration_test()
	print("\nRunning integration test...")

	local im_switch = require("im-switch.init")

	_G.vim = {
		fn = {
			has = function(feature)
				if feature == "mac" or feature == "macunix" then
					return 1
				elseif feature == "unix" then
					return 1
				end
				return 0
			end,
			executable = function(path)
				return 1
			end,
			fnamemodify = function(path, modifier)
				return "./im-switch"
			end,
		},
		tbl_deep_extend = function(behavior, ...)
			local result = {}
			for _, t in ipairs({ ... }) do
				for k, v in pairs(t) do
					result[k] = v
				end
			end
			return result
		end,
		api = {
			nvim_create_augroup = function(name, opts)
				return 1
			end,
			nvim_create_autocmd = function(events, opts) end,
		},
		notify = function(msg, level)
			print("NOTIFY: " .. msg)
		end,
		log = { levels = { WARN = 1 } },
		split = function(str, delimiter)
			local result = {}
			for match in (str .. delimiter):gmatch("(.-)" .. delimiter) do
				table.insert(result, match)
			end
			return result
		end,
	}

	_G.debug = {
		getinfo = function(level, what)
			return { source = "@" .. os.getenv("PWD") .. "/lua/im-switch/init.lua" }
		end,
	}

	_G.io = {
		popen = function(cmd)
			local handle = io.popen(cmd .. " 2>&1")
			if not handle then
				return {
					read = function()
						return ""
					end,
					close = function()
						return false
					end,
				}
			end
			return handle
		end,
	}

	local success, err = pcall(function()
		im_switch.setup({ debug = true })

		local current = im_switch.get_current()
		if current and current ~= "" then
			print("✓ get_current: " .. current)
		else
			print("✗ get_current returned empty or nil")
		end

		local inputs = im_switch.list_inputs()
		if inputs and #inputs > 0 then
			print("✓ list_inputs: found " .. #inputs .. " inputs")
			for i, input in ipairs(inputs) do
				if i <= 3 then
					print("  - " .. input)
				end
			end
			if #inputs > 3 then
				print("  ... and " .. (#inputs - 3) .. " more")
			end
		else
			print("✗ list_inputs returned empty or nil")
		end
	end)

	if not success then
		print("✗ Integration test failed: " .. err)
		return false
	end

	print("✓ Integration test completed")
	return true
end

local function main()
	local lua_success = run_lua_tests()
	local integration_success = run_integration_test()

	print("\n" .. string.rep("=", 50))
	if lua_success and integration_success then
		print("All tests passed!")
		os.exit(0)
	else
		print("Some tests failed!")
		os.exit(1)
	end
end

main()

