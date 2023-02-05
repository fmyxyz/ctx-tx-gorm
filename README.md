# go-ctx-tx

通过context.Context控制事务行为

## Installation

    go get github.com/fmyxyz/go-ctx-tx@latest

## Overview

已实现的事务行为:

1. Required 没有事务时创建，有事务时加入
2. Supported 有事务时加入
3. Nested 没有事务时创建，有事务时嵌套
4. New 总是新建事务

## Usage

#   c t x - t x - g o r m  
 