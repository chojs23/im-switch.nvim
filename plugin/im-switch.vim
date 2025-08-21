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

" Initialize plugin - LazyVim will handle the build
lua require('im-switch').setup()
