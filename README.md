# Failure
The failure package builds upon the [errors](https://github.com/pkg/errors) 
package, implementing a strategy called ``Opaque errors`` which I first learned 
about from an article [Don't just check errors handle the gracefully]
(https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
by [Dave Cheney](https://dave.cheney.net) who I believe is also the author of the `errors` package. 



Failure is geared towards describing errors that occur while using 
microservices that generally support rest api's and as such you will notice 
a slight bias towards that, although I try to separate concerns as much as 
possible. It is heavily influenced by [sls](https://github.com/rsb/sls) a 
library used to develop AWS Serverless applications using `Golang` and 
`Terraform`. 

