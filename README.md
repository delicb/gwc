# GWC (Go Web Client)
GWC is GoLang HTTP client based on [Cliware](https://github.com/delicb/cliware)
client middleware management library and its [middlewares](https://github.com/delicb/cliware-middlewares).

Basic idea is to use same middleware mechanism on client side as many projects
use it for server development. 

Because it is based on middlewares client is very pluggable. However, even out 
of box it should support most common use cases.

# State
This is early development, not stable, backward compatibility not guarantied.
**Not recommended for use in production yet**.

# Credits
Idea and bunch of implementation details were taken from cool GoLang HTTP client
[Gentleman](https://github.com/h2non/gentleman).

Difference is that GWC is based on Cliware, which supports `context` for client
requests and has more simple idea of middleware. Also, GWC is lacking some
features that Gentleman has (like mux). For now I do not plan on adding them
to GWC, but I might write middleware that support similar functionality in the
future. 
