# HTTP Server

## TODO List App

### Description

Aplikasi bernama todo app yang berfungsi untuk membuat list beserta status pekerjaan dari aktifitas yang kita kerjakan. Fitur dari aplikasi ini adalah:

- Register
- Login
- Create Todo list
- Read Todo list
- Clear Todo List
- Logout

Terdapat chain middleware untuk menghandle Method dan Authentication dengan menggunakan metode session based token kemudian menyimpan semua data user dan todo list di _in memory map_.