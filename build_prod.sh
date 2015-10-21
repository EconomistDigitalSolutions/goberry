cd service
GOOS=linux GOARCH=amd64 godep go build -o ../goberry -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`"