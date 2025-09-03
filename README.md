# util

A small collection of generic utilities for Go, grouped under the module gitlab.thepratts.info/util.

These helpers cover common needs such as averages, debouncing, email relay, formatting, JSON helpers, math comparisons, pub/sub, random queues, strings, time series, value tracking, vector math, and worker pools.

- Go version: 1.24+
- License: MIT (see LICENSE)

## Installation

- With modules:

  go get gitlab.thepratts.info/util

- Import:

  import "gitlab.thepratts.info/util"

## Contents

Below is an overview of each package utility with short examples.

### arrayUtils.go
- Count[T any](d []T, f func(T) bool) int64
  Counts elements in a slice that satisfy a predicate. Internally parallelizes over NumCPU blocks.

Example:

  nums := []int{1,2,3,4,5}
  evens := util.Count(nums, func(v int) bool { return v%2 == 0 })
  // evens == 2

Notes:
- Safe for large slices. Remainder processed on main goroutine.

### avg.go
- Average
  Running arithmetic mean of all observed values.
- MovingAverage
  Smoothed moving average (simple stable smoothing using factor = 1/n).

Example Average:

  a := util.NewAverage(10)
  a.Update(20)
  a.Update(30)
  mean := a.Get() // (10+20+30)/3 = 20

Example MovingAverage:

  ma := util.NewMovingAverage(10, 0)
  _ = ma.UpdateAndGet(5)
  cur := ma.Get()

Notes:
- Not concurrency-safe.
- NewMovingAverage panics if n == 0.

### debounce.go
- Debouncer[T any]
  Caches values for a given delay using a fetcher. Simple, single-value cache.
- PubSubDebouncer[T comparable]
  Debounced value with pub/sub notifications and background fetch loop.

Example Debouncer:

  d := util.NewDebouncer(time.Second, func() (int, error) { return expensiveCompute(), nil })
  v, _ := d.GetValue() // fetch; subsequent calls within delay use cached value

Example PubSubDebouncer:

  db, cancel := util.NewPubSubDebouncer(time.Second, fetchValue)
  defer cancel()
  ch := db.Register()
  go func() {
      for v := range ch {
          // handle updates
          _ = v
      }
  }()
  _ = db.GetValueMust()
  db.Unregister(ch)

Notes:
- PubSubDebouncer requires delay >= 100ms (panics otherwise).
- GetValueMust panics on fetch error.

### email.go
- EmailRelay
  Minimal SMTP relay sender using net/smtp.

Example:

  relay := util.NewEmailRelay("smtp.example.com")
  err := relay.Send("from@example.com", "to@example.com", "Subject: Hi\r\n\r\nHello!")

Notes:
- Uses port 587; provide fully formatted message body with headers.

### formatter.go
Preconfigured github.com/leekchan/accounting formatters:
- DollarsFormat ($, 2 decimals)
- CentsFormat (¢, 2 decimals)
- DollarsCentsFormat ($, 4 decimals)
- DollarsCentsFormat2 ($, 2 decimals)
- PercentFormat (%, 2 decimals)
- NumberFormat (thousands separators, 0 decimals)

Example:

  price := util.DollarsFormat.FormatMoney(123.456) // "$123.46"

### jsonHelpers.go
Custom JSON unmarshal helpers:
- UnixTimeFromIntString (JSON string holding Unix seconds) -> time.Time
- UnixTimeFromInt (JSON number Unix seconds) -> time.Time
- Uint64FromString (JSON string number) -> uint64

Example:

  type Obj struct {
      When util.UnixTimeFromIntString `json:"when"`
  }

  var o Obj
  _ = json.Unmarshal([]byte(`{"when":"1700000000"}`), &o)
  t := o.When.Value() // time.Time

### math.go
- AlmostEqual[T ~float32|~float64](a, b T) bool
  Absolute epsilon comparison at 1e-6.

Example:

  if util.AlmostEqual(0.300000, 0.1+0.2) { /* ... */ }

### pubsub.go
- PubSub[T any]
  Simple fan-out to registered channels.

Example:

  ps := util.NewPubSub[int]()
  ch := ps.Register(1)
  ps.Broadcast(42)
  v := <-ch // 42
  ps.Unregister(ch)

Notes:
- Register with buffer size 0 for synchronous send.

### randomQueue.go
- RandomQueue[T any]
  Mutex-protected bag with O(1) average random removal via swap-with-last.

Example:

  rq := util.NewRandomQueue[string]()
  rq.Add("a"); rq.Add("b"); rq.Add("c")
  x := rq.GetAndRemove() // random item
  n := rq.Len()

Notes:
- GetAndRemove panics if empty. Guard with Len().

### strings.go
- IsASCIIDigits(s string) bool
  True if s is non-empty and all runes are '0'..'9'.

Example:

  util.IsASCIIDigits("12345") // true
  util.IsASCIIDigits("12a45") // false

### timeSeries.go
- TimeSeries[T any]
  Map-like time series keyed by truncated time to a given precision.

Key methods:
- NewTimeSeriesMap(precision time.Duration) *TimeSeries[T]
- Get(t time.Time) *T
- Update(t time.Time, init func() *T, update func(*T))
- Len() int
- Clear()
- Keys() []time.Time (sorted ascending)

Example:

  ts := util.NewTimeSeriesMap[int](time.Minute)
  now := time.Now()
  ts.Update(now, func() *int { v := 1; return &v }, func(p *int) { *p += 1 })
  v := ts.Get(now)
  keys := ts.Keys()

Notes:
- Not concurrency-safe.

### valueTracker.go
- TrackedValue[T constraints.Float]
  Tracks a numeric value and renders a formatted string with a trend symbol using a MovingAverage.

Key functions:
- NewTrackedValue(v T, isCurrency bool, window uint) *TrackedValue[T]
- (t *TrackedValue[T]) Update(v T) bool
- (t *TrackedValue[T]) Value() T
- (t *TrackedValue[T]) String() string

Example:

  tv := util.NewTrackedValue[float64](123.45, true, 10)
  changed := tv.Update(124.00)
  s := tv.String() // e.g., "$124.00↑"

Notes:
- Uses DollarsFormat or NumberFormat depending on isCurrency.

### vector.go
- Vector2D[T Signed]
  Generic 2D vector supporting ints and floats with common operations.

Key methods:
- NewVector2D(dx, dy T) *Vector2D[T]
- Magnitude() float64
- ThetaRadians()/ThetaDegrees()
- Add(*Vector2D[T]) *Vector2D[T]
- Scale(f float64) *Vector2D[T]
- ScaleXY(fx, fy float64) *Vector2D[T]
- ScaleToMaxComponent() *Vector2D[T]
- NormalizeToUnit() *Vector2D[T]
- Rotate(theta float64) *Vector2D[T]
- Translate(dx, dy T) *Vector2D[T]
- X(), Y() T
- String() string

Example:

  v := util.NewVector2D(3.0, 4.0)
  mag := v.Magnitude()           // 5
  deg := v.ThetaDegrees()        // 53.13...
  v.Scale(2).Translate(-1, 0)    // chainable

### workerPool.go
- WorkerPool[W any, R any]
  Generic worker pool with bounded backlog and N workers. Tracks active count.

Key functions:
- NewWorkerPool(f func(W) R, backlog int, numWorkers int) *WorkerPool[W,R]
- (wp *WorkerPool[W,R]) Post(w W)
- (wp *WorkerPool[W,R]) Result() R
- (wp *WorkerPool[W,R]) Len() int32
- (wp *WorkerPool[W,R]) IsActive() bool
- (wp *WorkerPool[W,R]) Close()

Example:

  pool := util.NewWorkerPool(func(n int) int { return n*n }, 16, 4)
  for i := 0; i < 10; i++ { pool.Post(i) }
  go func() { pool.Close() }()
  for i := 0; i < 10; i++ { r := pool.Result(); _ = r }

Notes:
- Panics if backlog < 0 or numWorkers < 1.

## Running tests

  go test ./...

## Contributing

Issues and PRs are welcome.

## License

MIT © 2025 kpratt
