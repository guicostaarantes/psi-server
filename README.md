# Psi server

This is the main server application for Psi.

Psi is a project that aims to provide low-cost psychological treatment to people that are not able to afford the high prices of the market, and also to help people in vulnerable situations with volunteering, totally free appointments.

Feel free to open an issue or create a pull request. Please follow the rules:
- Be gentle with your colleagues 😊.
- Write code and documentation in English.
- Write issues and pull requests in English.
- Use variable names that explain what that var/interface/struct is supposed to do.
- Prepend every commit message with an emoji 😎 to help others understand what you are doing there (use https://gitmoji.dev as a reference).

## Instructions to help you get going

### Domain

- Domain logic is divided in modules, each one represented by a folder inside `./modules`.
- Inside each module, there is a models folder for representing the entities of the module.
- Inside each module, there is a services folder for executing business rules.
- Each service is a struct that receives the utils it needs to operate (dependency inversion is very important for testing).
- Each service has only one method, called Execute, that consumes the utils and runs the business logic of the service.

### Infrastructure

- Infrastructure logic stays inside `./utils`.
- Everytime you need a library to make something work (e.g. password hashing), instead of importing the library in the service you need, you have to create a util for that functionality.
- Each util stays inside a different folder in `./utils`.
- Each util must have an interface in a file called i_<name_of_util>.go. This interface dictates what that util can do.
- In the same folder, create the implementations of the util. Check created utils for reference.
- If the util requires an external service to run (e.g. database), you must also create a mock implementation that shall be used on testing.

### GraphQL

- The external world interacts with this application using GraphQL. All files related to this stay in `./graph`.
- After creating a service, you probably want users to trigger it via GraphQL. Let's allow that by changing the schema.
- Schema stays in `./graph/schema`. For convenience the full schema is broken in files with the same name as the modules they relate with the most.
- Create the new queries/mutations/subscriptions for your users to interact with. Then run `./gqlgen.sh`.
- Check that changes were applied to files in `./graph/resolvers`. New queries/mutations/subscriptions will lead to new methods in the resolvers, initially with no content, only a panic for "not implemented" which you should overwrite.
- Before linking your services to the resolvers, you must initialize the new services in `./graph/resolvers/resolver.go`. Create a new field for it inside the resolver struct, and create a method for getting or setting this service. Check the ones that are there already, is pretty straightforward.
- Now, go to the `./graph/resolvers/*.resolvers.go`, find the methods that resolve your new queries/mutations/subscriptions, and attach the necessary services via the instances you created in `./graph/resolvers/resolver.go`.

### Testing

- Currently, there are no unit tests for each service. If you'd like to contribute on it, that would be awesome 😁.
- Unit tests must stay in the same folder as the file they are testing. They also must have the same name plus the suffix "_test".
- There is a end to end test in `./e2e/e2e_test.go`. It's a junky file that tests a real-world scenario, but it works pretty well.
- After you created utils, services, schema and resolvers, please add tests to the new functionality in `./e2e/e2e_test.go`.
- Running `./test-e2e.sh` will execute the tests and output the coverage percentage. We are currently at 75% coverage. Most of the non-covered code is error handling which I don't really know how to tackle. If you do, please let me know.
- Tests should be fast. In my PC, `./test-e2e.sh` finishes in 45 milliseconds. Run the tests in your PC before you start coding as set it as a benchmark. Then run again after your functionality is ready and the new tests are created. If the time increases too much, it is likely that something is not working the way it should.
- Tests must not produce side effects. Despite speed, that's another reason to create mock utils for things such as databases, mail sending, etc.
- If you want to debug the e2e tests in VSCode, run the `Launch e2e tests` debug configuration. 

### Running locally

- The instructions for docker to run all the containers we need is inside `./docker-compose.yml`. If you have docker installed in your PC, run `docker compose up --build` and it will build this app and run it alongside the other containers needed.
- The containers that run in this file are:
  - `mongo`: the database for our application
  - `mongo-express`: a web application to explore mongo collections
  - `mailhog`: a SMTP server that will intercept all mails sent by it and show in a web application
- After docker compose is up and running, go to your web browser and open 3 tabs:
  - `localhost:8080`: GraphQL Playground to execute queries/mutations/subscriptions
  - `localhost:8081`: Mongo Express Web Application to monitor mongo collections
  - `localhost:8082`: Mailhog Web Application to monitor mails being sent

### Debugging locally

- If you want to debug your code in VSCode, you need to run `docker compose --file docker-compose-debug.yml up --build` instead.
- After it finishes loading, run the `Attach server` debug configuration.
- When the attachment to the debug session finishes, you will see `connect to http://localhost:8080/ for GraphQL playground` logged in the terminal.