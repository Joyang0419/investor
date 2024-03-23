// init.db.js
db = db.getSiblingDB('investor');

db.createCollection('stock_daily_price');