https://pkg.go.dev/golang.org/x/sync/singleflight

Imagine you have a function that is called by a lot of goroutines, and you want to avoid calling it multiple times, this because this function is expensive and you want to avoid wasting resources, or if you have multiple goroutines it can cause a race conditioning or a side effect in the result of the function.

A simple example is the following:

You have a function that get a key from a database, and you know that this key will be used a lot by different goroutines. You just need to pay the cost to get the key from the database once, and then you can use the result of the function for all the goroutines. Using a cache.

In a concurrent scenario, if you trigger multiple goroutines, all of them will call the function, and at some point, all of them could call this function at the same time and hit the database multiple times. Leaving your code to waste resources and time.

This is where singleflight comes in. Singleflight is a package that provides a way to call a function only once, and then use the result of the function for all the goroutines.

