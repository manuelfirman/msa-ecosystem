# user/Dockerfile
FROM node:16
WORKDIR /app

# Copiar y instalar dependencias
COPY package*.json ./
RUN npm install

# Copiar el resto del código
COPY . .

# Instalar dockerize
RUN apt-get update && apt-get install -y curl \
    && curl -sSLO https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
    && tar xzf dockerize-linux-amd64-v0.6.1.tar.gz -C /usr/local/bin \
    && rm dockerize-linux-amd64-v0.6.1.tar.gz

# Exponer el puerto
EXPOSE 5001

# Configurar el punto de entrada para esperar a la base de datos
CMD ["dockerize", "-wait", "tcp://user_db:3306", "-timeout", "30s", "node", "src/index.js"]
