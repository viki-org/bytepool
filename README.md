### BytePool
BytePool manages a thread-safe pool of `[]byte`. By using a pool of pre-allocated arrays, one reduces the number of allocations (and deallocations) as well as reducing memory fragmentation.

If the pool is empty, new items will be created on the fly, but the size of the pool will not grow. Furthermore, the returned items are fixed-length `[]byte` - they will not grow as needed. The idea is for you to favor over-allocation upfront.

#### NOTE
Perhaps a slightly more generic version of pool will be coming in Go 1.3, as in [sync.Pool](http://tip.golang.org/pkg/sync/#Pool)

### Installation

    go get github.com/viki-org/bytepool.v1

### Example
A common example is reading the body of an HTTP Request. The memory-unfriendly approach is to do:

    body, _ := ioutil.ReadFull(req.Body)

A slightly better approach would be to predefine the array length:

    body := make([]byte, req.ContentLength)
    io.ReadFull(req.Body, body)

While the 2nd example avoids any over-allocation as well reallocation from a dynamically growing buffer, it still creates a new array (a new array which will need to be garbage collected).

This allocation can be avoided by using a pool of `[]byte`:

    //pre-allocates 256MB (8K arrays of 32K bytes each)
    var pool = bytepool.New(8196, 32768)
    func handler(res http.ResponseWriter, req *http.Request) {
      buffer := pool.Checkout()
      defer buffer.Close()
      buffer.ReadFrom(req.Body)
      body := buffer.Bytes()
      ...
    }

The above generates a pool of 8K `[]byte` each of which can hold 32K of data. An array is retrieved via the `Checkout` method and returned back to the pool by calling `Close`.

### Methods
The item returned from the pool implements a number of common interfaces, such as `io.Closer`, `io.Writer`, `io.Reader` and `io.ReaderFrom`.

You can get the returned value as `Bytes()` or `String()`

### Json
If the buffer will be used to generate JSON, consider creating a `JsonPool` instead:

    var pool = bytepool.NewJson(8196, 32768)

Items returned from a `JsonPool` have a number of helper methods for writing JSON (in addition to most of the methods of the base type):

    buffer := pool.Checkout()
    buffer.BeginArray()
    for _, id := range ids {
      buffer.WriteInt(id)
    }
    buffer.EndArray()

The above will take care of properly delimiting the array. Similar behavior can be achieved with the `BeginObject`, `EndObject` and the various key-value helpers: `WriteKeyString`, `WriteKeyInt`, `WriteKeyBool`, `WriteKeyTime` (time.Time).

Key values are expected to be escaped. String values will automatically be escaped. This can be circumvented by using the alternative `WriteSafeString` and `WriteKeySafeString` methods.

You can also use the `WriteKeyValue` to append other JSON object/array that is passed as an string. You should use this method carefully. E.g:

  buffer.BeginObject()
  buffer.WriteKeyString("name":"tyler")
  buffer.WriteKeyValue("metadata", `{"age": 12}`)
  buffer.EndObject()
  println(buffer.String()) // outputs: {"name":"tyler","metadata":{"age":12}}

### Credits
Bytepool is open-sourced, used and maintained by [Viki](https://github.com/viki-org).
Much of the work goes to [Karl](https://github.com/karlseguin), [Cristobal](https://github.com/cviedmai)
