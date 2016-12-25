# Google Font Downloader

[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://mit-license.org/2016)

Google Fonts下载器，可以根据参数或样式表下载字体，并生成新的使用下载字体的样式表。

## 安装

1. `go get -u github.com/fate-lovely/google-font-downloader`
2. 在[Releases](https://github.com/fate-lovely/google-font-downloader/releases)页面下载

## 用法

### 概要

`google-font-downloader [flags...] [specs...]`

### 选项

- `-o, —output=output.css`

指定输出样式表的路径。

- `-f, —format=woff`

指定下载字体的格式，默认为`woff`，根据[caniuse.com](http://caniuse.com/#search=woff)的数据，`woff`目前拥有最好的兼容性，可选的格式为`eot`, `woff`, `woff2`, `svg`,`ttf`。

- `-i, —input`

指定输入样式表的路径，如果这一项不为空的话，那么后面的字体参数将被忽略，程序将分析这个样式表，下载样式表中的字体，然后输出使用下载字体的新样式表。当你在Google Fonts上挑选好字体后，可以使用这个选项来转换Google Fonts生成的样式表。

- `-l, —lang=latin`

指定字体的语言，使用逗号分隔，例如`-l latin,greek`，默认为`latin`。

### 字体

字体的格式为`[name]:[weights...]`，例如`Roboto:300,400,500`，可以使用`i`表示斜体，`b`表示粗体。

## 示例

- `google-font-downloader "Roboto:300,300i,400"`
- `google-font-downloader "Roboto:300,b"`
- `google-font-downloader "Open Sans" "Roboto"`，如果不指定weight，默认下载regular的字体
- `google-font-downloader -f svg "Roboto"`
- `google-font-downloader -i input.css`



