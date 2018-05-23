### High-performance concurrency

syossan27

---

### Concurrency is ... ?

---

#### Rob Pike says

![RobPike](assets/images/robpike.jpg)

---

***But when people hear the word concurrency they often think of parallelism, a related but quite distinct concept.***   
***In programming, concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of (possibly related) computations.***   
***Concurrency is about dealing with lots of things at once.***   
***Parallelism is about doing lots of things at once.***   

---

***Concurrencyという単語を聞くと、しばしばParallelismのことを考えるかもしれないが、関連こそすれ全く別の概念である。***   
***プログラミングでは、Concurrencyは「独立して実行されるプロセスによる構成」のことですが、Parallelismは「（おそらく関連する）処理の同時実行」のことです。***   
***つまりConcurrencyは、「一度に多くのことを扱うこと」***   
***Parallelismは「一度にたくさんのことをすること」 ***

---?image=assets/images/concurrency_parallelism1.png&size=auto 70%

---?image=assets/images/concurrency_parallelism2.png&size=auto 70%

---

### Simple Concurrency

---

```
for i := 0; i < 10; i++ {
  go func() {
    fmt.Println("Hello")
  }()
}
```

---

This code is NOOP😩

---

### True Simple Concurrency

---

```
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        fmt.Println("Hello")
        wg.Done()
    }()
}
wg.Wait()
```

---

```
Hello
Hello
Hello
Hello
Hello
Hello
Hello
Hello
Hello
Hello
```

---

### It's Concurrency👌

---

### High-performance Concurrency is ... ?

---

### Find mistake👀

---

```
type value struct {
    mu    sync.Mutex
    value int
}

var wg sync.WaitGroup
printSum := func(v1, v2 *value) {
    defer wg.Done()

    v1.mu.Lock()
    defer v1.mu.Unlock()

    time.Sleep(2 * time.Second)

    v2.mu.Lock()
    defer v2.mu.Unlock()

    fmt.Printf("sum=%v\n", v1.value+v2.value)
}

var a, b value
wg.Add(2)
go printSum(&a, &b)
go printSum(&b, &a)
wg.Wait()
```

---

```
fatal error: all goroutines are asleep - deadlock!
```

---

Why?🤔

---

### Coffman Conditions

- Mutual Exclusion: リソースは最大１つまでのプロセスにしか確保されないこと
- Wait For Condition:  リソースが確保済みの場合、要求している他のプロセスは待たなければならない
- No Preemption: リソースは確保したプロセスによってのみ解放される
- Circular Wait: リソースを確保しているプロセスAが、他のリソースを確保しているプロセスBのリソースを要求することにより循環待ちが発生する

これを全て満たすとDeadLockを引き起こす

---

### Other

- LiveLock
- Starvation
- MemoryLeak
- etc...

Concurrencyでは考慮しなければならないことが多い👿

---

### High-performance Concurrency = Safety💪 

---

### Basic Concurrency Pattern

- Confinement
- Preventing Goroutine Leaks
- Timeouts and Cancellation

---

### Confinement

---

### Confinement

goroutine内で使うデータに制限をかける手法   
開発者の認知負荷の軽減や、小さなクリティカルセクションに対して効果がある

---

### Type

- ad hoc
- lexical

---

### adhoc

チーム内での認識合意のみで実現する

---

```
data := []int{1, 2, 3, 4}

loopData := func(handleData chan<- int) {
    defer close(handleData)
    for i := range data {
        handleData <- data[i]
    }
}

handleData := make(chan int)
go loopData(handleData)

for num := range handleData {
    fmt.Println(num)
}
```
@[1](mainからもloopDataからも参照出来てしまう)

---

### 🙅

認識のみで縛る方法なので、非常に危険

---

### lexical

レキシカルスコープを利用して、変数へのアクセスを制限する

---

```
loopData := func(handleData chan<- int) {
    defer close(handleData)
    data := []int{1, 2, 3, 4}
    for i := range data {
        handleData <- data[i]
    }
}

handleData := make(chan int)
go loopData(handleData)
for num := range handleData {
    fmt.Println(num)
}
```
@[3](データの参照範囲を明確にする)

---

### 🙆

クリティカルセクションも無くなり、認知負荷が軽減

---

### Preventing Goroutine Leaks

goroutineがGCで解放されないパターンに対応する

---

### Paths to termination

1. 処理の終了
1. 回復不能なエラーの発生
1. 処理の停止

1, 2はGCが動くが、3は動かない

---

### Leak Pattern

---

```
doWork := func(strings <-chan string) <-chan interface{} {
    completed := make(chan interface{})
    go func() {
        defer close(completed)
        for s := range strings { fmt.Println(s) }
    }()
    return completed
}
doWork(nil)
time.Sleep(5 * time.Second)
fmt.Println("Done.")
```
@[12](nil channelは読み込みも書き込みもブロッキングされる)

---

### If long lifecycle application...

😥

---

### Parent goroutine manage child goroutine

---

```
doWork := func(
  done <-chan interface{},
  strings <-chan string,
) <-chan interface{} {
    completed := make(chan interface{})
    go func() {
        defer close(completed)
        for {
            select {
            case s := <-strings: fmt.Println(s)
            case <-done: return
            }
        }
    }()
    return completed
}

done := make(chan interface{})
completed := doWork(done, nil)

go func() {
    time.Sleep(1 * time.Second)
    fmt.Println("Canceling doWork goroutine...")
    close(done)
}()

<-completed
fmt.Println("Done.")
```
@[2](処理の終了を知らせるchannel)
@[21-25](mainから子goroutineへ処理の終了を伝える)
@[11](処理が正常に終了)

---

### Digression

GoのGCではgoroutineでヒープ領域に確保したメモリをOSに返さず、新しく生成されるgoroutineのために再利用しようとする性質があります。   
そのため、 `pprof` や `runtime.MemStats` で確認した時に単純にメモリ使用量が減らないからメモリリークしている、という勘違いをしないよう注意しましょう。   
真に確認するには `pprof` のスタックダンプや、 `leaktest` などのツールを使ったりするのが良いでしょう。

---

### Timeouts and Cancellation

---

### Merit by use Timeouts

- リトライ

---
