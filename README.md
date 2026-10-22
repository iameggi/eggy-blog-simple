# eggy.me

SSG kecil berbasis Go, `html/template`, dan [Goldmark](https://github.com/yuin/goldmark/tree/v1.8.6). HTML/CSS tanpa JavaScript dan tanpa proses build frontend.

## Jalankan

Dengan Go 1.22 atau lebih baru, jalankan dari direktori proyek:

```sh
go mod download
go run .
```

Buka `public/index.html` langsung di browser atau unggah isi `public/` ke document root `www.eggy.me`. Semua tautan navigasi memakai path relatif.

## Struktur

```text
main.go
go.mod
go.sum
template.html
static/
  images/
    about-me.jpg
    gearx.jpg
content/
  blog/
    (22-10-2026)membuat-gearx.md
  project/
    (22-10-2026)gearx.md
public/
  index.html
  images/
    about-me.jpg
    gearx.jpg
  blog/
    membuat-gearx.html
  project/
    gearx.html
```

Nama konten mengikuti `(DD-MM-YYYY)slug-judul.md`. Gunakan huruf ASCII, angka, dan tanda hubung antar kata pada slug. Judul ditampilkan persis sebagai slug, termasuk tanda hubung. Tanggal divalidasi dan ditampilkan sebagai `DD-MM-YYYY`; nama output hanya memakai slug.

Daftar blog dan proyek diurutkan dari tanggal terbaru, lalu slug untuk tanggal yang sama. Slug duplikat dalam kategori yang sama menghasilkan error, termasuk nama yang hanya berbeda huruf besar/kecil. Slug yang sama boleh digunakan pada kategori berbeda.

Folder konten kosong atau belum ada menghasilkan daftar kosong. Subfolder dan file selain `.md` diabaikan. Generator membaca konten dari satu tingkat direktori saja.

`template.html` mengatur layout, warna, serta bagian statis About Me dan GearX. Judul artikel sudah dibuat oleh template; mulai isi Markdown dengan paragraf atau heading `##`. Goldmark menggunakan pengaturan bawaan tanpa raw HTML, sehingga tag script dari Markdown tidak dirender.

About Me dan GearX berisi teks dummy Lorem ipsum dan gambar lokal. Ganti teks pada `template.html` dan gambar pada `static/images/`, lalu jalankan ulang `go run .`. Semua isi `static/` otomatis disalin ke `public/` saat build. Foto dummy berasal dari Lorem Picsum: [About Me](https://picsum.photos/id/0/info) dan [GearX](https://picsum.photos/id/180/info). Efek hitam-putih diterapkan melalui CSS.

File contoh dapat diedit atau dihapus. Build menulis ulang output dengan nama yang sama; setelah menghapus atau mengganti nama konten, hapus juga HTML lamanya dari `public/`.
