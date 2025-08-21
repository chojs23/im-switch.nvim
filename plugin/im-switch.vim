if exists('g:loaded_im_switch')
  finish
endif
let g:loaded_im_switch = 1

if !has('nvim')
  echohl WarningMsg
  echo 'im-switch requires Neovim'
  echohl None
  finish
endif

" Auto-build binary if it doesn't exist
let s:plugin_dir = expand('<sfile>:p:h:h')
let s:binary_path = s:plugin_dir . '/im-switch'

if !executable(s:binary_path)
  echohl WarningMsg
  echo 'Building im-switch binary...'
  echohl None
  
  let s:build_result = system('cd ' . shellescape(s:plugin_dir) . ' && make build')
  
  if v:shell_error != 0
    echohl ErrorMsg
    echo 'Failed to build im-switch binary: ' . s:build_result
    echohl None
    finish
  endif
  
  if !executable(s:binary_path)
    echohl ErrorMsg
    echo 'Binary not found after build: ' . s:binary_path
    echohl None
    finish
  endif
  
  echohl MoreMsg
  echo 'im-switch binary built successfully!'
  echohl None
endif

lua require('im-switch').setup()
