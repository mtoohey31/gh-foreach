# gh-foreach

Execute commands across multiple GitHub repositories. Currently a WIP.

## Usage

```bash
gh extension install mtoohey31/gh-foreach
gh foreach --help
gh foreach completion --help # see shell specific help pages for installation instructions
gh foreach run -a owner -v public -l go -s bash cat go.mod
gh foreach run -a owner -l rust -- bat --pager=never Cargo.toml
```
