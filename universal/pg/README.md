## Create database

### OS X

    sudo -i

If saas-go user does not exists 

    createuser -d -R -P saas-go -U postgres 

Create separated database for module universal

    createdb -O saas-go universal-test -U postgres