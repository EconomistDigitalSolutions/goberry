if [ -z "${SERVICE_NAME}" ]; then
  echo "Cleaning up local environment..."
  docker rmi -f $(docker images | grep 'goberry' | awk {'print $3'} | uniq)
else
  docker rmi -f $(docker images | grep ${SERVICE_NAME} | awk {'print $3'} | uniq)
fi
