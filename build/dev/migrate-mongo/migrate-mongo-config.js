// migrate-mongo-config.js
module.exports = {
    // MongoDB 的連接設定
    mongodb: {
        // 注意：此處 URL 需要根據實際環境設定，如果是使用 Docker Compose 管理容器，
        // 你應該使用服務名（這裡是 "mongodb"）作為主機名。
        url: "mongodb://root:root@mongodb:27017/admin",

        // MongoDB 的連接選項
        options: {
            useNewUrlParser: true,      // 使用新的 URL 字串解析器
            useUnifiedTopology: true,   // 使用統一的拓撲結構
        }
    },

    // 遷移腳本存放的目錄
    migrationsDir: "migrations",

    // 遷移狀態記錄的集合名稱
    changelogCollectionName: "changelog",

    // 遷移腳本文件的擴展名
    migrationFileExtension: ".js"
};
