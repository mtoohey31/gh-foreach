# gh-foreach

Automatically clone and execute commands across multiple GitHub repositories.

## Usage

Note: the minimum required version of the GitHub CLI (`gh`) is `v2.3.0`.

```bash
gh extension install mtoohey31/gh-foreach
gh foreach --help
gh foreach completion --help # see shell specific help pages for installation instructions
gh foreach run -a owner -v public -l go -s bash cat go.mod
gh foreach run -a owner -l rust -- bat --pager=never Cargo.toml
gh foreach run --no-confirm -v public -r '^\w{3}\d{3}.*project$' 'cat README.md | head -n 1'
gh foreach run -ci -l svelte 'pnpm install; fish'
```

![demo](https://user-images.githubusercontent.com/36740602/152471108-47cc7484-3f95-4da0-81e5-81b8fb35d1ea.gif)
