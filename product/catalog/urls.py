from django.urls import path, include
from rest_framework import routers
from catalog import views

router = routers.DefaultRouter()
router.register(r'categories', views.CategoryViewSet)
router.register(r'products', views.ProductViewSet)
router.register(r'images', views.ProductImageViewSet)

urlpatterns = [
    path('', include(router.urls)),
]
