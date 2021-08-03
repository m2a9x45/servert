# Servert

Thr website allows customers to rent virtual private servers (VPS). It integrates with stripe to support billing and I built
some internal tooling to support setting the servers up.

The backend is written in Go this was one of my first experinces writting Go. I used gorilla mux as easy way to spin up an API. The APi covers a lot of features, from receipt data
to the ordering flow to talking with stripe to actually pay for an order. This was also one of the first times I intergarted with Stripe.


