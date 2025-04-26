@echo off
REM 使用方法：
REM genModel.bat usercenter user
REM genModel.bat usercenter user_auth

set tables=%2
set modeldir=./genModel

set host=sh-cynosdbmysql-grp-o5q6mlei.sql.tencentcdb.com
set port=22106
set dbname=graph_med_%1
set username=root
set passwd=zxc?473027362

echo 开始创建库：%dbname% 的表：%tables%
goctl model mysql datasource -url="%username%:%passwd%@tcp(%host%:%port%)/%dbname%" -table="%tables%" -dir="%modeldir%" --cache=true --style=goZero