# Payments Gateway

This is a fairly simple payment's gateway, it follows 2 possible cycles of Auth->Void and Auth->capture->refund.

### Assumptions made
- If this service can understand the request and store it in db it means the customer has been charged. Having mock calls to another party to get an approval or to do the operations would cause this to become a much larger piece of work.
- HTTP errors imply a decline, this could be made nicer by having response bodies along with the http failure

### Areas for improvement
- logging using a log library like zap or logrus to create a lot more meaningful logs and make it easy to write to file most the logs act as debug statements for now, with a logging library you can have logging fields and log for example the transaction ID along with a message, and using a log aggregator tool see the lifespan of a transaction
- Idempotency on all requests
- Client Authentication
- Config loader needed for things like secrets etc. so It's not hard coded, and then we can pass db credentials as config instead
- Database integration testing
- Separate dockerfiles for testing locally vs 'production' we can add a mock for the testing files if external services need to be called during integration tests
- The code to decide to accept based on pan could have been re-framed as a fraud engine or something similar to give more context on why it was declined
- This was all done in http, grpc would be a nicer option
- Recovery scenarios, in the case perhaps the service has an issue after storing to the db in the real world we'd have to reverse the operation to make sure our 3rd parties are aware of the final state, could be done in an async fashion (publishing webhooks for example to the client to let them know it was successful) or simply just having a method to undo the operation depending on business needs
- A table/db for card information storage the pan is only stored for the purposes of fulfilling the requirements of the task, to take this further a separate table or db which owns all PCI information would be nice in reality a different service
- Using a state machine to handle transitions and technically 'charging' the customer, the service layer could then be responsible for more business specific logic like the Luhn check and the Pan check etc
- Metrics 
- Validation, at the moment this is only checking if currencies and amounts make sense, validation on all input fields should be used such as the CVV, and some kind of sanitization on every single request that comes in just to make sure there are no vulnerabilities 
- Small QoL improvement could be a makefile helpful once there are a lot more commands to run 

### Setup 
Setup is quite simple if you want to run this code simply do 
- docker-compose up --build \
This api will be available on localhost:8080 for any requests

### Testing
Unit tests: go test ./... \
All tests: go test -tags=integration ./...
