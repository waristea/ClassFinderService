# ClassFinderService
Web Service yang berisi informasi tentang mata kuliah yang terdapat di ITB. 
Web Service ini bekerja dengan cara meng-*scrap* data yang terdapat pada akademik.itb.ac.id menggunakan credentials yang disediakan. 

## Cara Pemakaian (Manual)
1. Download atau clone project Github ini
2. (Opsional) Bila ingin melakukan scrap, masukan credentials anda (nim, username, password) pada fungsi insertDB(), kemudian uncomment fungsi tersebut pada main().
3. Saat berada di direktori project jalankan *go run main.go dataFetch.go* dan *mongod --dbpath="./db"* pada dua terminal yang berbeda.
4. Lihat data yang ada dengan membuka **127.0.0.1:8080/schedules** pada browser
