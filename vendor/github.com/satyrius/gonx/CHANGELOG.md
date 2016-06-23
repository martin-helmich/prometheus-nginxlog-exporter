# Release History

## v1.3.0 (2015-12-19)

### Major features

- Introduce `Filter` interface #19. It acts similar to `Reducer`, but it is about limiting chain entries.
  - Has `Filter(*Entry) *Entry` method to check is entry meets filter condition
  - Implements `Reducer` interface too, to be used in chains
- `Datetime` filter #19. Its based on this PR #11 by @pshevtsov
- Introduce `StringParser` #24 by @pshevtsov


### Minor features

- Run test for `go` version up to `1.5`
- Linting fixes #13, #14 (thanks to @pshevtsov)
- More examples, espetially for reducers #17
- Use `goconvey` for tests

### Bugfixes

- Fixed `Reader` examples #21
- Long lines reading #23 by @pshevtsov
- Fix nginx conf parsing, deal with commented lines #25 by @jack1582

## v1.2.2 (2014-11-05)

### Bugfixes

- Parsing last value without quotes was fixed, #9
- Tested for Go v1.3

## v1.2.1 (2014-06-21)

### Bugfixes

- Fix issue #6 which causes parser to crash if some value in a log line was empty
- Fix `TestGroupByReducer` #4, it was random crashing, because we could not expect the order of entries readed from an output channel

## v1.2.0 (2014-03-30)

### Features and Improvements

* The aggregation reducers such as `Avg`, `Sum` and `Count` was introduces along with `Chain` and `GroupBy` reducers.
* `Entry` got some new methods
  * Getters `Field(name string)`, `FloatField(name string)`
  * Setters `SetField(name string, value string)`, `SetFloatField(name string, value float64)`, `SetUintField(name string, uint64)`
  * Utility methods `Merge(entry *Entry)`, `FieldsHash(fields []string)`, `Partial(fields []string)`

### Backward incompatibilities

* All functions deals with `*Entry` instead of `Entry`
* `MapReduce` returns `chan *Entry` instead of `chan interface{}` and all reducers accept output channel as `chan *Entry`
* `Entry` is a `struct`, not a `map[string]string` anymore and has two constructors `NewEntry` that accepts `Fields` and `NewEmptyEntry`
* `Entry.Get` was renamed to `Entry.Field`

## v1.1.0 (2013-11-20)

### Major feature

Implement function MapReduce to parse log file in asynchronous manner for speed improvement. `Reader.Read` and constructors signatures and behaviour still the same.

## v1.0.0 (2013-11-11)

Log reader type `Reader` with the following constructors
* `func NewReader(logFile io.Reader, format string) *Reader`
* `func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error)`

And one interface method
* `func (r *Reader) Read() (record Entry, err error)`
