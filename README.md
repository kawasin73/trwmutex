# Try Read Write Mutex for Go

Extended sync.RWMutex which have `TryLock()` and `TryRLock()`.

Use `sync.Mutex` and `sync.RWMutex` internally, not channel.

Of course `TryLock()` `TryRLock()` never blocks so long.

`trwmutex.TRWMutex` can be used as `sync.RWMutex` with same interface. 

## How to use

```bash
go get github.com/kawasin73/trwmutex
```

## API

- `Lock()`
- `TryLock() -> bool`
- `Unlock()`
- `RLock()`
- `TryRLock() -> bool`
- `RUnlock()`

## LICENSE

MIT
