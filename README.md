# shortcuts

A command-line tool meant to be used from vim, but equally useful in any shell piping situation.

Like many of my other projects, this is only intended to be suitable for use by the author, and is only on GitHub for my convenience. You are welcome to use it, or change it to make it better for your own needs, but it will probably never cater to a general audience.

# usage

## shell

```bash
shortcuts < old.txt > new.wxt
```

## vim 

From normal mode:

`:%! shortcuts`

## vim (automatic)

This is not recommended unless you know what you're doing. It _will_ overwrite your buffer in insert mode while you're in the midst of editing. You can lose work.

The vimscript function exists solely to restore the cursor to its previous position; by default the cursor is returned to the top of the file after the buffer is completely rewritten, which tends to interrupt the flow of coding. Otherwise, all this does is automate the call to `:%! shortcuts` shown above.

```vim
set updatetime=1000 
" pause in activity in insert mode
autocmd CursorHoldI * :call Shortcuts()
function Shortcuts()
    let pos = getcurpos()
    let offset = pos[2] 
    let args = {"file": expand("%:p"), "line": getline("."), "line_num": pos[1], "col_num": pos[2]}
    execute "%!shortcuts"
    execute args["line_num"] + "G"
    " execute ":normal! " . offset . "|"
    :call cursor(pos[1], pos[2]) 
    " execute ":normal! 0" . offset . "l"
endfunction
```

BSD licenced.
