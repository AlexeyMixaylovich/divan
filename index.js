const _ = require('lodash');
const fs = require('fs');

var cron = require('node-cron');
const { sendNewItem, sendUpdate, init } = require('./bot');
const { getNewData , writeDB} = require('./pages');

async function compare(data) { 
    const bestDB = require('./db.json') || [];
    const dbByAtr = _.keyBy(bestDB, 'atr')
    
    const updateItem = (item) => { 
        const index = bestDB.findIndex(({ atr }) => atr === item.atr);
        bestDB[index] = item;
    }
        
    data.forEach((data) => {
        const { atr, price } = data

        const fromDB = dbByAtr[atr];

        if (!fromDB) { 
            sendNewItem(data)
            console.error('Новая позиция', data);
            bestDB.push(data)
            return;
        }
        
        if (fromDB.price === price)
            return;
        
        if (fromDB.price > price) {
            sendUpdate(data, fromDB.price);
            updateItem(data);
            console.error('Цена упала', atr,  fromDB.price, price);
        }
        
        if (fromDB.price < price) {
            sendUpdate(data, fromDB.price);
            updateItem(data);
            console.error('Цена выросла', atr, fromDB.price, price);
        }
        
    })

    fs.writeFileSync('./db.json', JSON.stringify(bestDB, null, 2))

}
async function update() { 
    const data = await getNewData();
    await compare(data);
}

(async function run() { 
    init();
    // await writeDB()
    cron.schedule('*/5 * * * *', async () => {
        console.log('running');
        await update()
            .then(() => { console.log('done') })
            .catch(e => console.error(e));
    });

})()



