# gh-foreach

Automatically clone and execute commands across multiple GitHub repositories.

## Usage

Note: the minimum required version of the GitHub CLI (`gh`) is `v2.3.0`.

```bash
gh extension install mtoohey31/gh-foreach
gh foreach --help
gh foreach run -a owner -v public -l go -s bash cat go.mod
gh foreach run -a owner -l rust -- bat --pager=never Cargo.toml
gh foreach run --no-confirm -v public -r '^\w{3}\d{3}.*project$' 'cat README.md | head -n 1'
gh foreach run -ci -l svelte 'pnpm install; fish'
```

![demo](https://user-images.githubusercontent.com/36740602/152471108-47cc7484-3f95-4da0-81e5-81b8fb35d1ea.gif)

## Completions

Completions are supported via [kongplete](https://github.com/WillAbides/kongplete), and can be generated with `gh foreach install-completions`. However, they likely won't work as expected when `gh-foreach` is run as a subcommand of `gh`, since your shell will use `gh`'s completions. To work around this, you can add the `gh-foreach` binary to your path and run it directly.
