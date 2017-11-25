# bananamail

 This repository is able to send the mail you wanna send the contents text to bananaman no [bananamoonGOLD](https://www.tbsradio.jp/banana/).
 
 
 
## Description
 
It is possible to send it according to the theme is subject of bananamoon mail you want to send.

## Installation

```
$ go get github.com/soeyusuke/bananaman/bananamail
```

## Usage
You can use each sub-commands
Please see `$ bananamail help`.  

### help
```sh
$ bananamail -h

Usage: Gmail to bananamoon cli [--version] [--help] <command> [<args>]

Available commands are:
    ensyutu       This is 'ensyutu|演出' of Subject
    henken        This is 'henken|偏見' of Subject
    hiromenesu    This is 'hiromenesu|ヒロメネス' of Subject
    init          This is able to set your address, gmail app pass, your radio name
    sengen        This is 'sengen|宣言' of Subject
    theme         This is 'Theme|テーマ' of Subject


```

### init
set up the your gmail address, **your gmail app pass**, radio name

```
$ bananamail init 'your address' 'your pass' 'radio name'
```
If you don't know how to make gmail app password, you should see this site.
[Google Account Help](https://support.google.com/accounts/answer/185833?hl=ja&ctx=ch_b%2F0%2FDisplayUnlockCaptcha) 
