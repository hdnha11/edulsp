# Educational LSP

A Language Server built for the educational purpose of understanding **WHAT** LSP is and **HOW** it works.

It doesn't do anything special for any particular language, it is focussed on helping you understand what your tools **do**.

This LSP was tested with Neovim, but would likely work inside VS Code.

Video: https://www.youtube.com/watch?v=YsdlcQoHqPY&t

## Example Neovim config

```lua
local client = vim.lsp.start_client {
  name = "edulsp",
  cmd = { "/Users/nhahoang/workspaces/learning/edulsp/edulsp" },
  on_attach = require('user.lsp.handlers').on_attach,
}

if not client then
  vim.notify "hey, you didn't do the client thing good"
  return
end

vim.api.nvim_create_autocmd("FileType", {
  pattern = "markdown",
  callback = function ()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```
