const axios = require("axios");

class Client {
  constructor(host = 'http://localhost:8080') {
    this.host = host;
  }

  startBasket(userCustomerId) {
    return axios.post(`${this.host}/api/baskets`, {userCustomerId})
  }

  addItem(basketId, productId, quantity = 1) {
    return axios.put(
      `${this.host}/api/baskets/${basketId}/addItem`,
      {productId, quantity}
    )
  }
}

module.exports = {Client};
