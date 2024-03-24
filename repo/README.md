- 原則
  - repo 命名: 依照table name || collection name
  - 如果用sub query && 或者是臨時表 && join，原則不變，就是歸在「第一張表」。
  ```text
  # example: 這個就是歸類在 A repo
    SELECT * FROM A INNER JOIN B
  ```