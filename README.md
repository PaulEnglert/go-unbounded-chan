# go-unbounded-chan


[![version](https://img.shields.io/github/go-mod/go-version/PaulEnglert/go-unbounded-chan)](https://github.com/PaulEnglert/go-unbounded-chan)



This package provides an unbounded channel (ie. a channel with no buffer, but non-blocking write accessibility). The unbounded channel works by using two separate channels (In and Out) that are used for writing and reading, as well as an internal queue that is used for caching. The order in which the output channel receives items, is the same as the input channel is written to. Once the channel is shutdown with Close(), the queue will be drained automatically.

## Usage


Load the package:

    go get github.com/PaulEnglert/go-unbounded-chan


Use the `uchan` package
    
    import (
        "fmt"
        "sync"
        "github.com/PaulEnglert/go-unbounded-chan"
    )


    var wg sync.WaitGroup

    // create structure
    uc := uchan.NewUnboundedChannel()

    // spawn writer (will not block (!))
    wg.Add(1)
    go func(){
        for i := 0; i < 100; i++ {
            // use UnboundedChannel.In to write
            uc.In <- i
        }
        wg.Done()
    }()

    // do read till closed
    go func() {
        // use UnboundedChannel.Out to read
        for v := range uc.Out {
            vi := i.(int)
            fmt.Printf('Got: %d\n', vi)
        }
    }()

    // wait till all is written
    // then close
    wg.Wait()
    // Do not close In/Out manually, but use
    // Close() to shut down. This will drain
    // the queue automtaically and block until
    // all items have been read.
    uc.Close()


## Limitations

The channel works only with `interface{}` types for the in/out channels, casting on read is therefore always required.
