import '@testing-library/jest-dom';

// Mock nats.ws
jest.mock('nats.ws', () => ({
  connect: jest.fn().mockResolvedValue({
    closed: () => Promise.resolve(),
    close: jest.fn().mockResolvedValue(undefined),
    status: () => ({
      [Symbol.asyncIterator]: async function* () {
        yield { type: 'connect' };
      }
    }),
    subscribe: jest.fn().mockReturnValue({
      [Symbol.asyncIterator]: async function* () {
        // Mock subscription
      },
      unsubscribe: jest.fn(),
      closed: Promise.resolve()
    }),
    publish: jest.fn(),
    request: jest.fn().mockResolvedValue({ data: new Uint8Array() }),
    jetstream: jest.fn().mockReturnValue({
      publish: jest.fn().mockResolvedValue({})
    }),
    jetstreamManager: jest.fn().mockResolvedValue({})
  })
}));