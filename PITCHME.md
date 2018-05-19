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
func main() {
  for i := 0; i < 10; i++ {
    go func() {
      fmt.Println("Hello Concurrency")
    }()
  }
}
```

---

This code is NOOP :-(

---

### True Simple Concurrency

---

```
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            fmt.Println("Hello")
            wg.Done()
        }()
    }
    wg.Wait()
}
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

### Find mistake

---

```
func main() {
    type value struct {
        mu    sync.Mutex
        value int


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
}
```

---
