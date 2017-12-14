# ClassFinderService
Web Service yang berisi informasi tentang mata kuliah yang terdapat di ITB. 
Web Service ini bekerja dengan cara meng-*scrap* data yang terdapat pada akademik.itb.ac.id menggunakan credentials yang disediakan. 
Implementasi dari web ini dapat diakses pada: http://167.205.67.226:8080/schedules

## Dependencies
- github.com/gorilla/context
- golang.org/x/net/html
- gopkg.in/mgo.v2

## Cara Pemakaian (Manual)
1. Download atau clone project Github ini
2. (Opsional) Bila ingin melakukan scrap ulang, ketikkan 'y' pada prompt yang muncul. Kemudian masukkan credentials.
3. Saat berada di direktori project jalankan *go run main.go dataFetch.go*.
4. Lihat data yang ada dengan membuka **http://167.205.67.226:8080/schedules** pada browser

## Query String yang dapat dimasukkan
### code 
Kode mata kuliah
Contoh : http://167.205.67.226:8080/schedules?code=MA1101
### subject
Nama mata kuliah 
Contoh : http://167.205.67.226:8080/schedules?subjek=Pemrograman%20Integratif 
