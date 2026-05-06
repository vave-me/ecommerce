+-----------+         +------------+          +---------+
|  Client   |         |  gRPC/HTTP |          | Stripe  |
| (Browser) | <-----> |  Gateway   | <------> | API     |
+-----------+         +------------+          +---------+
       |                     |                     |
       | 1. Authorize Payment|                     |
       |-------------------->|                     |
       |                     | 2. Process Command  |
       |                     |-------------------->|
       |                     |                     | 3. Create Payment Intent
       |                     |                     |<--------------------|
       |                     | 4. Save Payment     |
       |                     |<--------------------|
       |                     |                     |
       | 5. Confirm Payment  |                     |
       |-------------------->|                     |
       |                     | 6. Process Command  |
       |                     |-------------------->|
       |                     |                     | 7. Confirm Payment Intent
       |                     |                     |<--------------------|
       |                     | 8. Update Payment   |
                                                                                                     |                     |<--------------------|
                                                                                                     |                     | 9. Emit Events      |
                                                                                                     |                     |-------------------->|
                                                                                                     |                     |                     |
                                                                                                     |                     |                     | (Event Handlers Act)
                                                                                                     |                     |                     |
                                                                                                     +-----------+         +------------+          +---------+
                                                                                                     |  Client   |         | Application |          | Event   |
                                                                                                     |           | <-----> | Services    | <------> | Handlers|
                                                                                                     +-----------+         +------------+          +---------+
