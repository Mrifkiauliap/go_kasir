# 1. Health Check
# Memastikan server berjalan dan merespons.
GET http://localhost:8080/health
HTTP 200
[Asserts]
jsonpath "$.status" == "OK"

# =============================================
# Produk API Tests
# =============================================

# 2. Get All Produk
GET http://localhost:8080/produk
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" count == 3


# 3. Create a New Produk
POST http://localhost:8080/produk
Content-Type: application/json
```json
{
  "nama": "Kopi Kapal Api",
  "harga": 2500,
  "stok": 200
}
```
HTTP 201
[Asserts]
jsonpath "$.nama" == "Kopi Kapal Api"
jsonpath "$.id" isInteger
[Captures]
newProdukId: jsonpath "$.id"


# 4. Get the Created Produk by ID
GET http://localhost:8080/produk/{{newProdukId}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{newProdukId}}
jsonpath "$.nama" == "Kopi Kapal Api"


# 5. Update the Produk
PUT http://localhost:8080/produk/{{newProdukId}}
Content-Type: application/json
```json
{
  "nama": "Kopi ABC Susu",
  "harga": 3000,
  "stok": 150
}
```
HTTP 200
[Asserts]
jsonpath "$.id" == {{newProdukId}}
jsonpath "$.nama" == "Kopi ABC Susu"
jsonpath "$.harga" == 3000


# 6. Delete the Produk
DELETE http://localhost:8080/produk/{{newProdukId}}
HTTP 204


# 7. Verify Produk Deletion
GET http://localhost:8080/produk/{{newProdukId}}
HTTP 404


# =============================================
# Kategori API Tests
# =============================================

# 8. Get All Kategori
GET http://localhost:8080/categories
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" count == 3


# 9. Create a New Kategori
POST http://localhost:8080/categories
Content-Type: application/json
```json
{
  "name": "Bumbu Dapur",
  "description": "Bumbu untuk memasak"
}
```
HTTP 201
[Asserts]
jsonpath "$.name" == "Bumbu Dapur"
jsonpath "$.id" isInteger
[Captures]
newCategoryId: jsonpath "$.id"


# 10. Get the Created Kategori by ID
GET http://localhost:8080/categories/{{newCategoryId}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{newCategoryId}}
jsonpath "$.name" == "Bumbu Dapur"


# 11. Update the Kategori
PUT http://localhost:8080/categories/{{newCategoryId}}
Content-Type: application/json
```json
{
  "name": "Bumbu Instan",
  "description": "Bumbu siap pakai"
}
```
HTTP 200
[Asserts]
jsonpath "$.name" == "Bumbu Instan"
jsonpath "$.description" == "Bumbu siap pakai"


# 12. Delete the Kategori
DELETE http://localhost:8080/categories/{{newCategoryId}}
HTTP 204


# 13. Verify Kategori Deletion
GET http://localhost:8080/categories/{{newCategoryId}}
HTTP 404
