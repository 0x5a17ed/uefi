# uefi

[![License: APACHE-2.0](https://img.shields.io/badge/license-APACHE--2.0-blue?style=flat-square)](https://www.apache.org/licenses/)

A UEFI library written in go to interact with efivars. Compatible with Windows and Linux.


## 📦 Installation

```shell
$ go get -u github.com/0x5a17ed/uefi@latest
```


## 🤔 Usage

```go
package main

import (
	"fmt"

	"github.com/0x5a17ed/uefi/efi/efivaraccess"
	"github.com/0x5a17ed/uefi/efi/efivars"
)

func main() {
	c := efivaraccess.NewDefaultContext()

	if err := efivars.BootNext.Set(c, 1); err != nil {
		fmt.Println(err)
		return
	}
}
```


## 💡 Features
- Works on Linux and on Windows
- Reading Boot options
- Setting next Boot option
- Extensible
- Simple API
