# order/Dockerfile
FROM openjdk:11
WORKDIR /app
COPY target/order_service.jar order_service.jar
EXPOSE 5003
CMD ["java", "-jar", "order_service.jar"]
