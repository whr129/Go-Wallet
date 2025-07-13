# Go-Wallet
Inspired by the simple bank project on Udemy, develop a microservice based fullstack project Go Wallet.<br>

## App Overview
Go wallet support basic account/ transaction functions<br>
Include email notification for lower balance<br>
Role-based management, include Admin, customer service, users<br>

## Service Breakdown

### Auth Service
User Info and Auth:
-   **JWT Auth** 
-   **Redis Cache**

### Wallet Service
Account Info and balances:
-   **Gin API**
-   **RabbitMQ** for low balance notification

### Transaction Service
Transactions and Trades:
-   **PostgreSQL Transactions**
-   **Handle deadlock**
-   **Concurrency Unit Test**

### Notification Service
Email Notification:
-   **RabbitMQ** for low balance notification

## Technical Stack
Frontend: React & Ant Design for frontend to build a user-friendly frontend<br>
MicroService: Using gRPC for internal communication between distributed services + MQ to notify user with email<br>
Backend: Go & Gin for a microservice based backend system<br>
Database: PostgreSQL for data storage<br>
Cache: Using Redis to store cache/ user login activity<br>
Message Queue: RabbitMQ for messaging and notification<br>
Deployment: docker, Kubernetes, AWS free tier<br>

## Future
Add MangoDB<br>
Add real-world crypto transaction<br>
Add logger<br>

## Reference
Go & Gin & gRPC & Docker & Deployment: https://github.com/techschool/simplebank<br>
Microservices & RabbitMQ: https://github.com/SmoothWay/udemy-go-microservices<br>
Architecture & Implementation: https://github.com/JordanMarcelino/learn-go-microservices/<br>
