# Core

## Working in a Monorepo

When you need to substitute a dependency with a local version try `replace`.
```
go mod edit -replace "github.com/devpies/core/tenant@v0.0.0 = ../tenant"
```
You may need to `require` the package before `replace` if it's not being used. 
```
go mod edit -require github.com/devpies/core/tenant@v0.0.0
```
Note: if there's no version tag, it defaults to `v0.0.0`.

## Tagging modules

When you're happy with your changes you need to tag the last commit.
```bash
git commit -m "feat(tenant): add new method"
git push origin feature_branch
git tag tenant/v0.1.0
git push --tags
```
Once the tag is pushed you can use it.
```bash
go get github.com/devpies/core/tenant@v0.1.0
```

### Removing tags
```
git tag -d tenant/v0.2.0
git push --delete origin tenant/v0.2.0
```
