const axios = require('axios').default;
const { parse } = require('node-html-parser');
const _ = require('lodash');
const fs = require('fs');

const domain = 'https://www.divan.ru';

const getUrl = (path) => `${domain}${path}`

async function getPageItems(path, selector) {
    const { data } = await axios.get(getUrl(path))
    return parse(data).querySelectorAll(selector);
}

async function getPagesPath() {
    const res = await getPageItems('/category/stok-mebeli?categories[]=2', '.dqBvL');
    const [{ childNodes }] = res;
    const pages = childNodes.reduce((t, item) => {
        const href = _.get(item, 'attrs.href') || '';
        if (href.includes('/category/stok-mebeli') && !t.includes(href))
            t.push(href)
        return t

    }, [])

    return pages

}


function getData(element) {
    const name = _.get(element, 'childNodes[0].childNodes[0].text')?.trim();
    const href = _.get(element, 'childNodes[0].attrs.href');
    const price = _.get(element, 'childNodes[2].childNodes[0].childNodes[0].text')?.trim();
    const oldPrice = _.get(element, 'childNodes[2].childNodes[1].childNodes[0].text')?.trim();
    const sale = _.get(element, 'childNodes[2].childNodes[2].childNodes[0].childNodes[0].text')?.trim();

    const url = getUrl(href);
    const atr = href.split('art--').pop()

    return {
        name, sale, url, atr,
        price: Number(price?.replace(' ', '') || 0),
        oldPrice: Number(oldPrice?.replace(' ', '') || 0),
        date: new Date(),
    }
}


async function getNewData() {
    // const pagesPaths = await getPagesPath();

    const pagesPaths = [
        '/category/stok-mebeli?categories[]=2',
        '/category/stok-mebeli/page-2?categories[]=2',
        '/category/stok-mebeli/page-3?categories[]=2'
    ];

    const items = [];

    for (const path of pagesPaths) {
        console.log('page', path);
        const pageItems = await getPageItems(path, '.lsooF')
        items.push(...pageItems)
    }

    return items.map(getData);

}

async function writeDB() {
    const data = await getNewData();
    console.log(data);
    fs.writeFileSync('./db.json', JSON.stringify(data, null, 2))


}

module.exports = {
    getNewData,
    writeDB
}