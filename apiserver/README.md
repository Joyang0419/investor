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
  - graph 的專案架構
    ```text
    ├── go.mod
    ├── go.sum
    ├── gqlgen.yml
    ├── graph
    │   ├── generated
    │   ├── model_gen
    │   └── modules
    │       ├── todo
    │       │   ├── mutator.go
    │       │   ├── resolver.go
    │       │   └── typedef.graphql
    │       └── user
    │           ├── mutator.go
    │           ├── resolver.go
    │           └── typedef.graphql
    └── server.go    
    ```
  - gqlgen 「無法解決」的議題
    - model_gen 無法依照schema分開放: https://github.com/99designs/gqlgen/issues/1265
      - 就先這樣吧。專注在實作modules 內的 mutator and resolver