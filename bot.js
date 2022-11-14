const TelegramBot = require('node-telegram-bot-api')
const fs = require('fs');

const { bot: { token } } = require('./config.json');


const Users = require('./users.json') || [];

const bot = new TelegramBot(token, { polling: true })

const options = { parse_mode: 'HTML', disable_web_page_preview: true };


async function addUser(chatId) { 
    if (!Users.includes(chatId))
        Users.push(chatId);
    
    fs.writeFileSync('./users.json', JSON.stringify(Users, null, 2))
}


function init() { 
    bot.on('message', (msg) => {
        const chatId = msg.chat.id
        bot.sendMessage(chatId, "А всё, Надо было раньше!")
        addUser(chatId)
    })
}



function getTex(item, before = []) {
    const { name, price, oldPrice, sale, url } = item;

    const text = [
        ...before,
        name, ' ',
        `<b>${price} ₽ </b>`,
        `(<s>${oldPrice}</s>)`,
        ' ', sale, '\n',
        url,
        // `<a href="${url}">Сайт</a>`,
        
    ].join('');
    return text;
 }

function sendUpdate(item, oldPrice) { 
    const { price } = item;
    const changeType = price > oldPrice
        ? 'увеличилась'
        : 'уменьшилась'
    
    const before = [
        'Цена ', changeType, ' ', '\n',
        oldPrice, ' -> ', price
    ]
    const text = getTex(item, before);
    Users.forEach(chatId => bot.sendMessage(chatId, text, options))
}

function sendNewItem(item) { 
    const text = getTex(item);
    Users.forEach(chatId => bot.sendMessage(chatId, text, options))
}


module.exports = {
    sendNewItem,
    sendUpdate,
    init
}