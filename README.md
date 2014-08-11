# simplegopherd [![Build Status](https://travis-ci.org/CatFriends/simplegopherd.svg?branch=master)](https://travis-ci.org/CatFriends/simplegopherd)

**Gopher** protocol is designed for distributing documents over the Internet. It was developed in 1991 in University of Minnesota. Because of licensing fees the protocol did not get much use, and now is almost forgotten. Nowadays it sounds like FIDO or BBS.

Development of the server is some kind of tribute to Gopher, and a try to bring new life into this old good technology.

## Prerequisites
This application uses [**gcfg** library][1] to handle `.ini` configuration files. You have to get it before making a build of the server:

```
go get code.google.com/p/gcfg
```

  [1]: https://code.google.com/p/gcfg/
