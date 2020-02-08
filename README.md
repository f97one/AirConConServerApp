# AirConCon

The app which manages the operation schedule of IoT for air conditioner control.

## before building the app

You need to create ssh key pair named by 'airconcon_jwt_rsa'.

```shell script
$ ssh-keygen -t rsa -f airconcon_jwt_rsa -P ""
```

You need to convert the public key format to PKCS8.

```shell script
$ ssh-keygen -f airconcon_jwt_rsa.pub -e -m PKCS8 > airconcon_jwt_rsa.pub.pkcs8
```

## Licensing

This app is under licensed by 2 Clause BSD license.