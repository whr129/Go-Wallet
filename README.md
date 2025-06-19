# Go-Wallet
Inspired by the simple bank project on Udemy, develop a microservice based fullstack project Go Wallet.

# App Overview
Go wallet support basic account/ transaction functions<br>
Include email notification for lower balance<br>
Role-based management, include Admin, customer service, users

# Service Breakdown
User Service: JWT Auth + User info + redis<br>
Wallet Service: Account/balances<br>
Transaction Service: tansactions<br>
Notification Service: Email Notification<br>

# Technical Stack
Frontend: React & Ant Design for frontend to build a user-friendly frontend<br>
Backend: Go & Gin for a microservice based backend system<br>
Database: PostgreSQL for data storage<br>
Cache: Using Redis to store cache/ user login activity<br>
Message Queue: RabbitMQ/Kafka (pending)<br>
Other: gRPC, docker, Kubernetes, AWS free tier<br>

# Future
Add MangoDB<br>
Add real-world crypto transaction<br>
