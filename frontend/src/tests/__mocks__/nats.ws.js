// Mock for nats.ws module
module.exports = {
  connect: jest.fn(() => Promise.resolve({
    subscribe: jest.fn(),
    publish: jest.fn(),
    drain: jest.fn(() => Promise.resolve()),
    close: jest.fn(() => Promise.resolve()),
  })),
  consumerOpts: jest.fn(() => ({
    deliverTo: jest.fn(),
    durable: jest.fn(),
    deliverAll: jest.fn(),
    manualAck: jest.fn(),
    ackExplicit: jest.fn(),
  })),
}; 