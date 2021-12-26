# gh-foreach

Execute commands across multiple GitHub repositories. Currently a WIP.

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

![demo](https://user-images.githubusercontent.com/36740602/145335345-03e00a82-168f-482e-9616-0a38c1c8649b.gif)
