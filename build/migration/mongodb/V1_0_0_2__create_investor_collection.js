module.exports = {
    async up(db, client) {
        // 從 client 對象中獲取 'investor' 數據庫實例
        const db_investor = client.db("investor");

        // 在 'investor' 數據庫中創建一個新的集合 'stock_daily_price'
        await db_investor.createCollection('investor');
    },

    async down(db, client) {
        // 從 client 對象中獲取 'investor' 數據庫實例
        const db_investor = client.db("investor");

        // 在 'investor' 數據庫中創建一個新的集合 'stock_daily_price'
        await db_investor.collection('investor').drop();
    }
};
