from django.shortcuts import render
from django.http import JsonResponse
from .models import Product

def product_list(request):
    products = list(Product.objects.values())
    return JsonResponse(products, safe=False)
