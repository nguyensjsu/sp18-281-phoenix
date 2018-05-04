var fs = require('fs');
var Client = require('node-rest-client').Client;
const syncClient = require('sync-rest-client');

var orderEndpoint = "http://ec2-user@ec2-52-52-199-60.us-west-1.compute.amazonaws.com:3000/order/";

module.exports = {
    getPage: function (page) {
      switch (classify(page)) {
        case "/order.html":
          return getOrderPage(page)
          break
        default:
          return fs.readFileSync("." + page)
      }
    }
};

function getOrderPage(page) {
  var url = classify(page)
  var id = extractData(page)
  var orderPage = fs.readFileSync("." + url, 'utf8');
  return processOrderPage(orderPage, id);
}

function processOrderPage(orderPage, id) {
  var orderTemplate = fs.readFileSync("./templates/order.template", 'utf8');
  var itemTemplate = fs.readFileSync("./templates/item.template", 'utf8');
  var order = getOrder(id)
  orderTemplate = orderTemplate.replaceAll("${orderID}", order.Id)
  orderTemplate = orderTemplate.replace("${totalPrice}", order.TotalPrice)

  var items = ""
  for (var i = 0; i < order.Items.length; i++ ) {
    var item = itemTemplate
    item = item.replace("${image}", getImage(order.Items[i].Drink))
    item = item.replace("${drink}", order.Items[i].Drink)
    item = item.replace("${price}", order.Items[i].Price)

    item = item.replace("${wid}", "w" + i)
    item = item.replace("${sid}", "s" + i)
    item = item.replace("${rsid}", "rs" + i)
    item = item.replace("${rtid}", "rt" + i)
    item = item.replace("${rgid}", "rg" + i)
    item = item.replaceAll("${optname}", "opt" + i)

    if (order.Items[i].Options.includes("Whipped Cream")) {
      item = item.replace("${whipped}", "checked")
    }
    if (order.Items[i].Options.includes("Low Sugar")) {
      item = item.replace("${low}", "checked")
    }

    if (order.Items[i].Size.includes("Short")) {
      item = item.replace("${short}", "checked")
    }
    else if (order.Items[i].Size.includes("Tall")) {
      item = item.replace("${tall}", "checked")
    }
    else if (order.Items[i].Size.includes("Grande")) {
      item = item.replace("${grande}", "checked")
    }

    items += item
  }

  orderTemplate = orderTemplate.replace("${items}", items)
  orderPage = orderPage.replace("${order}", orderTemplate)

  return Buffer.from(orderPage, 'utf8');
}

function getOrder(id) {
  var response = syncClient.get(orderEndpoint + id);
  return response.body
}

String.prototype.replaceAll = function(target, replacement) {
  return this.split(target).join(replacement);
};

function classify(page) {
  return page.split("?")[0]
}

function extractData(page) {
  return page.split("?")[1].split("=")[1]
}

function getImage(drink) {
  switch(drink) {
    case "Blonde Caffe Latte":
      return "images/blonde_caffe_latte.jpg"
    case "Ice Beverage":
      return "images/ice_beverage.jpg"
    case "Latte Macchiato":
      return "images/Latte_Macchiato.jpg"
    case "Strawberry Creme":
      return "images/strawberry_creme.jpg"
  }
}
