# Etapa de construcción
# FROM maven:3.9.2-eclipse-temurin-17 AS build
# WORKDIR /app
# COPY pom.xml .
# COPY src ./src
# RUN mvn clean package -DskipTests

# # Verifica si el JAR se ha creado correctamente
# RUN ls -l target/

# # Etapa de ejecución
# FROM openjdk:17-jdk-slim
# WORKDIR /app
# COPY --from=build /app/target/*.jar order_service.jar
# COPY env.sh /app/env.sh
# RUN chmod +x /app/env.sh
# EXPOSE 5003
# CMD ["sh", "-c", "source /app/env.sh && java -jar order_service.jar"]


# Etapa de construcción
FROM maven:3.9.2-eclipse-temurin-17 AS build
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

# Verifica si el JAR se ha creado correctamente
RUN ls -l target/

# Etapa de ejecución
FROM openjdk:17-jdk-slim
WORKDIR /app
COPY --from=build /app/target/*.jar order_service.jar
COPY env.sh /app/env.sh
RUN chmod +x /app/env.sh
EXPOSE 5003
CMD ["sh", "-c", "sleep 30 & . /app/env.sh && java -jar order_service.jar"]
