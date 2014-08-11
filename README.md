## simplegopherd [![Build Status](https://travis-ci.org/CatFriends/simplegopherd.svg?branch=master)](https://travis-ci.org/CatFriends/simplegopherd)

**Gopher** protocol is designed for distributing documents over the Internet. It was developed in 1991 in University of Minnesota. Because of licensing fees the protocol did not get much use, and is almost forgotten nowadays and sounds like FIDO or BBS.

Development of the server is some kind of tribute to Gopher, and a try to bring new life into this old good technology.

## Installation and Usage
This application uses [**gcfg** library][1] to handle `.ini` configuration files. You have to get it before making a build of the server:

```
go get code.google.com/p/gcfg
```

You will need a `go` compiler and a `git` client installed to build **simplegopherd** from source. Please refer to instructions of your operating system distribution to get these if needed.

Once everything is ready, type the following command in your console:

```
go get github.com/CatFriends/simplegopherd
```

### Configuration

**simplegopherd** uses `.ini`-style configuration files. You will need to create one for your instance:

```
; simpegopherd configuration file

[Network]
HostName = localhost         # bind address or hostname
PortNumber = 70              # port number

[Site]
BaseDirectory = C:/Temp      # root folder to serve
IndexFileName = index.csv    # index file name

[Gopher]
NewLineSequence = \n\r       # new line format
```

### Preparing site

Each and every directory from the Gopher site must contain `index.csv` file. It contains descriptions for files in that directory.

For example, if your site has the following layout:

```
+ root
|-- index.csv
|-- about.txt
|-- cat.gif
```

your `index.csv` might be like this:

```
"Welcome to Sample Gopher Server",""
"",""
"About us","about.txt"
"Our Felix's photo","cat.gif"
```

Keep in mind these rules when writing `index.csv` files:

  - First column is **description** of item, a free text limited with 70 characters
  - Second column is item **selector** and contains file name
  - If item have *empty* selector, it means that it is just an **information** item

### Running

After configuration steps is finished you can start your Gopher site using:

```
simplegopherd <configuration.ini>
```

  [1]: https://code.google.com/p/gcfg/