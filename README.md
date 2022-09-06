# pixiv-distributed-crawler

#### UID

$$
\rm \xrightarrow{given} uid \xrightarrow{following (cookie?)} uids(dalao)
$$

$$
\rm \xrightarrow{given} uid(dalao)
$$
#### Series

$$
\rm uid \xrightarrow{work} series
$$

$$
\rm uid \xrightarrow{bookmark (cookie?)} series
$$

$$
\rm pid \xrightarrow{} series
$$
#### PID

$$
\rm \xrightarrow{ranking} pids
$$

$$
\rm series \xrightarrow{include} pids
$$

$$
\rm uid \xrightarrow{work} pids
$$

$$
\rm uid \xrightarrow{bookmark (cookie?)} pids
$$

$$
\rm tag \xrightarrow{(thinking)} tags \xRightarrow{(filtering)} tags \xRightarrow{search (cookie)} pids
$$

#### Sync

```go
func sync() {
    old := read()
    if !old.exists {
        save()
        return
    }
    if info.timestamp != old.timestame {
        update_and_archive()
        return
    }
    return
}
```
