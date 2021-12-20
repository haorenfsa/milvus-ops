# Get started

prerequisite:
- go:1.13
- yarn

or install prerequisites on mac with homebrew:

```
brew update
brew install golang
brew install yarn
```



build frontend:

```
cd web && yarn install && yarn build
```


run backend:

```
go run ./cmd/milvus-ops
```

you can now find the page on http://localhost:8080/app/
