# GORM by Example

This project is for give examples for how to use GORM to do CRUD.

# Model Dependency Relationship

```
db -> nil
test -> db
models -> (db, test)
```

# JSON Field

https://github.com/go-gorm/datatypes
https://kejyuntw.gitbooks.io/mysql-learning-notes/content/query/query-json-contains.html

# Callback Notes

many2many 的中间表 callback不太好用。所以需要通过自己写Delete方法来cover callback。

```golang
user.Destroy() => 会删除所有的 user 的关联表。
```
