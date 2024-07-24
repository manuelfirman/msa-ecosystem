# Microservices Architecture Project

This repository showcases a comprehensive Microservices Architecture project that integrates various technologies to create a scalable and modular system. The project includes a Load Balancer, Authentication Service, User Service, Django Service, and a Java Service. Below is an overview of each component and the technology stack used.

## Project Overview

The system is designed to demonstrate how different microservices can work together in a distributed environment. Each service is responsible for a specific functionality and communicates with other services through HTTP requests.

## Technology Stack

### Load Balancer :traffic_light:
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
- **Language:** Go
- **Description:** Handles incoming requests and distributes them to the appropriate service. Implements load balancing and request forwarding.

<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,mysql,bash,docker" />
  </a>
</p>

### Authentication Service :lock:
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
- **Language:** Go
- **Description:** Manages user authentication, including registration and login functionalities. Secures endpoints and manages user sessions.
<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,mysql,bash,docker" />
  </a>
</p>

### User Service :busts_in_silhouette:
![Nodejs](https://img.shields.io/badge/Node.js-339933?style=flat&logo=nodedotjs&logoColor=white)
![Express](https://img.shields.io/badge/Express.js-000000?style=flat&logo=express&logoColor=white)
- **Language:** Node.js
- **Framework:** Express
- **Database:** MySQL (via Sequelize)
- **Description:** Manages user data and profiles. Provides RESTful APIs to interact with user information.
<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=nodejs,express,mysql,sequelize,docker" />
  </a>
</p>

### Product Service :package:
![Python](https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white)
![Django](https://img.shields.io/badge/Django-092E20?style=flat&logo=django&logoColor=white)
- **Language:** Python
- **Framework:** Django
- **Description:** Provides additional functionalities and serves specific business logic. Implements a RESTful API to interact with external services.
<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=python,django,sqlite,bash,docker" />
  </a>
</p>


### Order Service :shopping_cart:

![Java](https://img.shields.io/badge/Java-007396?style=flat&logo=java&logoColor=white)
![Spring Boot](https://img.shields.io/badge/Spring%20Boot-6DB33F?style=flat&logo=springboot&logoColor=white)
- **Language:** Java
- **Framework:** Spring Boot
- **Description:** Handles various business processes and provides a RESTful API for communication with other services.
<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=java,spring,bash,docker" />
  </a>
</p>


### Notifications Service :bell:
![Node.js](https://img.shields.io/badge/Node.js-339933?style=flat&logo=node.js&logoColor=white)
![Express](https://img.shields.io/badge/Express.js-000000?style=flat&logo=express&logoColor=white)
- **Language:** Node.js
- **Framework:** Express.js
- **Description:** A future service for handling notifications, to be implemented with Node.js and Express.

<p align="center">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=nodejs,express,mysql,sequelize,docker" />
  </a>
</p>


## Getting Started ▶️ 
1. **Clone the Repository**

```sh
git clone https://github.com/yourusername/your-repository.git
cd your-repository
```

2. **Configure Environment Variables**
- Review the `.env_example` and `env_example.sh` files.
- Set up your environment variables accordingly by creating a `.env` file or exporting the variables from `env_example.sh`.

3. **Set Up Docker Containers**
- Ensure Docker and Docker Compose are installed.
- Build and start the containers.
```sh
docker-compose up --build
```

