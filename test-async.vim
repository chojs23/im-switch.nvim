" Test script for async build
set rtp+=/Users/neo/Desktop/personal/im-switch
runtime plugin/im-switch.vim

" Wait a bit and check if plugin loads correctly
autocmd User AsyncBuildComplete echo "Build completed, testing plugin..."

" Add a timer to check binary after a few seconds
function! CheckBinaryExists()
  if executable('/Users/neo/Desktop/personal/im-switch/im-switch')
    echo "✓ Binary exists and is executable"
    lua print("Current input method:", require('im-switch').get_current())
  else
    echo "✗ Binary not found or not executable"
  endif
endfunction

call timer_start(3000, {-> CheckBinaryExists()})
