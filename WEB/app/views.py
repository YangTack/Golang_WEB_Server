from django.shortcuts import render, HttpResponse, redirect, Http404, HttpResponseRedirect
from django.http import StreamingHttpResponse
from django.contrib.auth import login, authenticate
import app.models as mo
import os
import json


def is_login(func):
    def inner(*args, **kwargs):
        if args[0].user.is_authenticated:
            return func(*args, **kwargs)
        else:
            return redirect(react_login)

    return inner


def react_login(request):
    if request.method == 'GET':
        return render(request, 'login.html')
    else:
        username = request.POST.get('user', '')
        passwd = request.POST.get('passwd', '')
        user = authenticate(username=username, password=passwd)
        if user is not None and user.is_active:
            login(request, user)
            return redirect(react_home)
        else:
            return render(request, 'login.html')


@is_login
def react_home(request):
    if request.method == 'GET':
        argv = ['error', 'direct']
        context = {}
        for i in argv:
            if i in request.session.keys():
                context[i] = request.session[i]
                del request.session[i]
        return render(request, '_form.html', context)
    else:
        return HttpResponse(status=404)


@is_login
def react_create(request):
    if request.method == 'POST':
        text = request.POST['text']
        file = request.FILES.get('file', '')
        if len(file) == 0 and len(text) == 0:
            request.session['error'] = '没有输入任何信息'
            return redirect(react_home)
        if len(mo.ModelsOne.objects.filter(file__exact=file)) != 0:
            request.session['error'] = '文件已存在'
            return redirect(react_home)
        if len(file) == 0:
            data = mo.ModelsOne(text=text, file='NOFILE')
        else:
            data = mo.ModelsOne(text=text, file=file)
        data.save()
        request.session['direct'] = True
        return redirect(react_home)
    else:
        return HttpResponse(status=400)


@is_login
def react_show(request):
    if request.method == 'GET':
        data = mo.ModelsOne.objects.all()
        return render(request, 'showall.html', {'mod': data})
    else:
        return HttpResponse(status=400)


@is_login
def react_delete(request, id):
    if request.method == 'POST':
        if len(id) == 0:
            request.session['direct'] = True
            return redirect(react_home)
        else:
            obj = mo.ModelsOne.objects.filter(id=int(id, 10))
            for i in obj:
                if i.file.name != 'NOFILE' and os.path.exists(i.file.path):
                    os.remove(i.file.path)
                i.delete()
            request.session['direct'] = True
            return redirect(react_home)
    else:
        return HttpResponse(status=404)


def download_func(request, path):
    def file_iter(file_path, chuck=1024):
        file_path = file_path.strip('/')
        with open('./' + file_path, 'rb') as f:
            while True:
                rec = f.read(chuck)
                if rec == b'':
                    raise StopIteration
                else:
                    yield rec

    response = StreamingHttpResponse(file_iter(path))
    response['Content-Type'] = 'application/octet-stream'
    response['Content-Disposition'] = 'attachment;filename="{0}"'.format(path.split('/')[-1])
    response['Content-Length'] = os.path.getsize('./' + path.strip('/'))
    return response


@is_login
def react_download(request, path):
    if request.method == 'GET':
        def file_iter(file_path, chuck=1024):
            file_path = file_path.strip('/')
            with open('./' + file_path, 'rb') as f:
                while True:
                    rec = f.read(chuck)
                    if rec == b'':
                        raise StopIteration
                    else:
                        yield rec

        response = StreamingHttpResponse(file_iter(path))
        response['Content-Type'] = 'application/octet-stream'
        response['Content-Disposition'] = 'attachment;filename="{0}"'.format(path.split('/')[-1])
        response['Content-Length'] = os.path.getsize('./' + path.strip('/'))
        return response
    else:
        return HttpResponse(status=404)
