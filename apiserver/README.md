# API SERVER

- [技術工具](#技術工具)
- [負責業務](#負責業務)
- [Others](#Others)

# 技術工具
- 選用Gin:
    - why?
      - 熱門且文件充足
      - doc: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
- 選用GraphQL
    - gateway的概念, 不用開一堆api path 出去，統一一個端口
    - GraphQL lib: https://gqlgen.com/
    - why?
        - 效能考量 && schema first && code gen 可以少打代碼
        - 範例多: https://github.com/99designs/gqlgen/tree/master/_examples
        - reference: https://blog.wu-boy.com/2020/04/switch-graphql-go-to-gqlgen-in-golang/

# 負責業務
- 透過GraphQL提供API給前端
- 接到前端的請求後，透過GraphQL的resolver, 轉打底層Microservice的Grpc服務
- api server 只會有連線Grpc, 不會連線任何Infra(ex: DB, cache), 唯一的連線Infra, 就是Elasitcsearch, 寫入log

# Others
- graphQL GUI tool: Altair: https://altairgraphql.dev/#download