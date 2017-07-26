"""untitled1 URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/1.11/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  url(r'^$', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  url(r'^$', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.conf.urls import url, include
    2. Add a URL to urlpatterns:  url(r'^blog/', include('blog.urls'))
"""
from django.conf.urls import url, include, static
import app.views as vs
import untitled1.settings as s

urlpatterns = [
    url(r'^$', vs.react_home),
    url(r'^create/', vs.react_create),
    url(r'^show/', vs.react_show),
    url(r'^delete/(.*?)/', vs.react_delete),
    url(r'^download/(.*)', vs.react_download),
    url(r'login/', include('app.loginurls')),
]