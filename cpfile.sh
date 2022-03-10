ssh root@andcloud.ddns.net "rm -rf ~/gopath/src/github.com/cowby123/Andcloud/*"
scp -r ./* root@andcloud.ddns.net:~/gopath/src/github.com/cowby123/Andcloud