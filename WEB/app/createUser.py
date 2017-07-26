from django.contrib.auth.models import User
import getopt
import sys

user_name = ''
email = ''
password = ''
is_staff = False
is_superuser = False
opts, args = getopt.getopt(sys.argv[1:-1], 'stu:p:e:', ['username', 'email', 'password', 'staff', 'superuser'])
for op, val in opts:
    if op == '-s' or op == '--superuser':
        is_superuser = True
    elif op == '-t' or op == '--staff':
        is_staff = True
    elif op == '-u' or op == '--username':
        user_name = val
    elif op == '-p' or op == '--password':
        password = val
    elif op == '-e' or op == '--email':
        email = val
if user_name == '':
    print('username must have a value')
    sys.exit(0)
elif password == '':
    print('password must have a value')
    sys.exit(0)
User.objects.create_user(username=user_name, email=email if email != '' else '/', password=password, is_staff=is_staff,
                         is_superuser=is_superuser)
