- 标量类型（Scalar Types）
  - Int：表示有符号32位整数。
  - Float：表示双精度浮点值。
  - String：表示UTF-8字符序列。
  - Boolean：表示true或false。
  - ID：表示一个唯一标识符，通常用于重新获取对象或作为缓存的键。尽管它被序列化为字符串，但不应该被看作是可读的。
- 对象类型（Object Types）
  - Query：查询类型，用于查询数据。
  - Mutation：变更类型，用于修改数据。
  - Subscription：订阅类型，用于实时数据推送。
```graphql
query {
  user(id: "123") {
    name
    friends {
      name
    }
  }
}

type Mutation {
  createUser(name: String!, email: String!): User
  updateUser(id: ID!, name: String, email: String): User
  deleteUser(id: ID!): User
}

type Subscription {
  newUser: User
}
```
- 接口（Interface Types）
  - Node：所有对象类型的基本接口。
```graphql
# 定义一个接口
interface Character {
  id: ID!
  name: String!
}

# 实现接口的具体类型
type Human implements Character {
  id: ID!
  name: String!
  height: Float
}

type Alien implements Character {
  id: ID!
  name: String!
  planet: String
}
```
- 联合类型（Union Types）
  - SearchResult：搜索结果类型。
```graphql
# 定义两种可能的对象类型
type Article {
  title: String!
  author: String!
  content: String!
}

type Video {
  title: String!
  director: String!
  duration: Int!
}

# 定义一个联合类型，包括上述两种类型
union SearchResult = Article | Video

# 定义查询类型
type Query {
  search(keyword: String!): [SearchResult!]
}
```
- 枚举类型（Enum Types）
```graphql
enum Season {
  SPRING
  SUMMER
  AUTUMN
  WINTER
}
```
- 输入类型（Input Types）
```graphql
input UserProfileInput {
  firstName: String
  lastName: String
  age: Int
}
```
- 列表和非空类型（List and Non-Null Types）
  - 列表：字段可以是列表类型，这表示该字段可以返回多个值。列表类型用方括号表示，如[Type]。
```graphql
type User {
  id: ID!
  emails: [String]
}
```
  - 非空：任何类型都可以被标记为非空，使用一个感叹号（!），表示这个字段在查询中返回的结果总是有值，从不为null。
```graphql
type User {
  id: ID!
  emails: [String]!
}
```

# Mutation
- 变更类型（Mutation）用于修改数据，它是GraphQL的一个特殊类型，用于执行写操作，如创建、更新、删除等。
- 不要去想「Nested Mutation」，撞過了，找不到答案，命名明確就好。
