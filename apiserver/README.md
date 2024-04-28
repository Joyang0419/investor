# API SERVER

- 選用Gin:
  - why?
    - 熱門且文件充足
  - doc: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
- 選用GraphQL
  - GraphQL lib: https://gqlgen.com/
  - why?
    - 效能考量 && schema first && code gen 可以少打代碼 
    - 範例多: https://github.com/99designs/gqlgen/tree/master/_examples
    - reference: https://blog.wu-boy.com/2020/04/switch-graphql-go-to-gqlgen-in-golang/
  - gqlgen 「無法解決」的議題
    - model_gen 無法依照schema分開放: 
    - https://github.com/99designs/gqlgen/issues/1265
    - https://github.com/99designs/gqlgen/issues/2399

# TODO example, 方便先熟悉GraphQL, 之後準備刪除, 找一款好用的GraphQL client
```graphql 
query {
  Investors(
    queryParams: {
      investorIDs: [], 
      loginAccounts: [], 
      pageSize: 10, 
      page: 1
    }
  ) {
    investorID
    loginAccount
  }
}

```

```graphql
mutation {
  createInvestor(
    input: {
      loginAccount: "Hello", 
      password: "qq1234",
    }
  	) 
  {
    investorID
    loginAccount
    password
  }
}
```