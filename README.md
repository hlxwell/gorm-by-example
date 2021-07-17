# GORM by Example

This project is for give examples for how to use GORM to do CRUD.

# Features

- Load Fixture
- Scope
- Many to many
- Has Many
- JSON Field
- Validate

# JSON Field

https://github.com/go-gorm/datatypes
https://kejyuntw.gitbooks.io/mysql-learning-notes/content/query/query-json-contains.html

# Callback Notes

many2many 的中间表 callback 不太好用。所以需要通过自己写 Delete 方法来 cover callback。
比如在例子中 user_role 属于中间表，你在 user_role 里面添加 AfterDelete 的时候，期待他被删除的时候会被触发，其实不会。

```golang
user.Destroy() => 会删除所有的 user 的关联表。

func (user *User) Destroy() {
  // 删除所有的 associated objects
}
```
