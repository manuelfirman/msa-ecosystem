# product/Dockerfile
FROM python:3.10
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
EXPOSE 5002

# Comando para ejecutar migraciones y levantar el servidor
CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:5002"]
