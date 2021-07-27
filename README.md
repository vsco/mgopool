## ***** DEPRECATION NOTICE ******
The mongo library that this project depends on has been deprecated ([mgo](https://github.com/globalsign/mgo)), in favor of the 
[official driver released by mongodb](https://github.com/mongodb/mongo-go-driver). The official driver includes an internal
mongo connection pool by default, which obviates this library.


# mgopool<br> [![GoDoc](https://godoc.org/github.com/vsco/mgopool?status.svg)](https://godoc.org/github.com/vsco/mgopool) [![Build Status](https://travis-ci.org/vsco/mgopool.svg?branch=master)](https://travis-ci.org/vsco/mgopool)

Package mgopool provides buffer implementations around [mgo.Session (v2)][mgo].

```
go get github.com/globalsign/mgo
go get github.com/vsco/mgopool
```

With mgo, an initial session is created with the server, and then either a [Copy][Session.Copy] or [Clone][Session.Clone]
of that initial session is used for each unit of work (HTTP request handler, worker loop, etc.), which is
[closed][Session.Close] upon completion. This allows mgo to efficiently manage a pool of sockets between those sessions;
however, it also results in the underlying sockets to constantly log in and out of the Mongo cluster.

mgopool avoids this issue by creating a Pool of copies from the initial session, storing them in a free list. Work units
request an existing copy (or one is created if it doesn't exist) and then return it to the pool on completion. Both a
leaky and capped version of the Pool is available, depending on need.

## Example Usage

```go
func Example() {
  // create the initial mgo.Session, defer closing it
  initial, _ := mgo.Dial("127.0.0.1")
  defer initial.Close()

  // create a new leaky pool with size of 3
  pool := mgopool.NewLeaky(initial, 3)
  defer pool.Close()

  wg := sync.WaitGroup{}

  // create a few workers to utilize the sessions in the pool
  for i := 0; i < 3; i++ {
    wg.Add(1)

    go func() {
      for j := 0; j < 2; j++ {
        // get a session from the pool
        session := pool.Get()

        // do something with that session, refresh it if there is an error
        if err := doSomething(session); err != nil {
          session.Refresh()
        }

        // return the session to the pool
        pool.Put(session)
      }
      wg.Done()
    }()
  }

  // wait for the workers to complete
  wg.Wait()

  // Output:
  // foobar
  // foobar
  // foobar
  // foobar
  // foobar
  // foobar
}

func doSomething(s *mgo.Session) error {
  fmt.Println("foobar")
  return s.Ping()
}
```

[mgo]:           https://godoc.org/github.com/globalsign/mgo
[Session.Copy]:  https://godoc.org/github.com/globalsign/mgo#Session.Copy
[Session.Clone]: https://godoc.org/github.com/globalsign/mgo#Session.Clone
[Session.Close]: https://godoc.org/github.com/globalsign/mgo#Session.Close

## Testing

Testing mgopool requires a running Mongo cluster. The default Mongo host is used (`127.0.0.1:27017`) and
can be overridden by setting the `MONGO_HOST` env variable. Execute the tests with the following command from the project root:

```sh
script/test
```

## License

The MIT License (MIT)

Copyright (c) 2016 Visual Supply, Co.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
