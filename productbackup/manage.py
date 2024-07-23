# product/manage.py
import os
import sys

if __name__ == "__main__":
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "product.settings")
    try:
        from django.core.management import execute_from_command_line
    except ImportError as exc:
        raise ImportError(
            "Couldn't import Django. Are you sure it's installed and "
            "available on your PYTHONPATH environment variable? Did you "
            "forget to activate a virtual environment?"
        ) from exc
    execute_from_command_line(sys.argv)

# # product/product/settings.py
INSTALLED_APPS = [
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'catalog',
]

from product.settings import *
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.mysql',
        'NAME': 'ms_product',
        'USER': 'root',
        'PASSWORD': 'root',
        'HOST': 'product_db',
        'PORT': '3306',
    }
}

# product/catalog/models.py
from django.db import models

class Product(models.Model):
    name = models.CharField(max_length=100)
    description = models.TextField()
    price = models.DecimalField(max_digits=10, decimal_places=2)

# product/catalog/views.py
from django.shortcuts import render
from django.http import JsonResponse
from product.catalog.models import Product

def product_list(request):
    products = list(Product.objects.values())
    return JsonResponse(products, safe=False)

# product/urls.py
from django.urls import path
from product.catalog.views import product_list

urlpatterns = [
    path('products/', product_list),
]

# product/settings.py
ROOT_URLCONF = 'urls'