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

``` markdown
# Simple File

\`\`\`bash
echo hello world
\`\`\`
```

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

``` markdown
# Simple File

\`\`\`bash
echo hello world
\`\`\`

<!-- notebook output start -->
<!-- notebook output modified -->

>        hello world

<!-- notebook output end -->
```

### Expanded

`cat <file> | notebook-md execute | notebook-md expand`

``` markdown
# Simple File

\`\`\`bash
bash << EOF
echo hello world
EOF
\`\`\`

<!-- notebook output start -->
<!-- notebook output modified -->

>        hello world

<!-- notebook output end -->
```

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
