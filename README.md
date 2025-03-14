# Notebook.md

Markdown code block execution.

## Usage

### Execution

To execute all code blocks in a file:

`cat <file> | notebook-md execute`

### Expansion

To expand code blocks using a different style for use in publishing:

`cat <file> | notebook-md expand`

## Example

### File

```` markdown
# Simple File

```bash
echo hello world
```
````

### Config

``` yaml
notebook:
  commands:
  - name: bash
    command: bash
    expand:
      style: heredoc
```

### Executed

`cat <file> | notebook-md execute`

```` markdown
# Simple File

```bash
echo hello world
```

<!-- notebook output start -->
<!-- notebook output modified -->

>        hello world

<!-- notebook output end -->
````

### Expanded

`cat <file> | notebook-md execute | notebook-md expand`

```` markdown
# Simple File

```bash
bash << EOF
echo hello world
EOF
```

<!-- notebook output start -->
<!-- notebook output modified -->

>        hello world

<!-- notebook output end -->
````

## Config

Commands must pass `--config <config>` or set
`NOTEBOOK_MD_CONFIG=<config>`

``` yaml
notebook:
  commands:
  - name: name-next-to-backticks
    command: to be execute
    display_style: RAW|QUOTE
    expand:
      block_name: next-to-backtick-override
      command_name: inline name override (style dependent)
      style: NONE|HIDE|LINE|ONCE|HEREDOC
```

## Vim

``` vim
function! RunNotebook(mode) range
    write
    echo "Executing..."
    if a:mode == 0
        " Execute the entire file
        execute "%!notebook-md execute"
    endif
    if a:mode == 1
        " Execute the selected lines
        execute "%!notebook-md execute --line " . a:firstline . "-" . a:lastline
    endif
    if a:mode == 2
        " Execute the current line and the rest of the file
        execute "%!notebook-md execute --line " . line('.') . "-"
    endif
    silent !tput bel
    redraw!
    redraw!
    echo "Done"
endfunction

autocmd Filetype markdown map <leader>E :call RunNotebook(0)<CR>
autocmd Filetype markdown map <leader>e :call RunNotebook(1)<CR>
autocmd Filetype markdown map <leader>ee :call RunNotebook(2)<CR>
```
