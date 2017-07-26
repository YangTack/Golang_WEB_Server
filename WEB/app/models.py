from django.db import models


class ModelsOne(models.Model):
    text = models.TextField(max_length=10)
    file = models.FileField(upload_to='./rec/', unique=True)
